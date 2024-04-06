//go:build !tinygo

package temp

import "embed"

//go:embed css go.mod html images js template
var fs embed.FS
