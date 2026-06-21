package perm

import (
	"fmt"
	"math"
	"slices"
	"strings"
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

func TestExample(t *testing.T) {
	p := New(3, 42)
	a := [3]int{p.At(0), p.At(1), p.At(2)}
	b := [3]int{2, 0, 1}
	if a != b {
		t.Logf("[0, 1, 2] -> expected %v, got %v", b, a)
	}
}

func TestPermSmall(t *testing.T) {
	n := 4
	s := make([]int, n)
	p := New(n, 0)
	for i := range s {
		s[i] = p.At(i)
		if idx, ok := p.Index(s[i]); !ok || idx != i {
			t.Fatalf("wrong index at %d", i)
		}
	}
	t.Logf("small: %v", s)
	slices.Sort(s)
	for i, v := range s {
		if v != i {
			t.Fatalf("wrong result at %d", i)
		}
	}
}

func TestPerm(t *testing.T) {
	n := 1000
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

func TestOddK(t *testing.T) {
	lengths := []int{5, 6, 7, 8, 17, 31, 32, 33, 64, 127, 128}
	seeds := []uint64{0, 1, 42, math.MaxUint64}
	for _, n := range lengths {
		for _, seed := range seeds {
			p := New(n, seed)
			for i := range n {
				v := p.At(i)
				if v < 0 || v >= n {
					t.Fatalf("n=%d seed=%d: At(%d)=%d out of range", n, seed, i, v)
				}
				idx, ok := p.Index(v)
				if !ok {
					t.Fatalf("n=%d seed=%d: Index(%d) returned false", n, seed, v)
				}
				if idx != i {
					t.Fatalf("n=%d seed=%d: At(%d)=%d, Index(%d)=%d (expected %d)", n, seed, i, v, v, idx, i)
				}
			}
		}
	}
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

func TestPermVisual(t *testing.T) {
	n := 64 * 64
	p := New(n, 1)
	t.Logf("Visual halves")
	for y := 1; y < 64; y++ {
		var line strings.Builder
		for x := 1; x < 64; x++ {
			if p.At(x*y) < n/2 {
				line.WriteString(".")
			} else {
				line.WriteString("o")
			}
		}
		t.Log(line.String())
	}
}

func TestPermVisualBar(t *testing.T) {
	n := 64
	p := New(n, 1)
	t.Logf("Visual bars")
	for i := range n {
		var line strings.Builder
		fmt.Fprintf(&line, "%3d ", i)
		v := p.At(i)
		for range v {
			line.WriteString("█")
		}
		t.Log(line.String())
	}
}

func BenchmarkPRF(b *testing.B) {
	for b.Loop() {
		prf(10, 0x23456789)
	}
}

func BenchmarkFeistel(b *testing.B) {
	p := New(100, 12345)
	rk, maskL, maskR, halfK := p.rk, p.lmask, p.rmask, p.rbits
	for b.Loop() {
		feistelEnc(10, halfK, maskL, maskR, &rk)
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
	n := 17
	p := New(n, 0)
	for b.Loop() {
		p.At(4)
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
