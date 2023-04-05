package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type MAME struct {
	exe    string
	romDir string
}

func NewMAMEFromPath(workdir string, romDir string) (*MAME, error) {
	exe := mameExe()
	localMameDir := filepath.Join(workdir, "bin")
	localROMDir := filepath.Join(workdir, romDir)
	exepath, err := exec.LookPath(exe)
	if err != nil {
		// not found in path, download and unpack.
		exepath = filepath.Join(localMameDir, "mame")
		if err := installMame(exepath); err != nil {
			return nil, fmt.Errorf("error installing mame: %w", err)
		}
		return nil, fmt.Errorf("mame executable not found, make sure it's installed: %w", err)
	}
	return &MAME{exe: exepath, romDir: localROMDir}, nil
}

func (m *MAME) Run() error {
	cmd := exec.Command(m.exe,
		"-rompath", m.romDir,
		"-window",
		"-resolution", "640x480",
		"-video", "opengl",
		"ie15",
		"-rs232", "null_modem", "-bitb", "socket.localhost:2323",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func mameExe() string {
	exe := "mame"
	if runtime.GOOS == "windows" {
		exe = "mame64.exe"
	}
	return exe
}