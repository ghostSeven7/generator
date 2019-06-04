package generator

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New(8)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerator(t *testing.T) {
	node, err := New(8)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		node.generator()
	}
}

func Benchmark_Generator(b *testing.B) {
	node, err := New(8)
	if err != nil {
		for i := 0; i < b.N; i++ {
			node.generator()
		}
	}
}
