package internal

import "embed"

var (
	//go:embed resources
	f embed.FS
)

func GetEmbeddedFS() embed.FS {
	return f
}
