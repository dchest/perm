package perm

import (
	"errors"
	"math/bits"
)

// Perm is a non-cryptographic permutation.
type Perm struct{ vmax, seed, halfK, mask uint32 }

// New creates a seeded permutation of [0..length-1].
func New(length, seed uint32) (Perm, error) {
	if length <= 0 {
		return Perm{}, errors.New("length cannot be negative or zero")
	}
	vmax := length - 1
	// Calculate the smallest even bit-length that can hold vmax
	k := bits.Len32(vmax)
	if k%2 != 0 {
		k++
	}
	halfK := uint32(k / 2)
	mask := uint32(1<<halfK) - 1
	return Perm{vmax: vmax, seed: seed, halfK: halfK, mask: mask}, nil
}

// At generates the value of the permutation [0...length-1]
// at the given index.
//
// If i is outside the range, it panics.
func (p Perm) At(index uint32) uint32 {
	if index > p.vmax || p.vmax == 0 {
		panic("index outside of the range")
	}
	for {
		index = feistel(index, p.halfK, p.mask, p.seed)
		if index <= p.vmax {
			return index
		}
		// otherwise, cycle walk to force v into the range.
	}
}

// prf mixes bits of v, MurmurHash/XXH32-style.
// Fixed point prf(0, 0)=0 doesn't matter.
func prf(v, k uint32) uint32 {
	const (
		// XXH32 constants, chosen empirically
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

// feistel is a 4-round Feistel network.
//
// It returns a permutation of v limited to halfK*2 bits
// mixed with the given seed.
func feistel(v, halfK, mask, seed uint32) uint32 {
	const (
		// SHA256/BLAKE IV xored into seed to generate round keys
		// "first 32 bits of the fractional parts of the square roots of the first primes"
		r1 = 0x6a09e667
		r2 = 0xbb67ae85
		r3 = 0x3c6ef372
		r4 = 0xa54ff53a
	)
	l, r := (v>>halfK)&mask, v&mask
	l, r = r, l^(prf(r, r1^seed)&mask)
	l, r = r, l^(prf(r, r2^seed)&mask)
	l, r = r, l^(prf(r, r3^seed)&mask)
	l, r = r, l^(prf(r, r4^seed)&mask)
	return (l << halfK) | r
}
