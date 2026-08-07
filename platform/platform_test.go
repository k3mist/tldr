package platform

import (
	"io"
	"os"
	"runtime"
	"testing"
)

func TestPlatformString(t *testing.T) {
	cases := []struct {
		plat Platform
		want string
	}{
		{UNKNOWN, "unknown"},
		{COMMON, "common"},
		{OSX, "osx"},
		{WINDOWS, "windows"},
		{LINUX, "linux"},
		{ANDROID, "android"},
		{FREEBSD, "freebsd"},
		{OPENBSD, "openbsd"},
		{NETBSD, "netbsd"},
		{SUNOS, "sunos"},
		{Platform(9999), "unknown"},
	}

	for _, c := range cases {
		if got := c.plat.String(); got != c.want {
			t.Errorf("Platform(%d).String() = %q, want %q", c.plat, got, c.want)
		}
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want Platform
	}{
		{"common", COMMON},
		{"osx", OSX},
		{"windows", WINDOWS},
		{"linux", LINUX},
		{"android", ANDROID},
		{"freebsd", FREEBSD},
		{"openbsd", OPENBSD},
		{"netbsd", NETBSD},
		{"sunos", SUNOS},
	}

	for _, c := range cases {
		if got := Parse(c.in); got != c.want {
			t.Errorf("Parse(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseInvalidFallsBackToActual(t *testing.T) {
	// Parse falls back to Actual() for unrecognized input, not UNKNOWN -
	// this is surprising since the zero value/UNKNOWN would be the more
	// obvious "not found" sentinel.
	got := Parse("not-a-real-platform")
	if got != Actual() {
		t.Errorf("Parse(invalid) = %v, want Actual() = %v", got, Actual())
	}
	if got == UNKNOWN {
		t.Errorf("Parse(invalid) = UNKNOWN, want it to fall back to Actual() instead")
	}
}

func TestParseFlag(t *testing.T) {
	cases := []string{"linux", "osx", "windows", "not-a-real-platform"}

	for _, in := range cases {
		if got, want := ParseFlag(in), Parse(in); got != want {
			t.Errorf("ParseFlag(%q) = %v, want %v (same as Parse)", in, got, want)
		}
	}
}

func TestGetPlatform(t *testing.T) {
	cases := []struct {
		in   Platform
		want Platform
	}{
		{UNKNOWN, Actual()},
		{LINUX, LINUX},
		{OSX, OSX},
		{WINDOWS, WINDOWS},
		{Platform(9999), Platform(9999)},
	}

	for _, c := range cases {
		if got := GetPlatform(c.in); got != c.want {
			t.Errorf("GetPlatform(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestPlatforms(t *testing.T) {
	got := Platforms()

	if len(got) != 9 {
		t.Errorf("len(Platforms()) = %d, want 9", len(got))
	}

	for _, name := range got {
		if name == "unknown" {
			t.Errorf("Platforms() contains %q, want it excluded", "unknown")
		}
	}

	want := map[string]bool{
		"common": false, "osx": false, "windows": false, "linux": false,
		"android": false, "freebsd": false, "openbsd": false, "netbsd": false,
		"sunos": false,
	}
	for _, name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("Platforms() contains unexpected name %q", name)
			continue
		}
		want[name] = true
	}
	for name, seen := range want {
		if !seen {
			t.Errorf("Platforms() missing expected name %q", name)
		}
	}
}

func TestActual(t *testing.T) {
	var want Platform
	switch runtime.GOOS {
	case "linux":
		want = LINUX
	case "freebsd":
		want = FREEBSD
	case "netbsd":
		want = NETBSD
	case "openbsd":
		want = OPENBSD
	case "darwin":
		want = OSX
	case "solaris":
		want = SUNOS
	case "windows":
		want = WINDOWS
	default:
		want = COMMON
	}

	got := Actual()
	if got != want {
		t.Errorf("Actual() = %v, want %v (mirrors runtime.GOOS = %q)", got, want, runtime.GOOS)
	}

	if got == UNKNOWN {
		t.Errorf("Actual() = UNKNOWN, want a non-UNKNOWN platform")
	}

	valid := map[string]bool{
		"common": true, "osx": true, "windows": true, "linux": true,
		"android": true, "freebsd": true, "openbsd": true, "netbsd": true,
		"sunos": true,
	}
	if s := got.String(); !valid[s] {
		t.Errorf("Actual().String() = %q, want one of the valid platform names", s)
	}
}

func TestPrintDoesNotPanic(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() failed: %v", err)
	}
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, r)
		close(done)
	}()

	Print()

	w.Close()
	<-done
	r.Close()
}
