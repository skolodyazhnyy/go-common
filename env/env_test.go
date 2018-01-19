package env

import (
	"os"
	"reflect"
	"testing"
)

type structure struct {
	StringVal  string  `env:"TEST_STRING"`
	IntegerVal int     `env:"TEST_INTEGER"`
	Float64Val float64 `env:"TEST_FLOAT64"`
	Level2     level2
}

type level2 struct {
	Value1 string  `env:"TEST_LEVEL2_VAL1"`
	Value2 float64 `env:"TEST_LEVEL2_VAL2"`
}
type environment map[string]string

type tests struct {
	env       environment
	recipient structure
	expected  structure
}

func TestUnmarshal(t *testing.T) {
	tests := buildTestCases()

	for _, test := range tests {
		t.Logf("Testing: %+v -> %+v\n", test.env, test.expected)

		// Setup expected environment
		for key, value := range test.env {
			os.Setenv(key, value)
		}

		Unmarshal(&test.recipient)

		if !reflect.DeepEqual(test.expected, test.recipient) {
			t.Errorf("Error populating structure. Expecting %+v got %+v", test.expected, test.recipient)
		}

		// Cleanup the environment
		for key := range test.env {
			os.Unsetenv(key)
		}
	}
}

func buildTestCases() []tests {
	return []tests{
		// Populating the keys that actually exist in the environment
		{
			env: environment{
				"TEST_STRING":      "ENV1VALUE",
				"TEST_INTEGER":     "32",
				"TEST_FLOAT64":     "12.12",
				"TEST_LEVEL2_VAL1": "Value1",
				"TEST_LEVEL2_VAL2": "88.88",
			},
			recipient: structure{},
			expected: structure{
				StringVal:  "ENV1VALUE",
				IntegerVal: 32,
				Float64Val: 12.12,
				Level2: level2{
					Value1: "Value1",
					Value2: 88.88,
				},
			},
		},
		// If a value already exists in recipient struct, it should be overridden
		// by whatever comes from the environment nontheless
		{
			env: environment{
				"TEST_STRING": "Will overwrite!",
			},
			recipient: structure{
				StringVal: "Original Value",
			},
			expected: structure{
				StringVal: "Will overwrite!",
			},
		},
	}
}
