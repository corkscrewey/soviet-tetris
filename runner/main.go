package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	dockerURL = "https://www.docker.com/products/docker-desktop/"
	imageName = "tetris-simh"
	docoHash  = "021160476b9249d01d00a4368aab98872749fe43acab4159da0b0c3c083ab81a"

	filesDir = "files"
)

type action func(ctx context.Context) error

type pipeline []action

var workdir = flag.String("workdir", "workdir", "working directory for the emulator")

func main() {
	flag.Parse()
	if err := run(context.Background(), *workdir); err != nil {
		log.Fatal(err)
	}
}

// edges represents the dependency graph
var edges = map[string]string{
	"docker-compose.yaml": "",
	"docker":              "docker-compose.yaml",
	"docker-version":      "docker",
	"docker-image":        "docker-version",
	"sdl2":                "docker-compose.yaml",
	"mame":                "sdl2",
}

func run(ctx context.Context, workdir string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// check if docker-compose is valid
	docoFile := filepath.Join(wd, "docker-compose.yaml")
	gotHash, err := filehash(docoFile)
	if err != nil {
		return err
	}
	if gotHash != docoHash {
		return errors.New("docker-compose.yml is invalid, please download it again")
	}
	log.Printf("docker-compose.yml hash: %s", gotHash)

	// check docker
	docker, err := NewFromPath()
	if err != nil {
		return errors.New("docker is not installed, install it first: " + dockerURL)
	}
	log.Printf("docker found: %s", docker.exe)

	if docker.Version() < "19.03.0" {
		return errors.New("docker version is too old, please upgrade to 19.03.0 or newer")
	}
	log.Printf("docker version: %s", docker.Version())

	var p pipeline

	// check if image exists, if not, add build action
	if !docker.ImageExists(imageName) {
		log.Printf("building image %s", imageName)
		if err := docker.Compose().Build(ctx, docoFile, "simh"); err != nil {
			return fmt.Errorf("error building emulator image: %w", err)
		}
		log.Printf("image built: %s", imageName)
	} else {
		log.Printf("image found: %s", imageName)
	}

	// check mame if not exist, add download action
	if !mameExists(workdir) {
		log.Printf("will install mame to %s", workdir)
		p = append(p, installMame(filepath.Join(filesDir, "mame"), workdir))
	}

	// run emulation
	p = append(p, runEmulation("bin", "tetris"))

	return p.run(ctx)
}

func (p pipeline) run(ctx context.Context) error {
	for _, a := range p {
		if err := a(ctx); err != nil {
			return err
		}
	}
	return nil
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

func runEmulation(bin, game string) action {
	return func(ctx context.Context) error {
		log.Println("running emulation")
		return nil
	}
	// return func(ctx context.Context) error {
	// 	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", "docker-compose.yml", "up")
	// 	cmd.Stdout = nil
	// 	cmd.Stderr = nil
	// 	return cmd.Run()
	// }
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
