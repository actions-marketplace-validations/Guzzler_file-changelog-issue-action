package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v29/github"

	"go.uber.org/zap"
)

// getAllPullRequestFiles returns all commit files in a pull request.
func getAllPullRequestFiles(
	ctx context.Context,
	client *github.Client,
	owner, repo string,
	number int,
) ([]*github.CommitFile, error) {
	var res []*github.CommitFile
	for offset := 0; ; offset++ {
		files, resp, err := client.PullRequests.ListFiles(
			ctx,
			owner, repo, number,
			&github.ListOptions{Page: offset, PerPage: 100},
		)
		if err != nil {
			logger.Error("Failed to remove a label from a pull request",
				zap.String("owner", owner),
				zap.String("repo", repo),
				zap.Int("number", number),
				zap.Int("offset", offset),
				zap.Error(err),
			)
			return nil, fmt.Errorf("list commit files: %w", err)
		}
		res = append(res, files...)
		if offset >= resp.LastPage {
			break
		}
	}
	return res, nil
}

// addToChangeLog appends the files to the changelog file and adds the file to the commit.
func addToChangeLog(
	ctx context.Context,
	client *github.Client,
	owner, repo, changeLogPath string,
	number int,
) (bool, error) {
	files, err := getAllPullRequestFiles(ctx, client, owner, repo, number)
	if err != nil {
		return false, fmt.Errorf("get all commit files: %w", err)
	}

	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, fmt.Errorf("openLogFile: %w", err)
	}
	defer f.Close()
	printableString := fmt.Sprintf("\nOwner: %s\n Number: %d\nFiles Changed:\n", owner, number)
	for _, file := range files {
		printableString = fmt.Sprintf("%s %s\n", printableString, file)
	}
	if _, err := f.WriteString(printableString); err != nil {
		return false, fmt.Errorf("get all commit files: %w", err)
	}
	return true, nil
}
