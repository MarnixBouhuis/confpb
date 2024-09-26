package testutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func CopyAllFromFS(t *testing.T, srcFS fs.FS, destinationPath string) {
	t.Helper()

	err := fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(destinationPath, path)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}
		return CopyFileFromFSToDisk(srcFS, path, destPath)
	})
	require.NoError(t, err)
}

func CopyFileFromFSToDisk(srcFS fs.FS, srcPath string, destPath string) error {
	src, err := srcFS.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to write file from fs to disk, open failed: %w", err)
	}
	defer func() {
		_ = src.Close()
	}()

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to write file from fs to disk, create failed: %w", err)
	}
	defer func() {
		_ = dest.Close()
	}()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("failed to write file from fs to disk, copy failed: %w", err)
	}
	return nil
}
