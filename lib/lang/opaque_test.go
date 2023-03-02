package lang

import (
	"testing"

	"github.com/hashicorp/nomad/ci"
	"github.com/shoenig/test/must"
)

func Test_OpaqueMapsEqual(t *testing.T) {
	ci.Parallel(t)

	type public struct {
		A int
	}

	type private struct {
		a int
	}

	type mix struct {
		A int
		b int
	}

	cases := []struct {
		name string
		a, b map[string]any
		exp  bool
	}{{
		name: "both nil",
		a:    nil,
		b:    nil,
		exp:  true,
	}, {
		name: "empty and nil",
		a:    nil,
		b:    make(map[string]any),
		exp:  true,
	}, {
		name: "same strings",
		a:    map[string]any{"a": "A"},
		b:    map[string]any{"a": "A"},
		exp:  true,
	}, {
		name: "same public struct",
		a:    map[string]any{"a": &public{A: 42}},
		b:    map[string]any{"a": &public{A: 42}},
		exp:  true,
	}, {
		name: "different public struct",
		a:    map[string]any{"a": &public{A: 42}},
		b:    map[string]any{"a": &public{A: 10}},
		exp:  false,
	}, {
		name: "different private struct",
		a:    map[string]any{"a": &private{a: 42}},
		b:    map[string]any{"a": &private{a: 10}},
		exp:  true, // private fields not compared
	}, {
		name: "mix same public different private",
		a:    map[string]any{"a": &mix{A: 42, b: 1}},
		b:    map[string]any{"a": &mix{A: 42, b: 2}},
		exp:  true, // private fields not compared
	}, {
		name: "mix different public same private",
		a:    map[string]any{"a": &mix{A: 42, b: 1}},
		b:    map[string]any{"a": &mix{A: 10, b: 1}},
		exp:  false,
	}, {
		name: "nil empty slice values",
		a:    map[string]any{"a": []string(nil)},
		b:    map[string]any{"a": make([]string, 0)},
		exp:  true,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := OpaqueMapsEqual(tc.a, tc.b)
			must.Eq(t, tc.exp, result)
		})
	}
}
