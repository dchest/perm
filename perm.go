// Package perm implements a seekable permutation.
package perm

import (
	"math/bits"
)

// Perm is a fast non-cryptographic seekable permutation.
//
// It maps an integer in a range into an integer from the same range.
//
// It is a virtual pseudorandom shuffle of all values 0, 1, 2, ... up to the
// given length, but generates values at indexes dynamically instead of
// shuffling in memory.
//
// It can also be considered a pseudorandom number generator, except none of
// the generated numbers repeat, e.g. given the length 3 indexes [0, 1, 2] will
// generate [0, 1, 2], [0, 2, 1], [1, 0, 2], [1, 2, 0], [2, 0, 1], or [2, 1, 0],
// depending on the seed.
//
//	p := perm.New(3, 42)
//	p.At(0) // -> 2
//	p.At(1) // -> 0
//	p.At(2) // -> 1
//
// It can also be considered an insecure format-preserving encryption, where
// At is encryption and Index is decryption in the domain of integers from 0 to
// length-1, and in fact is based on a cycle-walked 4-round Feistel network
// with a weak PRF.
type Perm struct {
	rk                  [4]uint32 // Feistel subkeys derived from seed
	rbits, lmask, rmask uint32    // helpers for forcing into the range
	maxi                uint      // maximum index (length-1)
}

// New creates a seeded permutation of [0..length-1].
//
// It panics if length is negative or zero.
func New(length int, seed uint64) Perm {
	if length <= 0 {
		panic("length cannot be negative or zero")
	}
	maxi := uint(length - 1)

	// Split bits into two halves (possibly differ by 1 bit).
	k := uint32(bits.Len(maxi))
	rbits := k / 2
	rmask := uint32((uint64(1) << rbits) - 1)
	lmask := uint32((uint64(1) << (k - rbits)) - 1)

	// SHA256/BLAKE IV xored into seed to generate round keys.
	// "first 32 bits of the fractional parts of the square
	// roots of the first primes"
	const (
		r0 = 0x6a09e667
		r1 = 0xbb67ae85
		r2 = 0x3c6ef372
		r3 = 0xa54ff53a
	)

	sl := uint32(seed)
	sh := uint32(seed >> 32)

	return Perm{
		maxi:  maxi,
		rbits: rbits,
		lmask: lmask,
		rmask: rmask,
		rk: [4]uint32{
			r0 ^ sl,
			r1 ^ sh,
			r2 ^ sl,
			r3 ^ sh,
		},
	}
}

// At generates the value of the permutation at the given index.
//
// It panics if the given index is outside the range.
func (p Perm) At(index int) int {
	if index < 0 || uint(index) > p.maxi {
		panic("Perm.At: index outside of the range")
	}
	v := uint(index)
	for {
		v = feistelEnc(v, p.rbits, p.lmask, p.rmask, &p.rk)
		// if in the range, return it, otherwise
		// cycle walk to force v into the range.
		if v <= p.maxi {
			return int(v)
		}
	}
}

// Index generates the inverse of At.
//
// In other terms, it "finds" the index of the given value.
//
// If the given value is outside the range, it returns -1, false.
func (p Perm) Index(value int) (index int, ok bool) {
	if value < 0 || uint(value) > p.maxi {
		return -1, false
	}
	idx := uint(value)
	for {
		idx = feistelDec(idx, p.rbits, p.lmask, p.rmask, &p.rk)
		if idx <= p.maxi {
			return int(idx), true
		}
	}
}

// feistelEnc encodes v with a 4-round Feistel network.
//
// It returns a permutation of v with the left half masked by lmask
// and the right half masked by rmask, mixed with the given round keys.
func feistelEnc(v uint, rbits, lmask, rmask uint32, rk *[4]uint32) uint {
	l, r := uint32(v>>rbits)&lmask, uint32(v)&rmask
	l, r = fr(l, r, rk[0], lmask)
	l, r = fr(l, r, rk[1], rmask)
	l, r = fr(l, r, rk[2], lmask)
	l, r = fr(l, r, rk[3], rmask)
	return uint(l)<<rbits | uint(r)
}

// feistelDec decodes v with a 4-round Feistel network.
func feistelDec(v uint, rbits, lmask, rmask uint32, rk *[4]uint32) uint {
	l, r := uint32(v>>rbits)&lmask, uint32(v)&rmask
	r, l = fr(r, l, rk[3], rmask)
	r, l = fr(r, l, rk[2], lmask)
	r, l = fr(r, l, rk[1], rmask)
	r, l = fr(r, l, rk[0], lmask)
	return uint(l)<<rbits | uint(r)
}

// fr applies a Feistel round.
func fr(l, r, rk, mask uint32) (L, R uint32) {
	return r, l ^ (prf(r, rk) & mask)
}

// prf mixes bits of v, MurmurHash/XXH32-style.
// Fixed point prf(0, 0)=0 doesn't matter.
// Also prf(x, y)==prf(y, x).
// It is reversible, thus in fact a PRP.
func prf(v, k uint32) uint32 {
	// XXH32 constants, chosen empirically
	const (
		prime2 = 0x85ebca77
		prime3 = 0xc2b2ae3d
	)
	h := v ^ k
	h ^= h >> 15
	h *= prime2
	h ^= h >> 13
	h *= prime3
	h ^= h >> 16
	return h
}
