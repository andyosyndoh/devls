package internals

import (
	"testing"
)

func Test_shouldShowFile(t *testing.T) {
	type args struct {
		name       string
		showHidden bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Regular file",
			args: args{name: "file.txt", showHidden: false},
			want: true,
		},
		{
			name: "Hidden file, show hidden false",
			args: args{name: ".hiddenfile", showHidden: false},
			want: false,
		},
		{
			name: "Hidden file, show hidden true",
			args: args{name: ".hiddenfile", showHidden: true},
			want: true,
		},
		{
			name: "Regular file with dot, not hidden",
			args: args{name: "file.with.dots.txt", showHidden: false},
			want: true,
		},
		{
			name: "Empty filename",
			args: args{name: "", showHidden: false},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShowFile(tt.args.name, tt.args.showHidden); got != tt.want {
				t.Errorf("shouldShowFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_joinPath(t *testing.T) {
	type args struct {
		elem []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Basic case with single element",
			args: args{elem: []string{"folder"}},
			want: "folder",
		},
		{
			name: "Multiple elements with no leading/trailing slashes",
			args: args{elem: []string{"folder", "subfolder", "file.txt"}},
			want: "folder/subfolder/file.txt",
		},
		{
			name: "Handles elements with leading slashes",
			args: args{elem: []string{"/folder", "subfolder", "file.txt"}},
			want: "/folder/subfolder/file.txt",
		},
		{
			name: "Handles elements with trailing slashes",
			args: args{elem: []string{"folder/", "subfolder/", "file.txt"}},
			want: "folder//subfolder//file.txt",
		},
		{
			name: "Handles both leading and trailing slashes",
			args: args{elem: []string{"/folder/", "/subfolder/", "file.txt"}},
			want: "/folder///subfolder//file.txt",
		},
		{
			name: "Handles empty elements",
			args: args{elem: []string{"folder", "", "file.txt"}},
			want: "folder//file.txt",
		},
		{
			name: "Handles all empty elements",
			args: args{elem: []string{"", "", ""}},
			want: "//",
		},
		{
			name: "Handles a mix of slashes and empty elements",
			args: args{elem: []string{"/", "", "file.txt"}},
			want: "///file.txt",
		},
		{
			name: "Handles root paths correctly",
			args: args{elem: []string{"/", "file.txt"}},
			want: "//file.txt",
		},
		{
			name: "Handles only slashes",
			args: args{elem: []string{"/", "/"}},
			want: "///",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPath(tt.args.elem...); got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
