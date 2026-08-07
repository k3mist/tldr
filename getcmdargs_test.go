package main

import (
	"os"
	"reflect"
	"testing"
)

func withArgs(t *testing.T, args []string) {
	t.Helper()
	orig := os.Args
	t.Cleanup(func() { os.Args = orig })
	os.Args = args
}

func TestGetCmdArgs(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "single word command",
			args:     []string{"tldr", "ln"},
			wantCmd:  "ln",
			wantArgs: nil,
		},
		{
			name:     "two non-flag tokens are hyphen-joined into the command",
			args:     []string{"tldr", "docker", "compose"},
			wantCmd:  "docker-compose",
			wantArgs: nil,
		},
		{
			name:     "command followed by a long flag",
			args:     []string{"tldr", "ln", "--help"},
			wantCmd:  "ln",
			wantArgs: []string{"--help"},
		},
		{
			name:     "command followed by a short flag with its value",
			args:     []string{"tldr", "ln", "-p", "linux"},
			wantCmd:  "ln",
			wantArgs: []string{"-p", "linux"},
		},
		{
			name:     "bare short flag as the sole argument does not consume a value",
			args:     []string{"tldr", "-p"},
			wantCmd:  "",
			wantArgs: []string{"-p"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			withArgs(t, c.args)

			gotCmd, gotArgs := getCmdArgs()
			if gotCmd != c.wantCmd {
				t.Errorf("cmd = %q, want %q", gotCmd, c.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, c.wantArgs) {
				t.Errorf("args = %#v, want %#v", gotArgs, c.wantArgs)
			}
		})
	}
}

// TestGetCmdArgsShortFlagAsLastArgPanics documents a real bug in getCmdArgs.
//
// When a short flag (e.g. "-p") is the LAST token in os.Args and there is
// more than one token total after the program name, getCmdArgs unconditionally
// does `os.Args[(i+1)+1]` to grab the flag's "value" -- with no bounds check
// against the actual end of os.Args. When the short flag is the last token,
// that index is one past the end of os.Args and the function panics with
// "index out of range".
//
// The len(os.Args[1:]) > 1 guard only protects the case where the short flag
// is the *only* argument in the whole invocation (e.g. `tldr -p`, see the
// "bare short flag as the sole argument" case in TestGetCmdArgs) -- it does
// NOT protect the case where the short flag is merely the last of several
// arguments (e.g. `tldr ln -p`), which is the case exercised here.
//
// This test does not fix the bug -- it uses recover() to characterize the
// current actual behavior without crashing the rest of the test binary.
func TestGetCmdArgsShortFlagAsLastArgPanics(t *testing.T) {
	withArgs(t, []string{"tldr", "ln", "-p"})

	paniced := func() (r any) {
		defer func() { r = recover() }()
		getCmdArgs()
		return nil
	}()

	if paniced == nil {
		t.Fatal("expected getCmdArgs() to panic (index out of range) when a short flag is the last argument with other tokens present -- if this now passes, the underlying indexing bug in getCmdArgs may have been fixed and this test should be updated")
	}

	t.Logf("confirmed bug: getCmdArgs() panics on trailing short flag: %v", paniced)
}
