package pkg

import (
	"os"
	"path/filepath"
)

var (
	// RootDir is the asset root directory.
	RootDir = func() string {
		if os.Getenv("DATA_DIR") != "" {
			return os.Getenv("DATA_DIR")
		}
		return filepath.Join(os.Getenv("HOME"), ".tsconfig")
	}()

	// BinDir is the binary dependencies directory.
	BinDir = filepath.Join(RootDir, "bin")
)
