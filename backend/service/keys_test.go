package service

import "testing"

func TestKeysCollision(t *testing.T) {
	iterations := 10000
	keylen := 5
	keyset := make(map[string]struct{}, iterations)
	for i := 0; i < iterations; i++ {
		key := randKey(keylen)
		_, exists := keyset[key]
		if exists {
			t.Fatalf("Collision for %v iterations", iterations)
		}
		keyset[key] = struct{}{}
	}
}
