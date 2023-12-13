package common

import (
	"reflect"
	"testing"
)

type someStruct struct {
	SomeInt    int    `json:"some_int"`
	SomeString string `json:"some_string"`
}

func TestFromJson(t *testing.T) {

	for _, tc := range []struct {
		name    string
		json    string
		want    *someStruct
		wantErr bool
	}{
		{
			name: "normal json",
			json: "{\"some_int\": 1234, \"some_string\": \"abcd\"}",
			want: &someStruct{1234, "abcd"},
		},
		{
			name:    "broken json",
			json:    "{wefwef",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FromJson[someStruct]([]byte(tc.json))
			if (err != nil && !tc.wantErr) || (err == nil && tc.wantErr) {
				t.Errorf("Unexpected error or missing an expected error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("expected FromJson to return %s, but got %s", tc.want, result)
			}
		})
	}
}
