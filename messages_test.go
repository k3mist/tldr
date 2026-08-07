package main

import (
	"strings"
	"testing"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

func TestBanner(t *testing.T) {
	orig := config.Config
	t.Cleanup(func() { config.Config = orig })
	config.Config = config.Options{}

	out := captureStdout(t, func() {
		banner()
	})

	if out == "" {
		t.Fatal("banner() produced no output")
	}
	if !strings.Contains(out, "https://tldr.sh") {
		t.Errorf("banner() output = %q, want it to contain %q", out, "https://tldr.sh")
	}
}

func TestVersion(t *testing.T) {
	orig := config.Config
	t.Cleanup(func() { config.Config = orig })
	config.Config = config.Options{}

	out := captureStdout(t, func() {
		version()
	})

	if out == "" {
		t.Fatal("version() produced no output")
	}
	if !strings.Contains(out, "2.1.0") {
		t.Errorf("version() output = %q, want it to contain the hardcoded version literal %q", out, "2.1.0")
	}
	if !strings.Contains(strings.ToLower(out), "version") {
		t.Errorf("version() output = %q, want it to contain %q (case-insensitive)", out, "version")
	}
}

func TestNoTldr(t *testing.T) {
	orig := config.Config
	t.Cleanup(func() { config.Config = orig })
	config.Config = config.Options{}

	out := captureStdout(t, func() {
		noTldr("frobnicate")
	})

	if !strings.Contains(out, "frobnicate") {
		t.Errorf("noTldr(%q) output = %q, want it to contain the cmd", "frobnicate", out)
	}
}

func TestNoLookup(t *testing.T) {
	origConfig := config.Config
	origWarn := flagWarn
	t.Cleanup(func() {
		config.Config = origConfig
		flagWarn = origWarn
	})

	t.Run("silent when both LookupWarnings and flagWarn are false", func(t *testing.T) {
		config.Config = config.Options{LookupWarnings: false}
		flagWarn = false

		out := captureStdout(t, func() {
			noLookup("frobnicate", "en", platform.LINUX)
		})

		if out != "" {
			t.Errorf("noLookup() output = %q, want empty output", out)
		}
	})

	t.Run("prints when flagWarn is true", func(t *testing.T) {
		config.Config = config.Options{LookupWarnings: false}
		flagWarn = true

		out := captureStdout(t, func() {
			noLookup("frobnicate", "en", platform.LINUX)
		})

		if out == "" {
			t.Fatal("noLookup() produced no output, want warning text since flagWarn is true")
		}
		if !strings.Contains(out, "frobnicate") {
			t.Errorf("noLookup() output = %q, want it to contain cmd %q", out, "frobnicate")
		}
		if !strings.Contains(out, "en") {
			t.Errorf("noLookup() output = %q, want it to contain lang %q", out, "en")
		}
		if !strings.Contains(out, platform.LINUX.String()) {
			t.Errorf("noLookup() output = %q, want it to contain platform %q", out, platform.LINUX.String())
		}
	})

	t.Run("prints when LookupWarnings is true", func(t *testing.T) {
		config.Config = config.Options{LookupWarnings: true}
		flagWarn = false

		out := captureStdout(t, func() {
			noLookup("frobnicate", "en", platform.LINUX)
		})

		if out == "" {
			t.Fatal("noLookup() produced no output, want warning text since LookupWarnings is true")
		}
		if !strings.Contains(out, "frobnicate") {
			t.Errorf("noLookup() output = %q, want it to contain cmd %q", out, "frobnicate")
		}
	})
}
