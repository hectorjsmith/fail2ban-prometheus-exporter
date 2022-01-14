package auth

import (
	"reflect"
	"testing"
)

func TestHashString(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"Happy path #1", "123", "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"},
		{"Happy path #2", "hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
		{"Happy path #3", "H3Ll0_W0RLD", "d58a27fe9a6e73a1d8a67189fb8acace047e7a1a795276a0056d3717ad61bd0e"},
		{"Blank string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashString(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HashString() = %v, want %v", got, tt.want)
			}
		})
	}
}
