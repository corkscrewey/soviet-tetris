package main

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func installMame(mameExe string) error {
	if fi, err := os.Stat(mameExe); err == nil {
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			return nil
		}
	}

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

	log.Printf("fetching mame from %s", mameZIP)
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
			f, err := os.Create(mameExe)
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
