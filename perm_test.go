package perm

import (
	"slices"
	"testing"
)

func TestPerm(t *testing.T) {
	t.Logf("%d\n", prf(0, 0))
	t.Logf("%d\n", prf(0, 1))
	t.Logf("%d\n", prf(0, 2))
	t.Logf("%d\n", prf(0, 3))
	n := uint32(10)
	s := make([]uint32, 0, n)
	p, err := New(n, 12345)
	if err != nil {
		panic(err)
	}
	for i := range n {
		x := p.At(i)
		s = append(s, x)
	}
	t.Logf("%v - ", s)
	// Check that all numbers are there
	slices.Sort(s)
	for i, v := range s {
		if uint32(i) != v {
			t.Fatalf("not ok: expected %d, got %d", i, v)
		}
	}
	t.Logf("ok\n")
}

func BenchmarkPRF(b *testing.B) {
	for b.Loop() {
		prf(10, 0x23456789)
	}
}

func BenchmarkFeistel(b *testing.B) {
	p, _ := New(100, 12345)
	seed, mask, halfK := p.seed, p.mask, p.halfK
	for b.Loop() {
		feistel(10, halfK, mask, seed)
	}
}

func BenchmarkPerm1(b *testing.B) {
	n := uint32(10)
	p, _ := New(n, 1234)
	for b.Loop() {
		p.At(1)
	}
}

func BenchmarkPerm2(b *testing.B) {
	n := uint32(1000000)
	p, _ := New(n, 1234)
	for b.Loop() {
		p.At(53232)
	}
}
