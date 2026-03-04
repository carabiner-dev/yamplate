// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package yamlplate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		data     string
		expect   map[string]any
		optFuncs []DecoderOption
		mustErr  bool
	}{
		{
			"nno-subs",
			"---\nid: test\nn: 3\n",
			map[string]any{"id": "test", "n": "3"},
			[]DecoderOption{},
			false,
		},
		{
			"1-subs",
			"---\nid: test\nn: 3\nsub: ${TEST}\n",
			map[string]any{"id": "test", "n": "3", "sub": "OK"},
			[]DecoderOption{WithVariables(map[string]string{"TEST": "OK"})},
			false,
		},
		{
			"2-subs",
			"---\nid: test\nn: 3\nsub: ${TEST}\nsub2: ${ MOAR_TEST }",
			map[string]any{"id": "test", "n": "3", "sub": "OK", "sub2": "YES"},
			[]DecoderOption{WithVariables(map[string]string{"TEST": "OK", "MOAR_TEST": "YES"})},
			false,
		},
		{
			"missing-value",
			"---\nid: test\nn: 3\nsub: ${TEST}\nsub2: ${ MOAR_TEST }",
			map[string]any{"id": "test", "n": "3", "sub": "OK", "sub2": "YES"},
			[]DecoderOption{WithVariables(map[string]string{"MOAR_TEST": "YES"})},
			true, // This one errs
		},
		{
			"missing-value-deduped",
			"---\nid: ${ORG}\nn: ${ORG}\nsub: ${ORG}\nsub2: ${ORG}",
			nil,
			[]DecoderOption{WithVariables(map[string]string{})},
			true, // This one errs
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := map[string]string{}

			err := Unmarshal([]byte(tc.data), &sut, tc.optFuncs...)
			if tc.mustErr {
				require.Error(t, err)
				if tc.name == "missing-value-deduped" {
					// Verify the error for ORG appears only once despite 4 uses
					require.Equal(t, 1, strings.Count(err.Error(), "no variable substitution defined for \"ORG\""))
				}
				return
			}
			require.NoError(t, err)

			for k, v := range tc.expect {
				require.Equal(t, v, sut[k])
			}
		})
	}
}
