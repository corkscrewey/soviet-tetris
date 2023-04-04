package main

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func installMame(distdir string, workdir string) action {
	mameZIP := "mame 252.zip"
	if runtime.GOARCH == "amd64" {
		mameZIP = "mame0252-arm64.zip"
	}
	return func(ctx context.Context) error {
		z, err := zip.OpenReader(filepath.Join(distdir, mameZIP))
		if err != nil {
			return err
		}
		for _, f := range z.File {
			if f.Name == "mame" {
				rc, err := f.Open()
				if err != nil {
					return err
				}
				defer rc.Close()
				f, err := os.Create(filepath.Join(workdir, "mame"))
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err := io.Copy(f, rc); err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
}
