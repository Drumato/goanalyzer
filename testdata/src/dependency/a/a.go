package a // want "dependency analyze finished"

import "dependency/b"

func f() {
	var _ = b.Zero()
}
