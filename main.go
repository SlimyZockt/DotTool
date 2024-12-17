package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
)

type Config struct {
	GitRepoDir string   `json:"git_repo_dir"`
	Dotfiles   []string `json:"dotfiles"`
}

var config Config
var id int

func main() {
	interrupt_chan := make(chan os.Signal, 1)
	signal.Notify(interrupt_chan, os.Interrupt)
	go func() {
		<-interrupt_chan

		err := os.WriteFile(config.GitRepoDir+"id", []byte(strconv.Itoa(id-1)), 0644)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(1)
	}()

	app()
}

func app() {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		if err != nil {
			b, err := json.Marshal(&config)
			if err != nil {
				log.Fatal(err)
			}
			os.WriteFile("./config.json", b, 0644)
		}
	}

	config_file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(config_file, &config)

	log.Println(config)

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

		err = os.WriteFile(config.GitRepoDir+stat.Name(), file, 0644)
		if err != nil {
			log.Fatal(err)
		}

	}

	id_data, err := os.ReadFile(config.GitRepoDir + "id")

	if os.IsNotExist(err) {
		err := os.WriteFile(config.GitRepoDir+"id", []byte("0"), 0644)

		if err != nil {
			log.Fatal(err)
		}

		id = 0
	} else if err == nil {
		new_id, err := strconv.Atoi(string(id_data))
		if err != nil {
			log.Fatal(err)
		}
		id = new_id + 1

	} else {
		log.Fatal(err)
	}

	err = os.WriteFile(config.GitRepoDir+"id", []byte(strconv.Itoa(id)), 0644)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("bash", "-c",
		"cd "+config.GitRepoDir+`
		git diff -U0
		git commit -am "`+
			strconv.Itoa(id)+`"
		git push
		`)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
