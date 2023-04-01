package static

import (
	"embed"
	_ "embed"
)

//go:embed dist/css.css
var Css []byte

//go:embed dist/index.html
var Index []byte

//go:embed dist/script.js
var Script []byte

//go:embed dist
var StatFiles embed.FS
