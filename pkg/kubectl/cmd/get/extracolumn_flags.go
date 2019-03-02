/*
Copyright 2018 The Kubernetes Authors.

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
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/scheme"
	"k8s.io/kubernetes/pkg/printers"
)

// ExtraColumnFlags provides default flags necessary for printing.
// Given the following flag values, a printer can be requested that knows
// how to handle printing based on these values.
type ExtraColumnFlags struct {
	ExtraColumns *[]string
	ShowKind     *bool
	ShowLabels   *bool
	SortBy       *string
	ColumnLabels *[]string

	// get.go-specific values
	NoHeaders bool

	Kind               schema.GroupKind
	AbsoluteTimestamps bool
	WithNamespace      bool
}

// AllowedFormats returns more customized formating options
func (f *ExtraColumnFlags) AllowedFormats() []string {
	return []string{"wide"}
}

// ToPrinter receives an outputFormat and returns a printer capable of
// handling human-readable output.
func (f *ExtraColumnFlags) ToPrinter(outputFormat string) (printers.ResourcePrinter, error) {

	if len(outputFormat) > 0 && outputFormat != "wide" {
		return nil, genericclioptions.NoCompatiblePrinterError{Options: f, AllowedFormats: f.AllowedFormats()}
	}

	decoder := scheme.Codecs.UniversalDecoder()

	showKind := false
	if f.ShowKind != nil {
		showKind = *f.ShowKind
	}

	showLabels := false
	if f.ShowLabels != nil {
		showLabels = *f.ShowLabels
	}

	columnLabels := []string{}
	if f.ColumnLabels != nil {
		columnLabels = *f.ColumnLabels
	}

	p, err := NewExtraColumnsPrinter(decoder, *f.ExtraColumns, printers.PrintOptions{
		Kind:          f.Kind,
		WithKind:      showKind,
		NoHeaders:     f.NoHeaders,
		Wide:          outputFormat == "wide",
		WithNamespace: f.WithNamespace,
		ColumnLabels:  columnLabels,
		ShowLabels:    showLabels,
	})

	if err != nil {
		return nil, err
	}

	return p, nil
}

// AddFlags receives a *cobra.Command reference and binds
// flags related to human-readable printing to it
func (f *ExtraColumnFlags) AddFlags(c *cobra.Command) {
	if f.ExtraColumns != nil {
		c.Flags().StringSliceVarP(f.ExtraColumns, "extra-columns", "E", *f.ExtraColumns, "Accepts a comma separated list of extra columns expressed as a spec with a JSONPath expression in the same vein as -o=custom-columns=<spec> (e.g. 'NAME:.metadata.name'). These columns will be displayed in addition to the default columns.")
	}
}

// NewExtraColumnFlags returns flags associated with
// human-readable printing, with default values set.
func NewExtraColumnFlags() *ExtraColumnFlags {
	extraColumns := []string{}

	return &ExtraColumnFlags{
		ExtraColumns: &extraColumns,
		NoHeaders:    false,
	}
}
