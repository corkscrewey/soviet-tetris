package main

import (
	_ "embed"
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

var dirstructure = []struct {
	Dir      string
	Name     string
	Contents *[]byte
	Mode     os.FileMode
}{
	{"tetris", "Dockerfile", &dockerfile, 0644},
	{"tetris", "docker-compose.yml", &dockercompose, 0644},
	{"roms", "ie15.zip", &ie15, 0644},
	{"roms", "ie15kbd.zip", &ie15kbd, 0644},
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
