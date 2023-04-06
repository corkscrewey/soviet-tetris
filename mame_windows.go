package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const mameURL = "https://github.com/mamedev/mame/releases/download/mame0253/mame0253b_64bit.exe"

func installMame(mameExe string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(mameExe)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(wd, abs)
	if err != nil {
		return err
	}
	log.Printf("Installing MAME to %s", rel)
	resp, err := http.Get(mameURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	tf, err := ioutil.TempFile("", "mame-dist*.exe")
	if err != nil {
		return err
	}
	defer os.Remove(tf.Name())
	defer tf.Close()
	if _, err = io.Copy(tf, resp.Body); err != nil {
		return err
	}
	if err = tf.Close(); err != nil {
		return err
	}
	log.Printf("Unpacking MAME to %s", filepath.Dir(rel))
	cmd := exec.Command(tf.Name(), "-o\""+filepath.Dir(rel)+`"`, "-y")
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}
