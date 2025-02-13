// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package yamlplate

var defaultOptions = Options{
	Variables:          map[string]string{},
	UnmarshalerVersion: "2",
}

// Options control the behavior of the yaml templater
type Options struct {
	// Variables holds the substitution table
	Variables map[string]string

	// UnmarshalerVersion determines is the unmarshal function
	// uses gopkg.in/yaml.v2 or gopkg.in/yaml.v3. Note that this
	// does not affect the Decoder object as those embed the
	// real yaml decoder each has its own type.
	UnmarshalerVersion string
}

type DecoderOption func(*Options)

// WithVariables holds the substitution table with the variable name in
// the map key.
func WithVariables(vars map[string]string) DecoderOption {
	return func(opts *Options) {
		opts.Variables = vars
	}
}

// WithUnmarshalVersion determines if the unmarshal function uses
// gopkg.in/yaml.v2 or gopkg.in/yaml.v3 to un marshal the resulting
// yaml data.
func WithUnmarshalerVersion(version string) DecoderOption {
	return func(opts *Options) {
		opts.UnmarshalerVersion = version
	}
}
