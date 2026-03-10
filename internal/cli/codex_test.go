package cli

import (
	"reflect"
	"testing"
)

func TestResolveStatusLineOptionSupportsBranchOnlyPreset(t *testing.T) {
	t.Parallel()

	got, err := resolveStatusLineOption("branch-only")
	if err != nil {
		t.Fatalf("resolveStatusLineOption() error = %v", err)
	}

	want := []string{"git-branch"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("resolveStatusLineOption() = %#v, want %#v", got, want)
	}
}
