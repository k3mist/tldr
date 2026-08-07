package language

import (
	"slices"
	"testing"

	"bitbucket.org/djr2/tldr/config"
)

// setConfigLanguage sets config.Config.Language for the duration of the
// (sub)test and restores the previous value on cleanup.
func setConfigLanguage(t *testing.T, lang string) {
	t.Helper()
	orig := config.Config.Language
	config.Config.Language = lang
	t.Cleanup(func() {
		config.Config.Language = orig
	})
}

// saveLanguages snapshots the package-level languages var so a test that
// mutates it (directly or via SetLang) can't leak state into later tests.
func saveLanguages(t *testing.T) {
	t.Helper()
	orig := languages
	t.Cleanup(func() {
		languages = orig
	})
}

func TestFindSuffix(t *testing.T) {
	cases := []struct {
		name   string
		locale []string
		want   string
	}{
		{"BR suffix kept", []string{"pt", "BR"}, "pt_BR"},
		{"BR suffix kept regardless of base", []string{"xx", "BR"}, "xx_BR"},
		// findSuffix only ever inspects keepSuffix[0] ("BR") because both
		// branches of the if/else return unconditionally on the first loop
		// iteration - the loop never advances to check "PT" or "TW". So
		// despite keepSuffix listing "PT" and "TW", those suffixes are dead
		// code and never actually get preserved. This looks like a bug, but
		// the assertions below intentionally pin the current (buggy)
		// behavior rather than the presumably-intended one.
		{"PT suffix NOT kept (bug: dead code)", []string{"pt", "PT"}, "pt"},
		{"TW suffix NOT kept (bug: dead code)", []string{"zh", "TW"}, "zh"},
		{"non-suffix two-part locale", []string{"en", "US"}, "en"},
		{"single-part locale", []string{"en"}, "en"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := findSuffix(c.locale); got != c.want {
				t.Errorf("findSuffix(%v) = %q, want %q", c.locale, got, c.want)
			}
		})
	}
}

func TestParseLang(t *testing.T) {
	cases := []struct {
		name string
		lang string
		want []string
	}{
		{"plain code", "en", []string{"en"}},
		{"locale with country", "en_US", []string{"en"}},
		{"locale with encoding suffix stripped", "en_US.UTF-8", []string{"en"}},
		{"BR suffix preserved", "pt_BR", []string{"pt_BR"}},
		{"BR suffix preserved with encoding", "pt_BR.UTF-8", []string{"pt_BR"}},
		{"PT suffix not preserved (bug)", "pt_PT", []string{"pt"}},
		{"TW suffix not preserved (bug)", "zh_TW", []string{"zh"}},
		{"multiple colon-separated locales", "en_US:fr_FR", []string{"en", "fr"}},
		{"multiple locales with BR preserved", "pt_BR:zh_TW", []string{"pt_BR", "zh"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseLang(c.lang); !slices.Equal(got, c.want) {
				t.Errorf("parseLang(%q) = %v, want %v", c.lang, got, c.want)
			}
		})
	}
}

func TestSetLangPrecedence(t *testing.T) {
	cases := []struct {
		name        string
		cfgLang     string
		envLang     string
		envLANGUAGE string
		want        []string
	}{
		{
			name:        "config.Config.Language takes precedence over both env vars",
			cfgLang:     "fr_FR",
			envLang:     "de_DE",
			envLANGUAGE: "es_ES",
			want:        []string{"fr", "en"},
		},
		{
			name:        "LANG used when config empty",
			cfgLang:     "",
			envLang:     "de_DE",
			envLANGUAGE: "es_ES",
			want:        []string{"de", "en"},
		},
		{
			name:        "LANGUAGE used when config and LANG empty",
			cfgLang:     "",
			envLang:     "",
			envLANGUAGE: "es_ES",
			want:        []string{"es", "en"},
		},
		{
			name:        "no source set falls back to en only",
			cfgLang:     "",
			envLang:     "",
			envLANGUAGE: "",
			want:        []string{"en"},
		},
		{
			name:        "en not duplicated when already present",
			cfgLang:     "en_US:fr_FR",
			envLang:     "",
			envLANGUAGE: "",
			want:        []string{"en", "fr"},
		},
		{
			name:        "multi-locale colon-separated list with en appended",
			cfgLang:     "de_DE:fr_FR:es_ES",
			envLang:     "",
			envLANGUAGE: "",
			want:        []string{"de", "fr", "es", "en"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			saveLanguages(t)
			setConfigLanguage(t, c.cfgLang)
			t.Setenv("LANG", c.envLang)
			t.Setenv("LANGUAGE", c.envLANGUAGE)

			SetLang()

			if got := GetLanguages(); !slices.Equal(got, c.want) {
				t.Errorf("SetLang() -> GetLanguages() = %v, want %v", got, c.want)
			}
			if !slices.Contains(GetLanguages(), "en") {
				t.Errorf("SetLang() result %v does not contain fallback %q", GetLanguages(), "en")
			}
		})
	}
}

func TestGetLanguage(t *testing.T) {
	saveLanguages(t)
	setConfigLanguage(t, "de_DE:fr_FR")
	t.Setenv("LANG", "")
	t.Setenv("LANGUAGE", "")

	SetLang()

	want := []string{"de", "fr", "en"}
	if got := GetLanguages(); !slices.Equal(got, want) {
		t.Fatalf("GetLanguages() = %v, want %v", got, want)
	}

	for i, w := range want {
		if got := GetLanguage(i); got != w {
			t.Errorf("GetLanguage(%d) = %q, want %q", i, got, w)
		}
	}
}

func TestGetLanguages(t *testing.T) {
	saveLanguages(t)
	setConfigLanguage(t, "pt_BR")
	t.Setenv("LANG", "")
	t.Setenv("LANGUAGE", "")

	SetLang()

	want := []string{"pt_BR", "en"}
	if got := GetLanguages(); !slices.Equal(got, want) {
		t.Errorf("GetLanguages() = %v, want %v", got, want)
	}
}

func TestHasLanguage(t *testing.T) {
	saveLanguages(t)
	setConfigLanguage(t, "de_DE:fr_FR")
	t.Setenv("LANG", "")
	t.Setenv("LANGUAGE", "")

	SetLang()

	cases := []struct {
		lang string
		want bool
	}{
		{"de", true},
		{"fr", true},
		{"en", true},
		{"es", false},
		{"de_DE", false},
	}

	for _, c := range cases {
		if got := HasLanguage(c.lang); got != c.want {
			t.Errorf("HasLanguage(%q) = %v, want %v", c.lang, got, c.want)
		}
	}
}
