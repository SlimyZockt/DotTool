package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	GitRepoDir string   `json:"git_repo_dir"`
	Dotfiles   []string `json:"dotfiles"`
}

var config Config

func main() {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(path)

	config_path := filepath.Join(dir, "config.json")

	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		if err != nil {
			b, err := json.Marshal(&config)
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile(config_path, b, 0644)
		}
	}

	config_file, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(config_file, &config)

	if err != nil {
		log.Fatal(err, " Cant convert to struct")
	}

	for _, path := range config.Dotfiles {
		file, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		out_dir := filepath.Join(
			config.GitRepoDir,
			filepath.Base(filepath.Dir(path)),
		)

		out_file := filepath.Join(
			out_dir,
			stat.Name(),
		)

		err = os.MkdirAll(out_dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(out_file, file, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

	}

	cmd := exec.Command("bash", "-c",
		"cd "+config.GitRepoDir+`
		git diff -U0
		git add .
		git commit -am "$(($(git log -1 --pretty=%B)+1))"
		git push
		`)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
