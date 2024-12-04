package internals

import "testing"

func Test_cleanPath(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "Single slash remains unchanged",
			arg:  "/home/user/docs",
			want: "/home/user/docs",
		},
		{
			name: "Multiple consecutive slashes in middle",
			arg:  "/home//user/docs",
			want: "/home/user/docs",
		},
		{
			name: "Multiple consecutive slashes at start",
			arg:  "///home/user/docs",
			want: "/home/user/docs",
		},
		{
			name: "Multiple consecutive slashes at end",
			arg:  "/home/user/docs//",
			want: "/home/user/docs/",
		},
		{
			name: "Multiple consecutive slashes throughout",
			arg:  "////home///user//docs//",
			want: "/home/user/docs/",
		},
		{
			name: "Path with no slashes",
			arg:  "homeuserdocs",
			want: "homeuserdocs",
		},
		{
			name: "Empty path",
			arg:  "",
			want: "",
		},
		{
			name: "Root path only",
			arg:  "/",
			want: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanPath(tt.arg); got != tt.want {
				t.Errorf("cleanPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

