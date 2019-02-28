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
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
)

type ExtraColumnsPrinter struct {
	columns []Column
	encoder runtime.Encoder
	decoder runtime.Decoder
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
		columns: columns,
		decoder: decoder,
	}
	return printer, nil
}

// PrintObj prints the obj in a human-friendly format according to the type of the obj.
func (e *ExtraColumnsPrinter) PrintObj(obj runtime.Object, output io.Writer) error {
	fmt.Println("ExtraColumns:", e.columns)

	// parsers := make([]*jsonpath.JSONPath, len(e.columns))
	// for ix := range e.columns {
	// 	parsers[ix] = jsonpath.New(fmt.Sprintf("column%d", ix)).AllowMissingKeys(true)
	// 	if err := parsers[ix].Parse(e.columns[ix].FieldSpec); err != nil {
	// 		return err
	// 	}
	// }

	// for _, column := range e.columns {
	// 	parser := jsonpath.New("sorting").AllowMissingKeys(true)
	// 	err := parser.Parse(column.FieldSpec)
	// 	if err != nil {
	// 		return fmt.Errorf("sorting error: %v", err)
	// 	}
	// }

	// fmt.Println("Obj:", obj)
	return nil
}
