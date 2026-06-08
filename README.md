# perm

    import "github.com/dchest/perm"

Package perm implements a seekable permutation.

## Usage

#### type Perm

```go
type Perm struct {
}
```

Perm is a fast non-cryptographic seekable permutation.

It maps an integer in a range into an integer from the same range.

It is a virtual pseudorandom shuffle of all values 0, 1, 2, ... up to the given
length, but generates values at indexes dynamically instead of shuffling in
memory.

It can also be considered a pseudorandom number generator, except none of the
generated numbers repeat, e.g. given the length 3 indexes [0, 1, 2] will
generate [0, 1, 2], [0, 2, 1], [1, 0, 2], [1, 2, 0], [2, 0, 1], or [2, 1, 0],
depending on the seed.

    p := perm.New(3, 42)
    p.At(0) // -> 2
    p.At(1) // -> 0
    p.At(2) // -> 1

It can also be considered an insecure format-preserving encryption, where At is
encryption and Index is decryption in the domain of integers from 0 to length-1,
and in fact is based on a cycle-walked 4-round Feistel network with a weak PRF.

#### func  New

```go
func New(length int, seed uint64) Perm
```
New creates a seeded permutation of [0..length-1].

It panics if length is negative or zero.

#### func (Perm) At

```go
func (p Perm) At(index int) int
```
At generates the value of the permutation at the given index.

It panics if the given index is outside the range.

#### func (Perm) Index

```go
func (p Perm) Index(value int) (index int, ok bool)
```
Index generates the inverse of At.

In other terms, it "finds" the index of the given value.

If the given value is outside the range, it returns -1, false.
