/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package get

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/jsonpath"
)

type ExtraColumnsPrinter struct {
	Columns   []Column
	Encoder   runtime.Encoder
	Decoder   runtime.Decoder
	NoHeaders bool
	// lastType records type of resource printed last so that we don't repeat
	// header while printing same type of resources.
	lastType reflect.Type
}

// NewExtraColumnPrinter creates a ExtraColumnPrinter.
// If encoder and decoder are provided, an attempt to convert unstructured types to internal types is made.
func NewExtraColumnsPrinter(decoder runtime.Decoder, spec []string) (*ExtraColumnsPrinter, error) {
	columns := make([]Column, len(spec))
	for ix := range spec {
		colSpec := strings.Split(spec[ix], ":")
		if len(colSpec) != 2 {
			return nil, fmt.Errorf("unexpected extra-columns spec: %s, expected <header>:<json-path-expr>", spec[ix])
		}
		spec, err := RelaxedJSONPathExpression(colSpec[1])
		if err != nil {
			return nil, err
		}
		columns[ix] = Column{Header: colSpec[0], FieldSpec: spec}
	}

	printer := &ExtraColumnsPrinter{
		Columns: columns,
		Decoder: decoder,
	}
	return printer, nil
}

// PrintObj prints the obj in a human-friendly format according to the type of the obj.
func (e *ExtraColumnsPrinter) PrintObj(obj runtime.Object, output io.Writer) error {
	//fmt.Println("ExtraColumns:", e.Columns)
	// fmt.Println("Object type:", reflect.TypeOf(obj))

	parsers := make([]*jsonpath.JSONPath, len(e.Columns))
	for ix := range e.Columns {
		parsers[ix] = jsonpath.New(fmt.Sprintf("column%d", ix)).AllowMissingKeys(true)
		if err := parsers[ix].Parse(e.Columns[ix].FieldSpec); err != nil {
			return err
		}
	}

	// Print default headers
	e.PrintTableHeaders(obj.(*metav1beta1.Table), output)

	// Print extra-column headers
	objType := reflect.TypeOf(obj)
	if !e.NoHeaders && objType != e.lastType {
		headers := make([]string, len(e.Columns))
		for ix := range e.Columns {
			headers[ix] = e.Columns[ix].Header
		}
		fmt.Fprintln(output, strings.Join(headers, "\t"))
		e.lastType = objType
	}

	includesTable := false
	includesRuntimeObjs := false

	switch t := obj.(type) {
	case *metav1beta1.Table:
		includesTable = true

		// Print columns
		if err := ParseAndPrint(t, parsers, output); err != nil {
			return err
		}
	default:
		includesRuntimeObjs = true
	}

	if includesRuntimeObjs && includesTable {
		return fmt.Errorf("sorting is not supported on mixed Table and non-Table object lists")
	}

	return nil
}

func (e *ExtraColumnsPrinter) PrintTableHeaders(table *metav1beta1.Table, output io.Writer) error {
	if !e.NoHeaders {
		// avoid printing headers if we have no rows to display
		if len(table.Rows) == 0 {
			return nil
		}

		first := true
		for _, column := range table.ColumnDefinitions {
			// if !e.Wide && column.Priority != 0 {
			if column.Priority != 0 {
				continue
			}
			if first {
				first = false
			} else {
				fmt.Fprint(output, "\t")
			}
			fmt.Fprint(output, strings.ToUpper(column.Name))
		}
		fmt.Fprintln(output)
	}

	return nil
}

func ParseAndPrint(table *metav1beta1.Table, parsers []*jsonpath.JSONPath, output io.Writer) error {
	columns := make([]string, len(parsers))

	for i, row := range table.Rows {

		// Print extra-columns
		for ix := range parsers {
			parser := parsers[ix]
			valueStrings := []string{}

			var values [][]reflect.Value
			var err error

			// Parse the JSON values from the table
			values, err = findJSONPathResults(parser, table.Rows[i].Object.Object)

			if err != nil {
				return err
			}

			if len(values) == 0 || len(values[0]) == 0 {
				valueStrings = append(valueStrings, "<none>")
			}
			for arrIx := range values {
				for valIx := range values[arrIx] {
					valueStrings = append(valueStrings, fmt.Sprintf("%v", values[arrIx][valIx].Interface()))
				}
			}
			columns[ix] = strings.Join(valueStrings, ",")
		}
		fmt.Fprintln(output, strings.Join(columns, "\t"))

		// Print default columns
		first := true
		for i, cell := range row.Cells {
			if i >= len(table.ColumnDefinitions) {
				// https://issue.k8s.io/66379
				// don't panic in case of bad output from the server, with more cells than column definitions
				break
			}
			column := table.ColumnDefinitions[i]
			// if !options.Wide && column.Priority != 0 {
			if column.Priority != 0 {
				continue
			}
			if first {
				first = false
			} else {
				fmt.Fprint(output, "\t")
			}
			if cell != nil {
				fmt.Fprint(output, cell)
			}
		}
		fmt.Fprintln(output)
	}

	return nil
}
