package main

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	//go:embed Dockerfile
	dockerfile []byte
	//go:embed docker-compose.yaml
	dockercompose []byte
	//go:embed files/rom/ie15.zip
	ie15 []byte
	//go:embed files/rom/ie15kbd.zip
	ie15kbd []byte
)

const (
	dirTetris = "tetris"
	dirROM    = "roms"
)

var dirstructure = []struct {
	Dir      string
	Name     string
	Contents *[]byte
	Mode     os.FileMode
}{
	{dirTetris, "Dockerfile", &dockerfile, 0644},
	{dirTetris, "docker-compose.yaml", &dockercompose, 0644},
	{dirROM, "ie15.zip", &ie15, 0644},
	{dirROM, "ie15kbd.zip", &ie15kbd, 0644},
}

func bootstrap(dir string) error {
	for _, d := range dirstructure {
		base := filepath.Join(dir, d.Dir)
		if err := os.MkdirAll(base, 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(base, d.Name), *d.Contents, d.Mode); err != nil {
			return err
		}
	}
	return nil
}

func filehash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return readerhash(f)
}

func readerhash(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func dirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("%s is not a directory", path)
	}
	return true, nil
}
