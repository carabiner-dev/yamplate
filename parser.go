// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package yamlplate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	yaml2 "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"
)

var variableRegex *regexp.Regexp

const variableRegexCode = `\$\{\s*([-A-Z0-9a-z_]+)\s*\}`

// Unmarshal emulates the yaml Unmarshal method but adding support for
// template substitutions.
func Unmarshal(in []byte, out interface{}, optsFncs ...DecoderOption) error {
	opts := Options{}
	for _, f := range optsFncs {
		f(&opts)
	}

	replacedin := ""
	scanner := bufio.NewScanner(bytes.NewReader(in))
	errs := []error{}
	for scanner.Scan() {
		subs := extractLineVariables(scanner.Text())
		line, ee := replaceSusbstitutions(scanner.Text(), &opts.Variables, subs)
		errs = append(errs, ee...)
		replacedin += line + "\n"
	}

	if err := errors.Join(errs...); err != nil {
		return err
	}

	switch opts.UnmarshalerVersion {
	case "", "2":
		return yaml2.Unmarshal([]byte(replacedin), out)
	case "3":
		return yaml3.Unmarshal([]byte(replacedin), out)
	default:
		return fmt.Errorf("Invalid yaml parser version (need to be 2 or 3)")
	}
}

// replaceSusbstitutions takes a string, a translation table and a substitution
// map and replaces the template entries with the values in the substitution map.
// It returs errors for all failed substitutions.
//
//nolint:gocritic // Yes we want a pointer here
func replaceSusbstitutions(line string, vars *map[string]string, subs []varSub) (string, []error) {
	errs := []error{}
	var value string
	var ok bool

	for _, sub := range subs {
		if value, ok = (*vars)[sub.Name]; !ok {
			errs = append(errs, fmt.Errorf("no variable substitution defined for %q", sub.Name))
			continue
		}
		line = strings.Replace(line, sub.Replace, value, 1)
	}
	return line, errs
}

// Decode reads the yaml data from the configured reader and decodes the YAML
// data using a real yaml decoder to z
func (d *DecoderV2) Decode(z any) error {
	return replaceAndDecode(z, d.Decoder, &d.pointers, &d.Options)
}

// Decode reads the yaml data from the configured reader and decodes the YAML
// data using a real yaml decoder to z
func (d *DecoderV3) Decode(z any) error {
	return replaceAndDecode(z, d.Decoder, &d.pointers, &d.Options)
}

// replaceAndDecode reads the yaml, data, replaces the variables and decodes
// the resulting yaml
func replaceAndDecode(z any, yd yamlDecoder, pts *decoderIoPointers, opts *Options) error {
	varErrs := []error{}

	go func() {
		for pts.scanner.Scan() {
			line := pts.scanner.Text()

			// Replace the lines as we read them
			subs := extractLineVariables(line)
			line, errs := replaceSusbstitutions(line, &opts.Variables, subs)
			varErrs = append(varErrs, errs...)

			if _, err := pts.writer.Write([]byte(line + "\n")); err != nil {
				varErrs = append(varErrs, fmt.Errorf("writing susbtituted line %w", err))
				return
			}
		}
		pts.writer.Close()
	}()

	if err := yd.Decode(z); err != nil {
		varErrs = append(varErrs, err)
	}

	return errors.Join(varErrs...)
}

type varSub struct {
	Name    string
	Replace string
}

// extractLineVariables reads a string and locates the template variable
// substitutions a list of variable names and the string to replace.
func extractLineVariables(line string) []varSub {
	if variableRegex == nil {
		variableRegex = regexp.MustCompile(variableRegexCode)
	}

	res := variableRegex.FindAllStringSubmatch(line, -1)
	ret := []varSub{}
	for _, m := range res {
		ret = append(ret, varSub{
			Name:    m[1],
			Replace: m[0],
		})
	}
	return ret
}
