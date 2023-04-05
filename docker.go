package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Docker struct {
	exe     string
	version string
}

type Compose struct {
	exe     string
	basedir string
}

func NewFromPath() (*Docker, error) {
	exe := "docker"
	if runtime.GOOS == "windows" {
		exe = "docker.exe"
	}
	s, err := exec.LookPath(exe)
	if err != nil {
		return nil, fmt.Errorf("docker executable not found, make sure it's installed: %w", err)
	}
	version, err := dockerver(s)
	if err != nil {
		return nil, fmt.Errorf("error checking docker version: %w, make sure it's running", err)
	}
	return &Docker{exe: s, version: version}, nil
}

func dockerver(exe string) (string, error) {
	out, err := exec.Command(exe, "version", "--format", "{{.Server.Version}}").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (d Docker) Version() string {
	return d.version
}

func (d Docker) ImageExists(name string) bool {
	out, err := exec.Command(d.exe, "images", "--format", "{{.Repository}}").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		if line == name {
			return true
		}
	}
	return false
}

func (d Docker) Compose(basedir string) *Compose {
	return &Compose{exe: d.exe, basedir: basedir}
}

func (d Compose) Build(ctx context.Context, file, service string) error {
	cmd := exec.CommandContext(ctx, d.exe, "compose", "-f", filepath.Join(d.basedir, file), "build", service)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d Compose) Up(ctx context.Context, file, service string) error {
	cmd := exec.CommandContext(ctx, d.exe, "compose", "-f", filepath.Join(d.basedir, file), "up", service)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func initTetrisDocker(ctx context.Context, workdir string, tetrisDir string) (*Compose, error) {
	docker, err := NewFromPath()
	if err != nil {
		return nil, err
	}
	log.Printf("docker found: %s", docker.exe)

	if docker.Version() < dockerVer {
		return nil, fmt.Errorf("docker version is too old, please upgrade to %s or newer", dockerVer)
	}
	log.Printf("docker version: %s", docker.Version())

	// check if image exists, if not, add build action
	imageName := dirTetris + "-simh"
	composeBase := filepath.Join(workdir, tetrisDir)
	if !docker.ImageExists(imageName) {
		log.Printf("building image %s", imageName)
		if err := docker.Compose(composeBase).Build(ctx, "docker-compose.yaml", "simh"); err != nil {
			return nil, fmt.Errorf("error building emulator image: %w", err)
		}
		log.Printf("image built: %s", imageName)
	} else {
		log.Printf("image found: %s", imageName)
	}

	return docker.Compose(composeBase), nil
}
