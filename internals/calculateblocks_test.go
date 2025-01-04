package internals

import (
	"testing"
)

func Test_max(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "a greater than b",
			args: args{a: 5, b: 3},
			want: 5,
		},
		{
			name: "b greater than a",
			args: args{a: 2, b: 7},
			want: 7,
		},
		{
			name: "a equal to b",
			args: args{a: 4, b: 4},
			want: 4,
		},
		{
			name: "negative numbers",
			args: args{a: -3, b: -5},
			want: -3,
		},
		{
			name: "zero and positive",
			args: args{a: 0, b: 8},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := max(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHidden(t *testing.T) {
    type args struct {
        name string
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
        {
            name: "Hidden file",
            args: args{name: ".hiddenfile"},
            want: true,
        },
        {
            name: "Non-hidden file",
            args: args{name: "visiblefile.txt"},
            want: false,
        },
        {
            name: "File with dot in the middle",
            args: args{name: "file.with.dots.txt"},
            want: false,
        },
        {
            name: "Single dot",
            args: args{name: "."},
            want: true,
        },
        {
            name: "Double dot",
            args: args{name: ".."},
            want: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := isHidden(tt.args.name); got != tt.want {
                t.Errorf("isHidden() = %v, want %v", got, tt.want)
            }
        })
    }
}
