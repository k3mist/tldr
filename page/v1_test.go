package page

import (
	"bytes"
	"testing"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

func TestPagev1Lines(t *testing.T) {
	lns := [][]byte{[]byte("# ln"), []byte("> desc")}
	p := &pagev1{lns: lns, buf: &bytes.Buffer{}}

	got := p.Lines()
	if len(got) != len(lns) {
		t.Fatalf("Lines() len = %d, want %d", len(got), len(lns))
	}
	for i := range lns {
		if !bytes.Equal(got[i], lns[i]) {
			t.Errorf("Lines()[%d] = %q, want %q", i, got[i], lns[i])
		}
	}
}

func TestPagev1Write(t *testing.T) {
	p := &pagev1{buf: &bytes.Buffer{}}
	p.Write([]byte("abc"))
	p.Write([]byte("def"))

	if got := p.buf.String(); got != "abcdef" {
		t.Errorf("buf = %q, want %q", got, "abcdef")
	}
}

func TestPagev1Header(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev1{lns: [][]byte{[]byte("# ln")}, buf: &bytes.Buffer{}}

	want := []byte(color.ColorBold(cfg.HeaderColor) + "ln")
	if got := p.Header(); !bytes.Equal(got, want) {
		t.Errorf("Header() = %q, want %q", got, want)
	}
}

func TestPagev1Description(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev1{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{"match", []byte("> Create a link."), []byte(color.Color(cfg.DescriptionColor) + "Create a link.")},
		{"no leading gt", []byte("Create a link."), nil},
		{"empty", []byte(""), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := p.Description(c.line)
			if !bytes.Equal(got, c.want) {
				t.Errorf("Description(%q) = %q, want %q", c.line, got, c.want)
			}
		})
	}
}

func TestPagev1Example(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev1{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{
			"match",
			[]byte("- Create a symbolic link:"),
			[]byte(color.Color(cfg.HyphenColor) + "- " + color.Color(cfg.ExampleColor) + "Create a symbolic link:"),
		},
		{"no hyphen prefix", []byte("Create a symbolic link:"), nil},
		{"hyphen no space", []byte("-Create"), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := p.Example(c.line)
			if !bytes.Equal(got, c.want) {
				t.Errorf("Example(%q) = %q, want %q", c.line, got, c.want)
			}
		})
	}
}

func TestPagev1Syntax(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev1{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{
			"match",
			[]byte("`ln {{target}} {{link}}`"),
			[]byte(color.Color(cfg.SyntaxColor) + "ln {{target}} {{link}}"),
		},
		{"no backticks", []byte("ln {{target}}"), nil},
		{"missing trailing backtick", []byte("`ln {{target}}"), nil},
		{"missing leading backtick", []byte("ln {{target}}`"), nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := p.Syntax(c.line)
			if !bytes.Equal(got, c.want) {
				t.Errorf("Syntax(%q) = %q, want %q", c.line, got, c.want)
			}
		})
	}
}

func TestPagev1Variable(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev1{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{
			"two variables",
			[]byte("ln {{target}} {{link}}"),
			[]byte("ln " + color.Color(cfg.VariableColor) + "target" + color.Color(cfg.SyntaxColor) + " " + color.Color(cfg.VariableColor) + "link" + color.Color(cfg.SyntaxColor)),
		},
		{"no variables", []byte("ln -s"), []byte("ln -s")},
		{
			// varRx's character class does not include '-', so a
			// hyphenated placeholder like {{-o}} is left untouched
			// rather than colorized. Characterizing current behavior.
			"hyphenated placeholder is not recognized as a variable",
			[]byte("cmd {{-o}}"),
			[]byte("cmd {{-o}}"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := p.Variable(c.line)
			if !bytes.Equal(got, c.want) {
				t.Errorf("Variable(%q) = %q, want %q", c.line, got, c.want)
			}
		})
	}
}
