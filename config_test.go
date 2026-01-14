package tagpr

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Songmu/gitconfig"
)

func TestConfig(t *testing.T) {
	tmpdir := t.TempDir()
	confPath := filepath.Join(tmpdir, defaultConfigFile)
	cfg := &config{
		conf:      confPath,
		gitconfig: &gitconfig.Config{GitPath: "git", File: confPath},
	}

	if err := cfg.Reload(); err != nil {
		t.Error(err)
	}
	if e, g := "", cfg.ReleaseBranch(); e != g {
		t.Errorf("got: %s, expext: %s", g, e)
	}
	if err := cfg.SetReleaseBranch("main"); err != nil {
		t.Error(err)
	}
	if e, g := "main", cfg.ReleaseBranch(); e != g {
		t.Errorf("got: %s, expext: %s", g, e)
	}
	if err := cfg.SetVersionFile(""); err != nil {
		t.Error(err)
	}
	if e, g := "", cfg.VersionFile(); e != g {
		t.Errorf("got: %s, expext: %s", g, e)
	}
	if e, g := []string{"major"}, cfg.MajorLabels(); !reflect.DeepEqual(e, g) {
		t.Errorf("got: %s, expext: %s", g, e)
	}
	if e, g := []string{"minor"}, cfg.MinorLabels(); !reflect.DeepEqual(e, g) {
		t.Errorf("got: %s, expext: %s", g, e)
	}
	if e, g := "[tagpr]", cfg.CommitPrefix(); !reflect.DeepEqual(e, g) {
		t.Errorf("got: %s, expext: %s", g, e)
	}

	b, err := os.ReadFile(confPath)
	if err != nil {
		t.Error(err)
	}

	var out string
	for line := range strings.SplitSeq(string(b), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			out += line + "\n"
		}
	}
	expect := `[tagpr]
	releaseBranch = main
	versionFile = -
`
	if out != expect {
		t.Errorf("got:\n%s\nexpect:\n%s", out, expect)
	}
}

func TestFixMajorVersion(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     *uint64
	}{
		{
			name:     "not set",
			envValue: "",
			want:     nil,
		},
		{
			name:     "numeric format",
			envValue: "1",
			want:     func() *uint64 { v := uint64(1); return &v }(),
		},
		{
			name:     "v-prefixed format",
			envValue: "v2",
			want:     func() *uint64 { v := uint64(2); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			confPath := filepath.Join(tmpdir, defaultConfigFile)
			cfg := &config{
				conf:      confPath,
				gitconfig: &gitconfig.Config{GitPath: "git", File: confPath},
			}

			if tt.envValue != "" {
				t.Setenv(envFixMajorVersion, tt.envValue)
			}

			if err := cfg.Reload(); err != nil {
				t.Fatal(err)
			}

			got := cfg.FixMajorVersion()
			if tt.want == nil {
				if got != nil {
					t.Errorf("got %v, want nil", *got)
				}
			} else {
				if got == nil {
					t.Errorf("got nil, want %v", *tt.want)
				} else if *got != *tt.want {
					t.Errorf("got %v, want %v", *got, *tt.want)
				}
			}
		})
	}
}
