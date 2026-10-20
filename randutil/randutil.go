// Package randutil implements random number/string generator helpers
package randutil

import (
	"math/rand"
	"sync"
	"time"
)

var rnd = struct {
	sync.Mutex
	r *rand.Rand
}{ //#nosec G404 -- math/rand is enough
	r: rand.New(rand.NewSource(time.Now().UnixNano())),
}

// alphabet is a combination of chars that could be a member of random string
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// String generates a non-cryptographic random string with given length.
// It is backed by math/rand and must not be used for secrets such as tokens
// or passwords; use crypto/rand for those.
func String(n int) string {
	ret := make([]byte, n)
	length := len(alphabet)
	for i := range n {
		ret[i] = alphabet[Int(length)]
	}
	return string(ret)
}

// Int returns random number in range [0, nmax)
// It will panic if input is invalid, <= 0
func Int(nmax int) int {
	rnd.Lock()
	defer rnd.Unlock()
	return rnd.r.Intn(nmax)
}
