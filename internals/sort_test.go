package internals

import (
	"reflect"
	"testing"
)

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

func TestSortStringsAscending(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty slice",
			args: args{slice: []string{}},
			want: []string{},
		},
		{
			name: "Single element",
			args: args{slice: []string{"a"}},
			want: []string{"a"},
		},
		{
			name: "Already sorted",
			args: args{slice: []string{"a", "b", "c"}},
			want: []string{"a", "b", "c"},
		},
		{
			name: "Reverse sorted",
			args: args{slice: []string{"c", "b", "a"}},
			want: []string{"a", "b", "c"},
		},
		{
			name: "Mixed case",
			args: args{slice: []string{"B", "a", "C"}},
			want: []string{"B", "C", "a"},
		},
		{
			name: "Numbers and letters",
			args: args{slice: []string{"3", "1", "2", "b", "a"}},
			want: []string{"1", "2", "3", "a", "b"},
		},
		{
			name: "Duplicate elements",
			args: args{slice: []string{"b", "a", "b", "c", "a"}},
			want: []string{"a", "a", "b", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortStringsAscending(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortStringsAscending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortStringsDescending(t *testing.T) {
    type args struct {
        slice []string
    }
    tests := []struct {
        name string
        args args
        want []string
    }{
        {
            name: "Empty slice",
            args: args{slice: []string{}},
            want: []string{},
        },
        {
            name: "Single element",
            args: args{slice: []string{"a"}},
            want: []string{"a"},
        },
        {
            name: "Already sorted descending",
            args: args{slice: []string{"c", "b", "a"}},
            want: []string{"c", "b", "a"},
        },
        {
            name: "Ascending sorted",
            args: args{slice: []string{"a", "b", "c"}},
            want: []string{"c", "b", "a"},
        },
        {
            name: "Mixed case",
            args: args{slice: []string{"B", "a", "C"}},
            want: []string{"a", "C", "B"},
        },
        {
            name: "Numbers and letters",
            args: args{slice: []string{"3", "1", "2", "b", "a"}},
            want: []string{"b", "a", "3", "2", "1"},
        },
        {
            name: "Duplicate elements",
            args: args{slice: []string{"b", "a", "b", "c", "a"}},
            want: []string{"c", "b", "b", "a", "a"},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := SortStringsDescending(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("SortStringsDescending() = %v, want %v", got, tt.want)
            }
        })
    }
}
