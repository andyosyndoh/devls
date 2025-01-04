package internals

import "testing"

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
