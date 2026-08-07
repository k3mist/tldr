package cache

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/djr2/tldr/config"
	"bitbucket.org/djr2/tldr/platform"
)

func TestNewCacher(t *testing.T) {
	c := newCacher("ls", "en", platform.LINUX, true)

	if c.name != "ls.md" {
		t.Errorf("name = %q, want %q", c.name, "ls.md")
	}
	if c.lang != "en" {
		t.Errorf("lang = %q, want %q", c.lang, "en")
	}
	if c.plat != platform.LINUX {
		t.Errorf("plat = %v, want %v", c.plat, platform.LINUX)
	}
	if !c.get {
		t.Errorf("get = false, want true")
	}
}

func TestFind(t *testing.T) {
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

	dir := filepath.Join(cacheDir, "pages.en", "common")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ls.md"), []byte("ls docs"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	f, lang, plat := Find("ls", "en", platform.COMMON, false)
	if f == nil {
		t.Fatalf("Find() file = nil, want non-nil")
	}
	defer f.Close()
	if lang != "en" {
		t.Errorf("Find() lang = %q, want %q", lang, "en")
	}
	if plat != platform.COMMON {
		t.Errorf("Find() plat = %v, want %v", plat, platform.COMMON)
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(buf) != "ls docs" {
		t.Errorf("content = %q, want %q", buf, "ls docs")
	}
}

func TestRemoveClearAll(t *testing.T) {
	origCacheDir := cacheDir
	defer func() { cacheDir = origCacheDir }()
	cacheDir = t.TempDir()

	dir := filepath.Join(cacheDir, "pages.en", "common")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ls.md"), []byte("x"), 0600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// exit=false so the clearall branch returns instead of calling
	// os.Exit(0); see cacher_test.go for why the normal (non-clearall)
	// removal path can't be exercised in-process at all.
	Remove("clearall", "en", platform.COMMON, false)

	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		t.Errorf("Remove(clearall) did not remove cacheDir, stat err = %v", err)
	}
}
