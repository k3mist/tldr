package cache

import (
	"testing"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

func TestPagesURL(t *testing.T) {
	config.Config = config.Options{}

	cases := []struct {
		lang string
		want string
	}{
		{"en", "https://raw.githubusercontent.com/tldr-pages/tldr/main/pages/common/ln"},
		{"es", "https://raw.githubusercontent.com/tldr-pages/tldr/main/pages.es/common/ln"},
	}

	for _, c := range cases {
		p := NewPage("ln", c.lang, platform.COMMON)
		if got := p.url(); got != c.want {
			t.Errorf("url() for lang %q = %q, want %q", c.lang, got, c.want)
		}
	}
}
