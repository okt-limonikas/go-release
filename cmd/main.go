package main

import (
	"log"
	"os"

	"github.com/okt-limonikas/go-release/config"
	"github.com/okt-limonikas/go-release/git"
	"github.com/okt-limonikas/go-release/utils"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	config, err := config.Load(args.path)
	if err != nil {
		log.Fatal(err)
	}

	// 1. Main dir
	tag := createRelease(*config)
	// 2. Go to dir and install
	prepareSync(*config, tag)
	defer cleanupSync(config.Paths.Sync)
	// 3. Build and release
	buildAndRelease(*config, tag)
}

func buildAndRelease(config config.Config, tag git.GitTag) {
	utils.ChangeDirectory(config.Paths.Build)
	utils.ExecuteCommand(config.Commands.Build)
	git.AddCommitAndPush(tag)
}

func cleanupSync(dir string) {
	utils.ChangeDirectory(dir)
	git.ResetStagingArea()
}

func prepareSync(config config.Config, tag git.GitTag) {
	utils.ChangeDirectory(config.Paths.Sync)
	git.CheckoutTag(tag.Tag)
	utils.ExecuteCommand(config.Commands.Install)
}

func createRelease(config config.Config) git.GitTag {
	utils.ChangeDirectory(config.Paths.Main)
	utils.ExecuteCommand(config.Commands.Install)
	utils.ExecuteMultiple(config.Commands.Checks)
	utils.ExecuteCommand(config.Commands.Release)

	return git.GetTagInfo()
}

type Args struct {
	path string
}

func parseArgs() (Args, error) {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the config file as an argument.")
	}

	path := os.Args[1]

	return Args{path}, nil
}
