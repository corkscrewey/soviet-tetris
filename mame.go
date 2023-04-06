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
	localMameDir := filepath.Join(workdir, "mame")
	localROMDir := filepath.Join(workdir, romDir)
	exepath, err := exec.LookPath(exe)
	if err != nil {
		// not found in PATH, download and unpack.
		exepath = filepath.Join(localMameDir, mameExe())
		if err := installMame(exepath); err != nil {
			return nil, fmt.Errorf("error installing mame: %w", err)
		}
	}
	return &MAME{exe: exepath, romDir: localROMDir}, nil
}

func (m *MAME) Run() error {
	cmd := exec.Command(m.exe,
		"-rompath", m.romDir,
		"-window",
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
		exe = "mame.exe"
	}
	return exe
}
