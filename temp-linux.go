//go:build !tinygo

package temp

import "embed"

//go:embed css go.mod js template
var fs embed.FS
