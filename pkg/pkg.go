package pkg

import (
	"os"
	"path"
)

var (
	// RootDir is the asset root directory.
	RootDir = path.Join(os.Getenv("HOME"), ".trustacks")

	// BinDir is the binary dependencies directory.
	BinDir = path.Join(RootDir, "bin")
)
