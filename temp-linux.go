//go:build !tinygo

package temp

import "embed"

//go:embed css js template
var fs embed.FS
