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

func Test_majorMinor(t *testing.T) {
    type args struct {
        rdev uint64
    }
    tests := []struct {
        name  string
        args  args
        want  uint64
        want1 uint64
    }{
        {
            name:  "Device number 0",
            args:  args{rdev: 0},
            want:  0,
            want1: 0,
        },
        {
            name:  "Device number 8",
            args:  args{rdev: 8},
            want:  0,
            want1: 8,
        },
        {
            name:  "Device number 513",
            args:  args{rdev: 513},
            want:  2,
            want1: 1,
        },
        {
            name:  "Device number 66305",
            args:  args{rdev: 66305},
            want:  259,
            want1: 1,
        },
        {
            name:  "Large device number",
            args:  args{rdev: 18446744073709551615},
            want:  4095,
            want1: 1048575,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1 := majorMinor(tt.args.rdev)
            if got != tt.want {
                t.Errorf("majorMinor() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("majorMinor() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}
