package cache

import (
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	Set("a", 1337)

	if val, _ := Get("a").Int(); val != 1337 {
		t.Error("Set not working")
	}
}

func TestExpire(t *testing.T) {
	Set("a", 1337)

	time.Sleep(5 * time.Second)

	if val, _ := Get("a").Int(); val != 1337 {
		t.Error("Expiration not working")
	}
}
