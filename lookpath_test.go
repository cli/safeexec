package safeexec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLookPath(t *testing.T) {
	root, wderr := os.Getwd()
	if wderr != nil {
		t.Fatal(wderr)
	}

	paths := []string{
		filepath.Join(root, "_fixtures", "nonexist"),
		filepath.Join(root, "_fixtures", "system"),
	}
	os.Setenv("PATH", strings.Join(paths, string(filepath.ListSeparator)))

	if err := os.Chdir(filepath.Join(root, "_fixtures", "cwd")); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc    string
		pathext string
		arg     string
		wants   string
		wantErr bool
	}{
		{
			desc:    "no extension",
			pathext: "",
			arg:     "ls",
			wants:   filepath.Join(root, "_fixtures", "system", "ls"+winonly(".exe")),
			wantErr: false,
		},
		{
			desc:    "with extension",
			pathext: "",
			arg:     "ls.exe",
			wants:   filepath.Join(root, "_fixtures", "system", "ls.exe"),
			wantErr: false,
		},
		{
			desc:    "with path",
			pathext: "",
			arg:     filepath.Join("..", "system", "ls"),
			wants:   filepath.Join("..", "system", "ls"+winonly(".exe")),
			wantErr: false,
		},
		{
			desc:    "with path+extension",
			pathext: "",
			arg:     filepath.Join("..", "system", "ls.bat"),
			wants:   filepath.Join("..", "system", "ls.bat"),
			wantErr: false,
		},
		{
			desc:    "no extension, PATHEXT",
			pathext: ".com;.bat",
			arg:     "ls",
			wants:   filepath.Join(root, "_fixtures", "system", "ls"+winonly(".bat")),
			wantErr: false,
		},
		{
			desc:    "with extension, PATHEXT",
			pathext: ".com;.bat",
			arg:     "ls.exe",
			wants:   filepath.Join(root, "_fixtures", "system", "ls.exe"),
			wantErr: false,
		},
		{
			desc:    "no extension, not found",
			pathext: "",
			arg:     "cat",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "with extension, not found",
			pathext: "",
			arg:     "cat.exe",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "no extension, PATHEXT, not found",
			pathext: ".com;.bat",
			arg:     "cat",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "with extension, PATHEXT, not found",
			pathext: ".com;.bat",
			arg:     "cat.exe",
			wants:   "",
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			os.Setenv("PATHEXT", tC.pathext)
			got, err := LookPath(tC.arg)

			if tC.wantErr != (err != nil) {
				t.Errorf("expects error: %v, got: %v", tC.wantErr, err)
			}
			if err != nil && !errors.Is(err, exec.ErrNotFound) {
				t.Errorf("expected exec.ErrNotFound; got %#v", err)
			}
			if got != tC.wants {
				t.Errorf("expected result %q, got %q", tC.wants, got)
			}
		})
	}
}

func winonly(s string) string {
	if runtime.GOOS == "windows" {
		return s
	}
	return ""
}
