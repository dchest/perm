package perm

import (
	"math"
	"slices"
	"testing"
)

func TestPRF(t *testing.T) {
	t.Logf("%d\n", prf(0, 0))
	t.Logf("%d\n", prf(0, 1))
	t.Logf("%d\n", prf(0, 2))
	t.Logf("%d\n", prf(0, 3))
	if a, b := prf(1, 2), prf(2, 2); a == b {
		t.Fatalf("prf returned the same result: %d", a)
	}
}

func TestPerm(t *testing.T) {
	n := 10
	s := make([]int, 0, n)
	p := New(n, 12345)
	for i := range n {
		x := p.At(i)
		s = append(s, x)
	}
	t.Logf("%v - ", s)
	// Check inverse
	for i, x := range s {
		y, ok := p.Index(x)
		if !ok {
			t.Fatalf("failed to get inverse")
		}
		if i != y {
			t.Fatalf("wrong inverse: expected %d, got %d", i, y)
		}
	}
	// Check that all numbers are there
	slices.Sort(s)
	for i, v := range s {
		if i != v {
			t.Fatalf("wrong value: expected %d, got %d", i, v)
		}
	}
	t.Logf("ok\n")
}

func TestLenOne(t *testing.T) {
	p := New(1, 0)
	if x := p.At(0); x != 0 {
		t.Errorf("At(0) is %d, should be 0", x)
	}
	x, ok := p.Index(0)
	if !ok {
		t.Errorf("Inv(0) returned not ok")
	}
	if x != 0 {
		t.Errorf("Inv(0) is %d, should be 0", x)
	}
}

func TestBig(t *testing.T) {
	n := math.MaxInt
	p := New(n, 219381293812938)
	for i := range 10 {
		t.Logf("%d", p.At(i))
	}
	t.Logf("%d", p.At(math.MaxInt-1))
}

func TestSeed(t *testing.T) {
	n := 256
	p1 := New(n, 0)
	p2 := New(n, 1)
	numSame := 0
	for i := range n {
		if p1.At(i) == p2.At(i) {
			numSame++
		}
	}
	if numSame > n/2 {
		t.Fatalf("different seeds generated similar sequences: %d of %d are the same", numSame, n)
	}
}

func BenchmarkPRF(b *testing.B) {
	for b.Loop() {
		prf(10, 0x23456789)
	}
}

func BenchmarkFeistel(b *testing.B) {
	p := New(100, 12345)
	rk, mask, halfK := p.rk, p.mask, p.halfK
	for b.Loop() {
		feistelEnc(10, halfK, mask, &rk)
	}
}

func BenchmarkPerm1(b *testing.B) {
	n := 10
	p := New(n, 1234)
	for b.Loop() {
		p.At(1)
	}
}

func BenchmarkPerm2(b *testing.B) {
	n := 1000000
	p := New(n, 1234)
	for b.Loop() {
		p.At(53232)
	}
}

func BenchmarkPerm3(b *testing.B) {
	n := 17 // requires ~4 cycle-walks
	p := New(n, 0)
	for b.Loop() {
		p.At(0)
	}
}

func BenchmarkNewPerm(b *testing.B) {
	for b.Loop() {
		New(math.MaxInt-1, 0).At(0)
	}
}

func BenchmarkIndex(b *testing.B) {
	n := 1000000
	p := New(n, 1234)
	for b.Loop() {
		p.Index(213132)
	}
}
