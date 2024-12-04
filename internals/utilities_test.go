package internals

import "testing"

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
