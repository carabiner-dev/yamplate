// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package yamlplate

import (
	"bufio"
	"io"

	yaml2 "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"
)

// Decoder aliases v2 as the default decoder
type Decoder = DecoderV2

type DecoderV2 struct {
	*yaml2.Decoder
	Options  Options
	pointers decoderIoPointers
}

type DecoderV3 struct {
	*yaml3.Decoder
	Options  Options
	pointers decoderIoPointers
}

type decoderIoPointers struct {
	scanner *bufio.Scanner
	writer  *io.PipeWriter
}

type yamlDecoder interface {
	Decode(any) error
}

// NewDecoder emulates the NewDecoder method from the YAML packages.
func NewDecoder(original io.Reader) *Decoder {
	return NewDecoderV2(original)
}

// NewDecoderYaml3 emulates the NewDecoder method from the YAML packages.
func NewDecoderV2(original io.Reader) *DecoderV2 {
	r, w := io.Pipe()
	return &DecoderV2{
		Decoder: yaml2.NewDecoder(r),
		Options: defaultOptions,
		pointers: decoderIoPointers{
			scanner: bufio.NewScanner(original),
			writer:  w,
		},
	}
}

// NewDecoderYaml3 emulates the NewDecoder method from the YAML packages.
func NewDecoderV3(original io.Reader) *DecoderV3 {
	r, w := io.Pipe()
	return &DecoderV3{
		Decoder: yaml3.NewDecoder(r),
		Options: defaultOptions,
		pointers: decoderIoPointers{
			scanner: bufio.NewScanner(original),
			writer:  w,
		},
	}
}
