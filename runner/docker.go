package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Docker struct {
	exe     string
	version string
}

type Composer struct {
	exe string
}

func NewFromPath() (*Docker, error) {
	exe := "docker"
	if runtime.GOOS == "windows" {
		exe = "docker.exe"
	}
	s, err := exec.LookPath(exe)
	if err != nil {
		return nil, err
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

func (d Docker) Compose() Composer {
	return Composer{exe: d.exe}
}

func (d Composer) Build(ctx context.Context, file, service string) error {
	cmd := exec.CommandContext(ctx, d.exe, "compose", "-f", file, "build", service)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
