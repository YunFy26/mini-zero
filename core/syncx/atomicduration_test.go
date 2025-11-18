package syncx

import (
	"testing"
	"time"
)

func TestAtomicDuration(t *testing.T) {
	t.Run("NewAtomicDuration", func(t *testing.T) {
		ad := NewAtomicDuration()
		if val := ad.Load(); val != 0 {
			t.Fatalf("Expected 0, got %v", val)
		}
	})

	t.Run("ForAtomicDuration", func(t *testing.T) {
		dur := time.Second
		ad := ForAtomicDuration(dur)
		if val := ad.Load(); val != dur {
			t.Fatalf("Expected %v, got %v", dur, val)
		}
	})

	t.Run("Set and Load", func(t *testing.T) {
		ad := NewAtomicDuration()
		dur := time.Minute
		ad.Set(dur)
		if val := ad.Load(); val != dur {
			t.Fatalf("Expected %v, got %v", dur, val)
		}
	})

	t.Run("CompareAndSwap", func(t *testing.T) {
		ad := ForAtomicDuration(time.Hour)

		// Successful swap
		swapped := ad.CompareAndSwap(time.Hour, time.Minute)
		if !swapped {
			t.Fatal("Expected swap to succeed")
		}
		if val := ad.Load(); val != time.Minute {
			t.Fatalf("Expected %v after swap, got %v", time.Minute, val)
		}

		// Failed swap
		swapped = ad.CompareAndSwap(time.Hour, time.Second)
		if swapped {
			t.Fatal("Expected swap to fail")
		}
		if val := ad.Load(); val != time.Minute {
			t.Fatalf("Expected value to remain %v, got %v", time.Minute, val)
		}
	})
}
