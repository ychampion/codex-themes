package codex_test

import (
	"reflect"
	"testing"

	"github.com/ychampion/codex-themes/internal/codex"
)

func TestStatusLinePresetSupportsBranchOnly(t *testing.T) {
	t.Parallel()

	got, err := codex.StatusLinePreset("branch-only")
	if err != nil {
		t.Fatalf("StatusLinePreset() error = %v", err)
	}

	want := []string{"git-branch"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("StatusLinePreset() = %#v, want %#v", got, want)
	}
}
