package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gen2brain/go-unarr"
	"github.com/schollz/progressbar/v3"
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

func extract7zv2(archive string, offset int64, dstdir string) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	} else if fi.Size() < offset {
		return fmt.Errorf("offset is too large")
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("seek failed: %w", err)
	}

	ra := io.NewSectionReader(f, offset, fi.Size()-offset)

	totalSz, err := getTotalSz(ra)
	if err != nil {
		return err
	}

	ar, err := unarr.NewArchiveFromReader(ra)
	if err != nil {
		return err
	}
	defer ar.Close()

	pb := progressbar.New64(totalSz)
	pb.Describe("Extracting MAME")
	// pb := progressbar.New(-1)
	if err := extract(ar, dstdir, func(_ string, size int) {
		pb.Add(size)
	}); err != nil {
		return err
	}
	pb.Finish()
	// if _, err := ar.Extract(dstdir); err != nil {
	// 	return err
	// }

	return nil
}

func getTotalSz(r io.ReadSeeker) (int64, error) {
	a, err := unarr.NewArchiveFromReader(r)
	if err != nil {
		return 0, err
	}
	defer a.Close()
	defer r.Seek(0, io.SeekStart)
	var n int64
	for {
		if err := a.Entry(); err != nil {
			if err == io.EOF {
				break
			}
			return n, err
		}
		n += int64(a.Size())
	}
	return n, nil
}

func extract(a *unarr.Archive, dstdir string, fn func(filename string, size int)) (err error) {
	for {
		e := a.Entry()
		if e != nil {
			if e == io.EOF {
				break
			}

			err = e
			return
		}

		name := a.Name()
		size := a.Size()
		if fn != nil {
			fn(name, size)
		}
		data, e := a.ReadAll()
		if e != nil {
			err = e
			return
		}

		dirname := filepath.Join(dstdir, filepath.Dir(name))
		os.MkdirAll(dirname, 0755)

		e = os.WriteFile(filepath.Join(dirname, filepath.Base(name)), data, 0644)
		if e != nil {
			err = e
			return
		}
	}

	return
}
