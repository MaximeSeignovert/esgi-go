package main

import "testing"

func TestSumDivisors(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "one", n: 1, want: 1},
		{name: "six", n: 6, want: 12},
		{name: "prime", n: 7, want: 8},
		{name: "ten", n: 10, want: 18},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := sumDivisors(test.n)
			if got != test.want {
				t.Fatalf("sumDivisors(%d) = %d, want %d", test.n, got, test.want)
			}
		})
	}
}
