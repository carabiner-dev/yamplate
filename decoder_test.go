// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package yamlplate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		yaml    string
		mustErr bool
		eval    func(*testing.T, map[string]any)
		opts    Options
	}{
		{
			"no-subs",
			`---
id: Si
lana: "no"
`, false, func(t *testing.T, a map[string]any) {
				require.Equal(t, "Si", a["id"])
				require.Equal(t, "no", a["lana"])

				fmt.Printf("\nDecoded es:\n%+v\n\n", a)
			},
			Options{},
		},
		{
			"1-sub",
			`---
id: ${VAL}
lana: "no"
`, false, func(t *testing.T, a map[string]any) {
				require.Equal(t, "chido1", a["id"])
				require.Equal(t, "no", a["lana"])

				fmt.Printf("\nDecoded es:\n%+v\n\n", a)
			},
			Options{
				Variables: map[string]string{
					"VAL": "chido1",
				},
			},
		},
		{
			"2-subs-same-line",
			`---
id: ${VAL}.${OTHER}
lana: "no"
`, false, func(t *testing.T, a map[string]any) {
				require.Equal(t, "chido1.com", a["id"])
				require.Equal(t, "no", a["lana"])

				fmt.Printf("\nDecoded es:\n%+v\n\n", a)
			},
			Options{
				Variables: map[string]string{
					"VAL":   "chido1",
					"OTHER": "com",
				},
			},
		},
		{
			"2-subs-different-lines",
			`---
id: ${VAL}
lana: ${LINE2}
`, false, func(t *testing.T, a map[string]any) {
				fmt.Printf("\nDecoded es:\n%+v\n\n", a)

				require.Equal(t, "chido1", a["id"])
				require.Equal(t, "oooy!", a["lana"])
			},
			Options{
				Variables: map[string]string{
					"VAL":   "chido1",
					"LINE2": `"oooy!"`,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dec2 := NewDecoder(bytes.NewReader([]byte(tc.yaml)))
			dec2.Options = tc.opts
			a := map[string]any{}
			err := dec2.Decode(&a)
			if tc.mustErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tc.eval(t, a)
		})
	}
}

func TestVariableExtract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		line   string
		expect []varSub
	}{
		{"no-subs", " Just a line", []varSub{}},
		{"1-sub", ` Just a line ${HALO}`, []varSub{{"HALO", `${HALO}`}}},
		{"1-sub-spaces", ` Just a line ${ HALO }`, []varSub{{"HALO", `${ HALO }`}}},
		{"1-sub-mixedcase", ` Just a line ${ Halo }`, []varSub{{"Halo", `${ Halo }`}}},
		{"2-subs", ` Just a line ${HALO} ${ BYE }`, []varSub{{"HALO", `${HALO}`}, {"BYE", `${ BYE }`}}},
		{"no-dollar", ` Just a line {HALO} `, []varSub{}},
		{"nunmbers", `lana: ${LINE2}`, []varSub{{`LINE2`, `${LINE2}`}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vars := extractLineVariables(tc.line)
			require.Len(t, vars, len(tc.expect))
			for i := range vars {
				require.Equal(t, tc.expect[i].Name, vars[i].Name)
				require.Equal(t, tc.expect[i].Replace, vars[i].Replace)
			}
		})
	}
}

func TestReplaceSusbstitutions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		line    string
		expect  string
		vars    *map[string]string
		subs    []varSub
		mustErr bool
	}{
		{"no-subs", "no substitutions here", "no substitutions here", &map[string]string{}, []varSub{}, false},
		{"1-sub", "Hello, {$PLANET}!", "Hello, World!", &map[string]string{"PLANET": "World"}, []varSub{{"PLANET", `{$PLANET}`}}, false},
		{"2-sub", "{ $GREET }, {$PLANET}!", "Hello, World!", &map[string]string{"PLANET": "World", "GREET": "Hello"}, []varSub{{"PLANET", `{$PLANET}`}, {"GREET", `{ $GREET }`}}, false},
		{"missing-val", "{ $NAME }, {$CHEESWE}!", "Hello, World!", &map[string]string{"PLANET": "World", "GREET": "Hello"}, []varSub{{"PLANET", `{$PLANET}`}, {"CHEESWE", `{ $CHEESWE }`}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			line, errs := replaceSusbstitutions(tc.line, tc.vars, tc.subs)
			if tc.mustErr {
				require.True(t, len(errs) > 0)
				return
			}

			require.Equal(t, tc.expect, line)
		})
	}
}
