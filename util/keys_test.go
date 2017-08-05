package util

import "testing"

func TestRandKey(t *testing.T) {
	key1 := RandKey(5)
	if len(key1) != 5 {
		t.Errorf("Expected size: %d, actual: %d", 5, len(key1))
	}
	key2 := RandKey(5)
	if len(key2) != 5 {
		t.Errorf("Expected size: %d, actual: %d", 5, len(key2))
	}
	if key1 == key2 {
		t.Errorf("Collision: %s - %s", key1, key2)
	}
	key3 := RandKey(10)
	if len(key3) != 10 {
		t.Errorf("Expected size: %d, actual: %d", 10, len(key3))
	}
}

