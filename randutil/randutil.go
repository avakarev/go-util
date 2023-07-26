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

// String generates secure random string with given length
func String(n int) (string, error) {
	ret := make([]byte, n)
	length := len(alphabet)
	for i := 0; i < n; i++ {
		ret[i] = alphabet[Int(length)]
	}
	return string(ret), nil
}

// Int returns random number in range [0, max]
// It will panic if input is invalid, <= 0
func Int(max int) int {
	rnd.Lock()
	defer rnd.Unlock()
	return rnd.r.Intn(max)
}
