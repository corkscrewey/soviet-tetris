package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	dockerVer = "19.03.0"
)

var workdir = flag.String("workdir", "workdir", "working directory for the emulator")

func main() {
	flag.Parse()
	if err := run(context.Background(), *workdir); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, workdir string) error {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }

	if _, err := dirExists(workdir); err != nil {
		return err
	}
	if err := bootstrap(workdir); err != nil {
		return err
	}

	// check docker
	docker, err := NewFromPath()
	if err != nil {
		return err
	}
	log.Printf("docker found: %s", docker.exe)

	if docker.Version() < dockerVer {
		return fmt.Errorf("docker version is too old, please upgrade to %s or newer", dockerVer)
	}
	log.Printf("docker version: %s", docker.Version())

	// // check if image exists, if not, add build action
	// if !docker.ImageExists(imageName) {
	// 	log.Printf("building image %s", imageName)
	// 	if err := docker.Compose().Build(ctx, docoFile, "simh"); err != nil {
	// 		return fmt.Errorf("error building emulator image: %w", err)
	// 	}
	// 	log.Printf("image built: %s", imageName)
	// } else {
	// 	log.Printf("image found: %s", imageName)
	// }

	// // check mame if not exist, add download action
	// if !mameExists(workdir) {
	// 	log.Printf("will install mame to %s", workdir)
	// 	if err := installMame(filepath.Join(filesDir, "mame"), workdir); err != nil {
	// 		return err
	// 	}
	// }

	return runEmulation("bin", "tetris")

}

func mameExists(workdir string) bool {
	if _, err := os.Stat(filepath.Join(workdir, mameExe())); err != nil {
		return false
	}
	return true
}

func mameExe() string {
	if runtime.GOOS == "windows" {
		return "mame64.exe"
	}
	return "mame"
}

func runEmulation(bin, game string) error {
	log.Println("running emulation")
	return nil
}

func filehash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
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
