package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("failed to close pipe reader: %v", err)
	}

	return buf.String()
}

func TestLogWriterWrite(t *testing.T) {
	debugFlag := flagSet.Lookup("debug")
	origDebug := debugFlag.Value.String()
	t.Cleanup(func() {
		if err := debugFlag.Value.Set(origDebug); err != nil {
			t.Fatalf("failed to restore debug flag: %v", err)
		}
	})

	t.Run("debug disabled prints raw text", func(t *testing.T) {
		if err := debugFlag.Value.Set("disable"); err != nil {
			t.Fatalf("failed to set debug flag: %v", err)
		}

		var w logWriter
		var n int
		var writeErr error
		out := captureStdout(t, func() {
			n, writeErr = w.Write([]byte("hello"))
		})

		if writeErr != nil {
			t.Fatalf("Write returned error: %v", writeErr)
		}
		if out != "hello" {
			t.Errorf("output = %q, want %q", out, "hello")
		}
		if n != len(out) {
			t.Errorf("n = %d, want %d (len of output written)", n, len(out))
		}
	})

	t.Run("debug enabled prepends UTC timestamp", func(t *testing.T) {
		if err := debugFlag.Value.Set("enable"); err != nil {
			t.Fatalf("failed to set debug flag: %v", err)
		}

		var w logWriter
		var n int
		var writeErr error
		out := captureStdout(t, func() {
			n, writeErr = w.Write([]byte("hello"))
		})

		if writeErr != nil {
			t.Fatalf("Write returned error: %v", writeErr)
		}
		if !strings.HasPrefix(out, "[") {
			t.Errorf("output = %q, want it to start with %q", out, "[")
		}
		if !strings.Contains(out, "] hello") {
			t.Errorf("output = %q, want it to contain %q", out, "] hello")
		}
		if n != len(out) {
			t.Errorf("n = %d, want %d (len of output written)", n, len(out))
		}
	})

	t.Run("any non-disable value is treated as enabled", func(t *testing.T) {
		if err := debugFlag.Value.Set("verbose"); err != nil {
			t.Fatalf("failed to set debug flag: %v", err)
		}

		var w logWriter
		out := captureStdout(t, func() {
			w.Write([]byte("hi"))
		})

		if !strings.HasPrefix(out, "[") || !strings.Contains(out, "] hi") {
			t.Errorf("output = %q, want timestamp-prefixed %q since debug != %q", out, "hi", "disable")
		}
	})
}

func TestSetLogDebug(t *testing.T) {
	debugFlag := flagSet.Lookup("debug")
	origDebug := debugFlag.Value.String()
	origLogFlags := log.Flags()
	t.Cleanup(func() {
		if err := debugFlag.Value.Set(origDebug); err != nil {
			t.Fatalf("failed to restore debug flag: %v", err)
		}
		log.SetFlags(origLogFlags)
	})

	if err := debugFlag.Value.Set("disable"); err != nil {
		t.Fatalf("failed to set debug flag: %v", err)
	}
	setLogDebug()
	if got := log.Flags(); got != 0 {
		t.Errorf("log.Flags() after setLogDebug() with debug=disable = %d, want %d", got, 0)
	}

	if err := debugFlag.Value.Set("enable"); err != nil {
		t.Fatalf("failed to set debug flag: %v", err)
	}
	setLogDebug()
	if got := log.Flags(); got != log.Lshortfile {
		t.Errorf("log.Flags() after setLogDebug() with debug=enable = %d, want %d", got, log.Lshortfile)
	}
}
