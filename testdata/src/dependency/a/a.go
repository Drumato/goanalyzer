package a // want "analyze finished"

import "dependency/b"

func f() {
	var _ = b.Zero()
}

