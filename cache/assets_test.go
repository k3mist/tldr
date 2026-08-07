package cache

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/djr2/tldr/config"
)

// getAssets() builds `page := Pages{}` directly instead of via NewPage(), so
// Pages' unexported cfg field is always the zero value config.Options{} --
// config.Config.ZipURI is silently ignored and the hardcoded zipUri constant
// (a real github.com URL) is always requested instead. This is a genuine bug
// (flagged, not fixed -- out of scope for this task). Because of it, pointing
// config.Config.ZipURI at a local httptest server -- the approach that works
// for cacher.go's download path via NewPage() -- would not work here and
// would instead make an uncontrolled real request to github.com. Instead,
// http.DefaultClient's transport is swapped for a fake one that answers
// in-memory, which both proves the bug (by inspecting the URL actually
// requested) and lets the extraction logic be verified without any real
// network access.

type fakeZipRoundTripper struct {
	requestedURL string
	body         []byte
}

func (rt *fakeZipRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.requestedURL = req.URL.String()
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(rt.body)),
		Header:     make(http.Header),
	}, nil
}

// failingRoundTripper fails the test immediately if it is ever invoked. It
// is used instead of the "unreachable URI" trick because that trick relies
// on config.Config.ZipURI actually being read -- which, per the bug above,
// it is not -- so it would provide no real protection here.
type failingRoundTripper struct {
	t *testing.T
}

func (rt *failingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.t.Fatalf("unexpected HTTP request during getAssets() fresh-cache branch: %s", req.URL.String())
	return nil, nil
}

func buildTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatalf("zip create %q: %v", name, err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatalf("zip write %q: %v", name, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func TestGetAssetsDownloadsAndExtracts(t *testing.T) {
	origCacheDir := cacheDir
	origConfig := config.Config
	origClient := http.DefaultClient
	cacheDir = t.TempDir()
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
		http.DefaultClient = origClient
	}()

	config.Config = config.Options{ZipURI: "http://127.0.0.1:1/should-be-ignored.zip"}

	zipBytes := buildTestZip(t, map[string]string{
		"pages/common/foo.md":    "foo docs",
		"pages.es/common/bar.md": "bar docs",
	})

	rt := &fakeZipRoundTripper{body: zipBytes}
	http.DefaultClient = &http.Client{Transport: rt}

	getAssets()

	if rt.requestedURL != zipUri {
		t.Errorf("getAssets() requested %q, want hardcoded %q (config.Config.ZipURI %q was ignored, see bug note above)", rt.requestedURL, zipUri, config.Config.ZipURI)
	}

	foo, err := os.ReadFile(filepath.Join(cacheDir, "pages", "common", "foo.md"))
	if err != nil {
		t.Fatalf("read extracted foo.md: %v", err)
	}
	if string(foo) != "foo docs" {
		t.Errorf("foo.md content = %q, want %q", foo, "foo docs")
	}

	bar, err := os.ReadFile(filepath.Join(cacheDir, "pages.es", "common", "bar.md"))
	if err != nil {
		t.Fatalf("read extracted bar.md: %v", err)
	}
	if string(bar) != "bar docs" {
		t.Errorf("bar.md content = %q, want %q", bar, "bar docs")
	}

	if _, err := os.Stat(filepath.Join(cacheDir, "assets.zip")); err != nil {
		t.Errorf("assets.zip not written to cacheDir: %v", err)
	}
}

func TestGetAssetsSkipsWhenFresh(t *testing.T) {
	origCacheDir := cacheDir
	origConfig := config.Config
	origClient := http.DefaultClient
	cacheDir = t.TempDir()
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
		http.DefaultClient = origClient
	}()

	config.Config = config.Options{ZipURI: "http://127.0.0.1:1/unused.zip"}

	zipFile := filepath.Join(cacheDir, "assets.zip")
	if err := os.WriteFile(zipFile, []byte("stale placeholder"), 0600); err != nil {
		t.Fatalf("seed assets.zip: %v", err)
	}
	now := time.Now()
	if err := os.Chtimes(zipFile, now, now); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	http.DefaultClient = &http.Client{Transport: &failingRoundTripper{t: t}}

	getAssets()

	buf, err := os.ReadFile(zipFile)
	if err != nil {
		t.Fatalf("read assets.zip: %v", err)
	}
	if string(buf) != "stale placeholder" {
		t.Errorf("assets.zip was modified, want untouched, got %q", buf)
	}
}
