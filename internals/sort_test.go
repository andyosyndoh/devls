package internals

import "testing"

func Test_compareStrings(t *testing.T) {
    type args struct {
        a string
        b string
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
        {
            name: "Equal strings",
            args: args{a: "hello", b: "hello"},
            want: false,
        },
        {
            name: "a < b",
            args: args{a: "abc", b: "def"},
            want: true,
        },
        {
            name: "a > b",
            args: args{a: "xyz", b: "abc"},
            want: false,
        },
        {
            name: "Empty strings",
            args: args{a: "", b: ""},
            want: false,
        },
        {
            name: "Case sensitivity",
            args: args{a: "ABC", b: "abc"},
            want: true,
        },
        {
            name: "Numbers",
            args: args{a: "123", b: "456"},
            want: true,
        },
        {
            name: "Mixed alphanumeric",
            args: args{a: "a1b2", b: "a1c3"},
            want: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := compareStrings(tt.args.a, tt.args.b); got != tt.want {
                t.Errorf("compareStrings() = %v, want %v", got, tt.want)
            }
        })
    }
}
