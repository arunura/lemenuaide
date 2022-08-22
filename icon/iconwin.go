//go:build windows

package icon

import (
	_ "embed"
)

//go:embed lemonade.ico
var LEMONADE []byte
