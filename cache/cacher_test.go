package cache

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

func TestCacherPlatformDir(t *testing.T) {
	origCacheDir := cacheDir
	cacheDir = "/tmp/tldr-test-cache"
	defer func() { cacheDir = origCacheDir }()

	c := &cacher{plat: platform.LINUX, lang: "en", name: "ls.md"}
	want := "/tmp/tldr-test-cache/pages.en/linux"
	if got := c.platformDir(); got != want {
		t.Errorf("platformDir() = %q, want %q", got, want)
	}
}

func TestCacherFile(t *testing.T) {
	origCacheDir := cacheDir
	cacheDir = "/tmp/tldr-test-cache"
	defer func() { cacheDir = origCacheDir }()

	c := &cacher{plat: platform.OSX, lang: "es", name: "ls.md"}
	want := "/tmp/tldr-test-cache/pages.es/osx/ls.md"
	if got := c.file(); got != want {
		t.Errorf("file() = %q, want %q", got, want)
	}
}

func TestCacherCmd(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"ls.md", "ls"},
		{"ls", "ls"},
		{"docker-compose.md", "docker-compose"},
	}

	for _, c := range cases {
		cr := &cacher{name: c.name}
		if got := cr.cmd(); got != c.want {
			t.Errorf("cmd() for name %q = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestCacherCreateDirReadDir(t *testing.T) {
	origCacheDir := cacheDir
	cacheDir = t.TempDir()
	defer func() { cacheDir = origCacheDir }()

	c := &cacher{plat: platform.COMMON, lang: "en", name: "foo.md"}
	c.createDir()

	if _, err := os.Stat(c.platformDir()); err != nil {
		t.Fatalf("createDir() did not create %q: %v", c.platformDir(), err)
	}

	if err := os.WriteFile(filepath.Join(c.platformDir(), "foo.md"), []byte("content"), 0600); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	entries := c.readDir()
	var found bool
	for _, e := range entries {
		if e.Name() == "foo.md" {
			found = true
		}
	}
	if !found {
		t.Errorf("readDir() did not list seeded file foo.md, got %v", entries)
	}
}

func TestCacherFind(t *testing.T) {
	origCacheDir := cacheDir
	cacheDir = t.TempDir()
	defer func() { cacheDir = origCacheDir }()

	c := &cacher{plat: platform.COMMON, lang: "en", name: "foo.md"}

	if f := c.find(); f != nil {
		f.Close()
		t.Errorf("find() = non-nil for missing file, want nil")
	}

	c.createDir()
	if err := os.WriteFile(c.file(), []byte("content"), 0600); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	f := c.find()
	if f == nil {
		t.Fatalf("find() = nil, want file")
	}
	defer f.Close()
	if filepath.Base(f.Name()) != "foo.md" {
		t.Errorf("find() returned file %q, want basename foo.md", f.Name())
	}
}

func TestCacherSearchFreshCacheHit(t *testing.T) {
	origCacheDir := cacheDir
	origConfig := config.Config
	cacheDir = t.TempDir()
	// PagesURI is unreachable on purpose: if search() incorrectly decided
	// the cache was stale and tried to download, this would fail fast
	// instead of silently hitting the real network.
	config.Config = config.Options{
		ExtendedSearch:  false,
		CacheExpiration: 30,
		PagesURI:        "http://127.0.0.1:1",
	}
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
	}()

	c := &cacher{plat: platform.COMMON, lang: "en", name: "foo.md"}
	c.createDir()
	if err := os.WriteFile(c.file(), []byte("cached content"), 0600); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	f := c.search()
	if f == nil {
		t.Fatalf("search() = nil, want cached file")
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read cached file: %v", err)
	}
	if string(buf) != "cached content" {
		t.Errorf("search() returned content %q, want %q", buf, "cached content")
	}
}

func TestCacherSearchStaleTriggersRedownload(t *testing.T) {
	origCacheDir := cacheDir
	origConfig := config.Config
	cacheDir = t.TempDir()
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
	}()

	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Path != "/common/foo.md" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fresh content"))
	}))
	defer srv.Close()

	config.Config = config.Options{
		ExtendedSearch:  false,
		CacheExpiration: 0,
		PagesURI:        srv.URL,
	}

	// get is left false here: the stale-cache re-download in search() is
	// not gated by the get flag at all, only extendedSearch()'s
	// cache-miss download is. get:false demonstrates that.
	c := &cacher{plat: platform.COMMON, lang: "en", name: "foo.md"}
	c.createDir()
	if err := os.WriteFile(c.file(), []byte("stale content"), 0600); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	staleTime := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(c.file(), staleTime, staleTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	f := c.search()
	if f == nil {
		t.Fatalf("search() = nil, want re-downloaded file")
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(buf) != "fresh content" {
		t.Errorf("search() returned content %q, want %q", buf, "fresh content")
	}
	if hits != 1 {
		t.Errorf("server hit count = %d, want 1", hits)
	}
}

func TestCacherSearchCacheMissDownloadsWhenGetTrue(t *testing.T) {
	origCacheDir := cacheDir
	origConfig := config.Config
	cacheDir = t.TempDir()
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
	}()

	var hits []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits = append(hits, r.URL.Path)
		if r.URL.Path == "/linux/foo.md" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("downloaded content"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	config.Config = config.Options{
		ExtendedSearch:  true,
		CacheExpiration: 30,
		PagesURI:        srv.URL,
	}

	c := &cacher{plat: platform.LINUX, lang: "en", name: "foo.md", get: true}

	f := c.search()
	if f == nil {
		t.Fatalf("search() = nil, want downloaded file")
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(buf) != "downloaded content" {
		t.Errorf("search() returned content %q, want %q", buf, "downloaded content")
	}
	if len(hits) == 0 {
		t.Errorf("expected at least one HTTP request, got none")
	}
}

func TestCacherSearchCacheMissWithoutExtendedSearchNeverDownloads(t *testing.T) {
	// Characterizes existing behavior (not fixed, out of scope): a
	// cache-miss only ever triggers a download from within
	// extendedSearch(), which itself only runs when
	// config.Config.ExtendedSearch is true. With ExtendedSearch=false,
	// get=true has no effect and search() returns nil without attempting
	// any download.
	origCacheDir := cacheDir
	origConfig := config.Config
	cacheDir = t.TempDir()
	config.Config = config.Options{
		ExtendedSearch:  false,
		CacheExpiration: 30,
		PagesURI:        "http://127.0.0.1:1",
	}
	defer func() {
		cacheDir = origCacheDir
		config.Config = origConfig
	}()

	c := &cacher{plat: platform.COMMON, lang: "en", name: "missing.md", get: true}
	f := c.search()
	if f != nil {
		f.Close()
		t.Errorf("search() = non-nil, want nil (ExtendedSearch=false should skip download even with get=true)")
	}
}

// cacher.remove() with a non-"clearall" name cannot be exercised here: it
// calls log.Fatal (os.Exit(1)) if the page isn't cached, and unconditionally
// calls os.Exit(0) at the end of that branch when it is -- regardless of the
// exit parameter, which only gates the clearall branch below. Either path
// would kill the test binary, so it is intentionally left untested.

func TestCacherRemoveClearAll(t *testing.T) {
	origCacheDir := cacheDir
	defer func() { cacheDir = origCacheDir }()
	cacheDir = t.TempDir()

	seedFile := filepath.Join(cacheDir, "pages.en", "common", "foo.md")
	if err := os.MkdirAll(filepath.Dir(seedFile), 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(seedFile, []byte("x"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	c := &cacher{name: "clearall.md"}
	c.remove(false)

	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		t.Errorf("remove(clearall) did not remove cacheDir, stat err = %v", err)
	}
}
