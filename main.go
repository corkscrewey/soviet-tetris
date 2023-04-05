package main

import (
	"context"
	"flag"
	"log"

	"golang.org/x/sync/errgroup"
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
	if _, err := dirExists(workdir); err != nil {
		return err
	}
	if err := bootstrap(workdir); err != nil {
		return err
	}

	// check docker

	doco, err := initTetrisDocker(ctx, workdir, dirTetris)
	if err != nil {
		return err
	}

	// check mame if not exist, add download action
	// check systemwide mame
	mame, err := NewMAMEFromPath(workdir, dirROM)
	if err != nil {
		return err
	}
	log.Println(mame.exe)
	return runEmulation(ctx, doco, mame, workdir, mameExe(), dirTetris, dirROM)
}

func runEmulation(ctx context.Context, doco *Compose, mame *MAME, workdir string, mameExe, tetrisDir, romDir string) error {
	log.Println("running emulation")
	// run docker compose
	var eg errgroup.Group
	eg.Go(func() error {
		return doco.Up(ctx, "docker-compose.yaml", "simh")
	})
	eg.Go(func() error {
		return mame.Run()
	})

	return eg.Wait()
}
