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

package printers

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
)

type ExtraColumnsPrinter struct {
	options PrintOptions
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewExtraColumnPrinter creates a ExtraColumnPrinter.
// If encoder and decoder are provided, an attempt to convert unstructured types to internal types is made.
func NewExtraColumnsPrinter(decoder runtime.Decoder, options PrintOptions) *ExtraColumnsPrinter {
	printer := &ExtraColumnsPrinter{
		options: options,
		decoder: decoder,
	}
	return printer
}

// PrintObj prints the obj in a human-friendly format according to the type of the obj.
func (h *ExtraColumnsPrinter) PrintObj(obj runtime.Object, output io.Writer) error {
	fmt.Println("Printing, yay!")
	return nil
}
