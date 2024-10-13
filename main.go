package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	Dotfiles []string `json:"dotfiles"`
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

	if _, err := os.Stat("./dotfiles"); os.IsNotExist(err) {
		err = os.Mkdir("./dotfiles", 0755)
		if err != nil {
			log.Fatal(err, "Error")
		}

		exec.Command("bash", "-c", `
			cd ./dotfiles/
			git init
			git add .
			git commit -m "0"
			git branch -M main
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
