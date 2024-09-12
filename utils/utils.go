package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/okt-limonikas/go-release/config"
)

func Execute(name string, args []string, env map[string]string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up environment
	cmd.Env = os.Environ() // Start with current environment
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running %s command: %v\n", name, err)
	}
}

func ChangeDirectory(dir string) {
	log.Printf("Changing directory to: %s\n", dir)
	if err := os.Chdir(dir); err != nil {
		log.Fatal("Error changing directory:", err)
	}
}

func ExecuteCommand(cmd config.Command) {
	if len(cmd.Cmd) < 1 {
		log.Fatal("Must provide name of the command to run")
	}

	Execute(cmd.Cmd[0], cmd.Cmd[1:], cmd.Env)
}

func ExecuteMultiple(commands map[string]config.Command) {
	var wg sync.WaitGroup

	wg.Add(len(commands))
	for _, command := range commands {
		go func(cmd config.Command) {
			defer wg.Done()
			log.Printf("Executing command: %s %v\n", cmd.Cmd[0], cmd.Cmd[1:])
			ExecuteCommand(cmd)
		}(command)
	}
	wg.Wait()
}
