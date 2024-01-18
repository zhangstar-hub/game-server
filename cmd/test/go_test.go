package main

import (
	"strings"
	"testing"
)

func BenchmarkAddString(b *testing.B) {
	a := strings.Repeat("a", 100)
	c := strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		AddString(a, c)
	}
}

func BenchmarkBuildString(b *testing.B) {
	a := strings.Repeat("a", 100)
	c := strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		BuildString(a, c)
	}
}

func BenchmarkSprint(b *testing.B) {
	a := strings.Repeat("a", 100)
	c := strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		SprintString(a, c)
	}
}
