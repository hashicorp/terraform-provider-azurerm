package tf

import "testing"

func TestAccRandTimeInt(t *testing.T) {
	t.Run("Rand Date int", func(t *testing.T) {
		ri := AccRandTimeInt()

		if ri < 100000000000000000 {
			t.Fatalf("AccRandTimeInt returned a value (%d) shorter then expected", ri)
		}

		if ri > 999999999999999999 {
			t.Fatalf("AccRandTimeInt returned a value (%d) longer then expected", ri)
		}
	})
}
