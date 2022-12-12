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

	if err := os.Chdir(filepath.Join(root, "_fixtures", "cwd")); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc    string
		path    []string
		pathext string
		arg     string
		wants   string
		wantErr bool
	}{
		{
			desc: "no extension",
			path: []string{
				filepath.Join(root, "_fixtures", "nonexist"),
				filepath.Join(root, "_fixtures", "system"),
			},
			pathext: "",
			arg:     "ls",
			wants:   filepath.Join(root, "_fixtures", "system", "ls"+winonly(".exe")),
			wantErr: false,
		},
		{
			desc:    "with extension",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: "",
			arg:     "ls.exe",
			wants:   filepath.Join(root, "_fixtures", "system", "ls.exe"),
			wantErr: false,
		},
		{
			desc:    "with path",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: "",
			arg:     filepath.Join("..", "system", "ls"),
			wants:   filepath.Join("..", "system", "ls"+winonly(".exe")),
			wantErr: false,
		},
		{
			desc:    "with path+extension",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: "",
			arg:     filepath.Join("..", "system", "ls.bat"),
			wants:   filepath.Join("..", "system", "ls.bat"),
			wantErr: false,
		},
		{
			desc:    "no extension, PATHEXT",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: ".com;.bat",
			arg:     "ls",
			wants:   filepath.Join(root, "_fixtures", "system", "ls"+winonly(".bat")),
			wantErr: false,
		},
		{
			desc:    "with extension, PATHEXT",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: ".com;.bat",
			arg:     "ls.exe",
			wants:   filepath.Join(root, "_fixtures", "system", "ls.exe"),
			wantErr: false,
		},
		{
			desc: "no extension, not found",
			path: []string{
				filepath.Join(root, "_fixtures", "nonexist"),
				filepath.Join(root, "_fixtures", "system"),
			},
			pathext: "",
			arg:     "cat",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "with extension, not found",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: "",
			arg:     "cat.exe",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "no extension, PATHEXT, not found",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: ".com;.bat",
			arg:     "cat",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "with extension, PATHEXT, not found",
			path:    []string{filepath.Join(root, "_fixtures", "system")},
			pathext: ".com;.bat",
			arg:     "cat.exe",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "relative path",
			path:    []string{filepath.Join("..", "system")},
			pathext: "",
			arg:     "ls",
			wants:   filepath.Join("..", "system", "ls"+winonly(".exe")),
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			setenv(t, "PATH", strings.Join(tC.path, string(filepath.ListSeparator)))
			setenv(t, "PATHEXT", tC.pathext)
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

func setenv(t *testing.T, name, newValue string) {
	oldValue, hasOldValue := os.LookupEnv(name)
	if err := os.Setenv(name, newValue); err != nil {
		t.Errorf("error setting environment variable %s: %v", name, err)
	}
	t.Cleanup(func() {
		if hasOldValue {
			_ = os.Setenv(name, oldValue)
		} else {
			_ = os.Unsetenv(name)
		}
	})
}

func winonly(s string) string {
	if runtime.GOOS == "windows" {
		return s
	}
	return ""
}
