package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embedFS embed.FS

// DistFS returns the embedded filesystem for the ui/dist directory
func DistFS() (fs.FS, error) {
	return fs.Sub(embedFS, "dist")
}
