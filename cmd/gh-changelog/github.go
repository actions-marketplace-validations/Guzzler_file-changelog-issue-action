package main

import (
	"context"
	"fmt"
	"time"

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
func commentChangeLogChanges(
	ctx context.Context,
	client *github.Client,
	owner, repo string,
	number int,
) (bool, error) {
	files, err := getAllPullRequestFiles(ctx, client, owner, repo, number)
	if err != nil {
		return false, fmt.Errorf("get all commit files: %w", err)
	}
	currentTime := time.Now().String()
	printableString := fmt.Sprintf("\nTo be added to change Log on %s: \n Owner: %s\n Number: %d\nFiles Changed:\n", currentTime, owner, number)
	for _, file := range files {
		printableString = fmt.Sprintf("%s %s\n", printableString, *file.Filename)
	}
	prOptions := new(github.ListOptions)
	prOptions.Page = 1
	prOptions.PerPage = 1
	allCommits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, number, prOptions)
	prComment := new(github.PullRequestComment)
	lineToComment := 1
	prComment.Body = &printableString
	prComment.CommitID = allCommits[0].SHA
	prComment.Line = &lineToComment
	prComment.Path = files[0].Filename
	if _, _, err := client.PullRequests.CreateComment(ctx, owner, repo, number, prComment); err != nil {
		return false, fmt.Errorf("get all commit files: %w", err)
	}
	return true, nil
}
