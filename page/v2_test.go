package page

import (
	"bytes"
	"testing"

	"bitbucket.org/djr2/tldr/color"
	"bitbucket.org/djr2/tldr/config"
)

func TestPagev2LinesShiftsHeaderUnderline(t *testing.T) {
	// First call detects the "====" underline at lns[1] and drops it: it
	// copies lns[0] ("ln") into lns[1]'s slot, then reslices off the
	// original index 0. The net effect is the underline line is removed
	// while the command name becomes the new index 0. A second call is a
	// no-op since the new lns[1] ("> desc") no longer matches headerRxV2.
	p := &pagev2{
		lns: [][]byte{[]byte("ln"), []byte("===="), []byte("> desc"), []byte("    ln {{a}}")},
		buf: &bytes.Buffer{},
	}

	want := [][]byte{[]byte("ln"), []byte("> desc"), []byte("    ln {{a}}")}

	got := p.Lines()
	if len(got) != len(want) {
		t.Fatalf("first Lines() len = %d, want %d (%q)", len(got), len(want), got)
	}
	for i := range want {
		if !bytes.Equal(got[i], want[i]) {
			t.Errorf("first Lines()[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	got2 := p.Lines()
	if len(got2) != len(want) {
		t.Fatalf("second Lines() len = %d, want %d (%q)", len(got2), len(want), got2)
	}
	for i := range want {
		if !bytes.Equal(got2[i], want[i]) {
			t.Errorf("second Lines()[%d] = %q, want %q", i, got2[i], want[i])
		}
	}
}

func TestPagev2LinesNoShiftWhenSecondLineIsNotUnderline(t *testing.T) {
	lns := [][]byte{[]byte("ln"), []byte("> desc"), []byte("    ln {{a}}")}
	p := &pagev2{lns: lns, buf: &bytes.Buffer{}}

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

func TestPagev2Write(t *testing.T) {
	p := &pagev2{buf: &bytes.Buffer{}}
	p.Write([]byte("abc"))
	p.Write([]byte("def"))

	if got := p.buf.String(); got != "abcdef" {
		t.Errorf("buf = %q, want %q", got, "abcdef")
	}
}

func TestPagev2Header(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev2{lns: [][]byte{[]byte("ln")}, buf: &bytes.Buffer{}}

	want := []byte(color.ColorBold(cfg.HeaderColor) + "ln")
	if got := p.Header(); !bytes.Equal(got, want) {
		t.Errorf("Header() = %q, want %q", got, want)
	}
}

func TestPagev2Description(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev2{buf: &bytes.Buffer{}}

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

func TestPagev2Example(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev2{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{
			"match leading word",
			[]byte("Create a symbolic link:"),
			[]byte(color.Color(cfg.HyphenColor) + "- Create" + color.Color(cfg.ExampleColor) + " a symbolic link:"),
		},
		{
			"match leading digits and underscore",
			[]byte("_private1 thing"),
			[]byte(color.Color(cfg.HyphenColor) + "- _private1" + color.Color(cfg.ExampleColor) + " thing"),
		},
		{"header underline does not match", []byte("===="), nil},
		{"empty", []byte(""), nil},
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

func TestPagev2Syntax(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev2{buf: &bytes.Buffer{}}

	cases := []struct {
		name string
		line []byte
		want []byte
	}{
		{
			"exactly four spaces",
			[]byte("    ln -s {{a}}"),
			[]byte(color.Color(cfg.SyntaxColor) + "ln -s {{a}}"),
		},
		{
			"more than four spaces keeps extra indent",
			[]byte("      ln -s {{a}}"),
			[]byte(color.Color(cfg.SyntaxColor) + "  ln -s {{a}}"),
		},
		{"three spaces does not match", []byte("   ln -s {{a}}"), nil},
		{"no indent", []byte("ln -s {{a}}"), nil},
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

func TestPagev2Variable(t *testing.T) {
	config.Config = config.Options{}
	cfg := config.Config
	p := &pagev2{buf: &bytes.Buffer{}}

	line := []byte("ln {{target}} {{link}}")
	want := []byte("ln " + color.Color(cfg.VariableColor) + "target" + color.Color(cfg.SyntaxColor) + " " + color.Color(cfg.VariableColor) + "link" + color.Color(cfg.SyntaxColor))

	if got := p.Variable(line); !bytes.Equal(got, want) {
		t.Errorf("Variable(%q) = %q, want %q", line, got, want)
	}
}
