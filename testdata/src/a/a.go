package a // want "analyze finished"

import "b"

func f() {
	var _ = b.Zero()
}

