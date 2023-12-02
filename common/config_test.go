package common

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

type Config struct {
	StringField        *string `json:"StringField"`
	AnotherStringField *string `json:"AnotherStringField"`
	IntegerField       *int    `json:"IntegerField"`
}

func (c *Config) String() string {
	return fmt.Sprintf("Config {%q %q %d}", *(c.StringField), *(c.AnotherStringField), *(c.IntegerField))
}

func ptr[T any](value T) *T {
	result := new(T)
	*result = value
	return result
}

func TestConfig(t *testing.T) {

	for _, tc := range []struct {
		name string
		file string
		args []string
		want interface{}
	}{
		{
			name: "flag only config",
			args: []string{
				"--StringField=SomeValue",
				"--AnotherStringField=/another/string/field/",
				"--IntegerField=123",
			},
			want: &Config{
				StringField:        ptr("SomeValue"),
				AnotherStringField: ptr("/another/string/field/"),
				IntegerField:       ptr(123),
			},
		},
		{
			name: "file only config",
			file: "sample_config.json",
			args: []string{},
			want: &Config{
				StringField:        ptr("Whatever"),
				AnotherStringField: ptr("Anyway"),
				IntegerField:       ptr(12345678),
			},
		},
		{
			name: "all flags and file config",
			file: "sample_config.json",
			args: []string{
				"--StringField=SomeValue",
				"--AnotherStringField=/another/string/field/",
				"--IntegerField=123",
			},
			want: &Config{
				StringField:        ptr("SomeValue"),
				AnotherStringField: ptr("/another/string/field/"),
				IntegerField:       ptr(123),
			},
		},
		{
			name: "some flags and file config",
			file: "sample_config.json",
			args: []string{
				"--StringField=SomeValue",
				"--IntegerField=123",
			},
			want: &Config{
				StringField:        ptr("SomeValue"),
				AnotherStringField: ptr("Anyway"),
				IntegerField:       ptr(123),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.file != "" {
				tc.args = append(tc.args, "--config", filepath.Join("testdata", tc.file))
			}
			config, err := GetConfig[Config](tc.args, "config")
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			} else if !reflect.DeepEqual(tc.want, config) {
				t.Errorf("Expected config for file %q and flags %q should be %v, but got %v",
					tc.file, tc.args, tc.want, config)
			}
		})
	}
}
