package internals

import (
	"os"
	"testing"
)

func Test_GetFileColor(t *testing.T) {
	// Mock lsColors map and helper functions for testing
	lsColors = map[string]string{
		"tw": "01;34;42",
		"ow": "01;34;43",
		"di": "01;34",
		"ln": "01;36",
		"pi": "40;33",
		"so": "01;35",
		"bd": "40;33;01",
		"cd": "40;33;01",
		"su": "37;41",
		"sg": "30;43",
		"ex": "01;32",
	}

	tests := []struct {
		name     string
		mode     os.FileMode
		fileName string
		want     string
	}{
		{
			name:     "Directory with sticky and writable by others",
			mode:     os.ModeDir | 0o002 | 0o010,
			fileName: "",
			want:     "\033[01;34;42m",
		},
		{
			name:     "Directory writable by others",
			mode:     os.ModeDir | 0o002,
			fileName: "",
			want:     "\033[01;34;43m",
		},
		{
			name:     "Normal directory",
			mode:     os.ModeDir,
			fileName: "",
			want:     "\033[01;34m",
		},
		{
			name:     "Symlink",
			mode:     os.ModeSymlink,
			fileName: "",
			want:     "\033[01;36m",
		},
		{
			name:     "Named pipe",
			mode:     os.ModeNamedPipe,
			fileName: "",
			want:     "\033[40;33m",
		},
		{
			name:     "Socket",
			mode:     os.ModeSocket,
			fileName: "",
			want:     "\033[01;35m",
		},
		{
			name:     "Block device",
			mode:     os.ModeDevice,
			fileName: "",
			want:     "\033[40;33;01m",
		},
		{
			name:     "Character device",
			mode:     os.ModeCharDevice,
			fileName: "",
			want:     "\033[40;33;01m",
		},
		{
			name:     "Setuid",
			mode:     os.ModeSetuid,
			fileName: "",
			want:     "\033[37;41m",
		},
		{
			name:     "Setgid",
			mode:     os.ModeSetgid,
			fileName: "",
			want:     "\033[30;43m",
		},
		{
			name:     "Executable file",
			mode:     0o111,
			fileName: "",
			want:     "\033[01;32m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileColor(tt.mode, tt.fileName); got != tt.want {
				t.Errorf("GetFileColor() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_getFileExtension(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "File with single extension",
			arg:  "example.txt",
			want: "txt",
		},
		{
			name: "File with multiple extensions",
			arg:  "archive.tar.gz",
			want: "gz",
		},
		{
			name: "Hidden file with extension",
			arg:  ".env",
			want: "env",
		},
		{
			name: "File with no extension",
			arg:  "README",
			want: "",
		},
		{
			name: "Hidden file with no extension",
			arg:  ".hiddenfile",
			want: "hiddenfile",
		},
		{
			name: "Empty string",
			arg:  "",
			want: "",
		},
		{
			name: "File with trailing dot",
			arg:  "file.",
			want: "",
		},
		{
			name: "File with multiple dots and no final extension",
			arg:  "file.name.",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileExtension(tt.arg); got != tt.want {
				t.Errorf("getFileExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
