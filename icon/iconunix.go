//go:build linux || darwin

package icon

import (
	_ "embed"
)

//go:embed lemonade.png
var LEMONADE []byte
