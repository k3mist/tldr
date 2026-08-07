package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"

	"bitbucket.org/djr2/tldr/color"
)

func withTempHome(t *testing.T) string {
	t.Helper()

	original, hadHome := os.LookupEnv("HOME")
	tempDir := t.TempDir()

	if err := os.Setenv("HOME", tempDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	// homedir.Dir() caches its result, so the cache must be cleared or the
	// HOME override above has no effect on the next Dir() call.
	homedir.Reset()

	t.Cleanup(func() {
		if hadHome {
			os.Setenv("HOME", original)
		} else {
			os.Unsetenv("HOME")
		}
		homedir.Reset()
	})

	return tempDir
}

func TestLoadCreatesDefaultConfig(t *testing.T) {
	tempDir := withTempHome(t)
	Config = Options{}

	// create() only creates the config.json file, not the .tldr directory
	// itself, so the directory must pre-exist (as it does on a real machine)
	// or Load() calls log.Fatal and kills the test process.
	dir := filepath.Join(tempDir, ".tldr")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create %s: %v", dir, err)
	}

	Load()

	f := filepath.Join(tempDir, ".tldr", "config.json")
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("expected config file at %s, got error: %v", f, err)
	}

	cases := []struct {
		name string
		got  any
		want any
	}{
		{"PagesURI", Config.PagesURI, ""},
		{"ZipURI", Config.ZipURI, ""},
		{"Language", Config.Language, ""},
		{"CacheExpiration", Config.CacheExpiration, 30},
		{"ExtendedSearch", Config.ExtendedSearch, true},
		{"LookupWarnings", Config.LookupWarnings, false},
		{"BannerColor1", Config.BannerColor1, color.Cyan},
		{"BannerColor2", Config.BannerColor2, color.Blue},
		{"TLDRColor", Config.TLDRColor, color.White},
		{"HeaderColor", Config.HeaderColor, color.Blue},
		{"HeaderDecorColor", Config.HeaderDecorColor, color.White},
		{"PlatformColor", Config.PlatformColor, color.DarkGray},
		{"PlatformAltColor", Config.PlatformAltColor, color.BrightPurple},
		{"DescriptionColor", Config.DescriptionColor, color.Normal},
		{"ExampleColor", Config.ExampleColor, color.Cyan},
		{"HyphenColor", Config.HyphenColor, color.Normal},
		{"SyntaxColor", Config.SyntaxColor, color.Red},
		{"VariableColor", Config.VariableColor, color.Normal},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestLoadReadsExistingConfig(t *testing.T) {
	tempDir := withTempHome(t)
	Config = Options{}

	dir := filepath.Join(tempDir, ".tldr")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create %s: %v", dir, err)
	}

	custom := Options{
		PagesURI:        "https://example.com/pages",
		ZipURI:          "https://example.com/zip.zip",
		Language:        "es",
		CacheExpiration: 99,
		ExtendedSearch:  false,
		LookupWarnings:  true,
	}

	b, err := json.Marshal(custom)
	if err != nil {
		t.Fatalf("failed to marshal custom config: %v", err)
	}

	f := filepath.Join(dir, "config.json")
	if err := os.WriteFile(f, b, 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", f, err)
	}

	Load()

	cases := []struct {
		name string
		got  any
		want any
	}{
		{"PagesURI", Config.PagesURI, custom.PagesURI},
		{"ZipURI", Config.ZipURI, custom.ZipURI},
		{"Language", Config.Language, custom.Language},
		{"CacheExpiration", Config.CacheExpiration, custom.CacheExpiration},
		{"ExtendedSearch", Config.ExtendedSearch, custom.ExtendedSearch},
		{"LookupWarnings", Config.LookupWarnings, custom.LookupWarnings},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}
