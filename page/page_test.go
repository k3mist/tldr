package page

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/language"
	"bitbucket.org/djr2/tldr/platform"
)

const v1SampleContent = "# ln\n\n> Create hard or symbolic links.\n> More information: https://example.com.\n\n- Create a link to a file:\n\n`ln -s {{source}} {{dest}}`\n"

const v2SampleContent = "ln\n====\n\n> Create hard or symbolic links.\n> More information: https://example.com.\n\nCreate a link to a file:\n\n    ln -s {{source}} {{dest}}\n"

func TestMain(m *testing.M) {
	language.SetLang()
	os.Exit(m.Run())
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}

	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("w.Close() error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error: %v", err)
	}
	return buf.String()
}

func writeTempPage(t *testing.T, content string) *os.File {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "tldr-page-*")
	if err != nil {
		t.Fatalf("os.CreateTemp() error: %v", err)
	}

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("f.WriteString() error: %v", err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("f.Seek() error: %v", err)
	}
	return f
}

func TestNewDispatchV1(t *testing.T) {
	config.Config = config.Options{}
	f := writeTempPage(t, v1SampleContent)

	pg := New(f, "fr", platform.COMMON)

	if _, ok := pg.(*pagev1); !ok {
		t.Fatalf("New() returned %T, want *pagev1", pg)
	}

	out := captureStdout(t, pg.Print)

	wantSubstrings := []string{"ln", "common", "fr", "Create hard or symbolic links.", "Create a link to a file:", "source", "dest"}
	for _, want := range wantSubstrings {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q\noutput: %q", want, out)
		}
	}
}

func TestNewDispatchV2(t *testing.T) {
	config.Config = config.Options{}
	f := writeTempPage(t, v2SampleContent)

	pg := New(f, "fr", platform.COMMON)

	if _, ok := pg.(*pagev2); !ok {
		t.Fatalf("New() returned %T, want *pagev2", pg)
	}

	out := captureStdout(t, pg.Print)

	wantSubstrings := []string{"ln", "common", "fr", "Create hard or symbolic links.", "a link to a file:", "source", "dest"}
	for _, want := range wantSubstrings {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q\noutput: %q", want, out)
		}
	}
}

func TestParseV1(t *testing.T) {
	config.Config = config.Options{}
	lines := bytes.Split([]byte(v1SampleContent), []byte("\n"))
	p := &pagev1{lns: lines, buf: &bytes.Buffer{}}

	Parse(p, "fr", platform.COMMON)

	got := p.buf.String()
	wantSubstrings := []string{
		"ln",
		"common",
		"fr",
		"Create hard or symbolic links.",
		"More information: https://example.com.",
		"Create a link to a file:",
		"source",
		"dest",
	}
	for _, want := range wantSubstrings {
		if !strings.Contains(got, want) {
			t.Errorf("Parse() output missing %q\noutput: %q", want, got)
		}
	}
}

func TestParseV2(t *testing.T) {
	config.Config = config.Options{}
	lines := bytes.Split([]byte(v2SampleContent), []byte("\n"))
	p := &pagev2{lns: lines, buf: &bytes.Buffer{}}

	Parse(p, "fr", platform.COMMON)

	got := p.buf.String()
	wantSubstrings := []string{
		"ln",
		"common",
		"fr",
		"Create hard or symbolic links.",
		"More information: https://example.com.",
		"a link to a file:",
		"source",
		"dest",
	}
	for _, want := range wantSubstrings {
		if !strings.Contains(got, want) {
			t.Errorf("Parse() output missing %q\noutput: %q", want, got)
		}
	}
}

func TestToB(t *testing.T) {
	got := toB("abc")
	want := []byte("abc")
	if !bytes.Equal(got, want) {
		t.Errorf("toB(%q) = %q, want %q", "abc", got, want)
	}
}
