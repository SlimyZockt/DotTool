package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	GitRepoDir string   `json:"git_repo_dir"`
	ID         int      `json:"id"`
	Dotfiles   []string `json:"dotfiles"`
}

func main() {
	var config Config

	config_file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(config_file, &config)
	if err != nil {
		log.Fatal(err, " Cant convert to struct")
	}

	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		if err != nil {
			log.Fatal(err, "Error")
		}

		exec.Command("bash", "-c", `
			`).Run()
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

		os.WriteFile("./dotfiles/"+stat.Name(), file, 0644)
	}

	log.Print(config.ID)

	cmd := exec.Command("bash", "-c", `
		cd ./dotfiles/
		git add .
		git diff
		git diff
		git push
	`)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
