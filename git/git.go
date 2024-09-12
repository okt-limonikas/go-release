package git

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/okt-limonikas/go-release/utils"
)

type GitTag struct {
	Tag     string
	TagNote string
	Sha     string
}

func GetTagInfo() GitTag {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Error getting tag info:", err)
	}
	tag := strings.TrimSpace(string(output))

	cmd = exec.Command("git", "tag", "-n1", tag)
	output, err = cmd.Output()
	if err != nil {
		log.Fatal("Error getting tag note:", err)
	}
	tagNote := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(string(output)), tag+" "))

	cmd = exec.Command("git", "rev-parse", "HEAD")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal("Error getting commit sha:", err)
	}
	sha := strings.TrimSpace(string(output))

	return GitTag{Tag: tag, Sha: sha, TagNote: tagNote}
}

func AddCommitAndPush(tag GitTag) {
	fullMessage := fmt.Sprintf("chore(release): %s", tag.Tag)

	utils.Execute("git", []string{"add", "."}, nil)
	utils.Execute("git", []string{"commit", "-m", fullMessage}, nil)
	utils.Execute("git", []string{"tag", "-a", tag.Tag, "-m", tag.TagNote}, nil)
	utils.Execute("git", []string{"push"}, nil)
	utils.Execute("git", []string{"push", "origin", "tag", tag.Tag}, nil)
	fmt.Printf("Files added and committed with tag %s\n", tag)
}

func ResetStagingArea() {
	cmd := exec.Command("git", "reset", "--hard", "HEAD")
	if err := cmd.Run(); err != nil {
		log.Fatal("Error resetting staging area:", err)
	}

	fmt.Println("Staging area reset successfully.")
}

func CheckoutTag(tag string) {
	utils.Execute("git", []string{"fetch", "--all"}, nil)
	utils.Execute("git", []string{"checkout", tag}, nil)
	fmt.Printf("Checked out tag: %s\n", tag)
}
