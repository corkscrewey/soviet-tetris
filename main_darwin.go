package main

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func installMame(mameExe string) error {
	mameZIP := "https://github.com/corkscrewey/soviet-tetris/releases/download/files/mame_252.zip"
	if runtime.GOARCH == "arm64" {
		mameZIP = "https://github.com/corkscrewey/soviet-tetris/releases/download/files/mame0252-arm64.zip"
	}

	if err := os.MkdirAll(filepath.Dir(mameExe), 0755); err != nil {
		return err
	}

	tmpf, err := os.CreateTemp("", "mame.zip")
	if err != nil {
		return err
	}
	defer tmpf.Close()
	defer os.Remove(tmpf.Name())

	resp, err := http.Get(mameZIP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	size, err := io.Copy(tmpf, resp.Body)
	if err != nil {
		return err
	}

	z, err := zip.NewReader(tmpf, size)
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
			f, err := os.Create(filepath.Join(mameExe, "mame"))
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(f, rc); err != nil {
				return err
			}
			if err := f.Chmod(0755); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("mame binary not found in zip")
}

func mameExists(workdir string) bool {
	if _, err := os.Stat(filepath.Join(workdir, mameExe())); err != nil {
		return false
	}
	return true
}
