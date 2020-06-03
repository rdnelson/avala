package main

import (
	"path/filepath"
	"testing"
)

func TestChangeExtention(t *testing.T) {
	cases := []struct {
		in_file, in_ext, want string
	}{
		{filepath.Join("test", "file"), "html", filepath.Join("test", "file.html")},
		{filepath.Join("test", "file.ext"), "html", filepath.Join("test", "file.html")},
		{filepath.Join("test.dir", "file"), "html", filepath.Join("test.dir", "file.html")},
		{filepath.Join("test.dir", "file.ext"), "html", filepath.Join("test.dir", "file.html")},
		{"", "html", ""},
	}

	for _, c := range cases {
		if got := changeExtention(c.in_file, c.in_ext); got != c.want {
			t.Errorf("changeExtention(%q, %q) == %q, expected: %q", c.in_file, c.in_ext, got, c.want)
		}
	}
}
