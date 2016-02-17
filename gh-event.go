package main

import (
	"fmt"
	"strings"
)

/// Github types

type GithubUser struct {
	Name  string
	Email string
}

type GithubRepo struct {
	Name     string
	Fullname string `json:"full_name"`
	Url      string
	Private  bool
	Owner    GithubUser
}

type GithubCommit struct {
	Id        string
	Message   string
	Timestamp string
	Url       string
	Author    GithubUser
}

type GithubPayload struct {
	Ref        string
	Repository GithubRepo
	Commits    []GithubCommit
}

func (g *GithubPayload) RepoName() string {
	return g.Repository.Fullname
}

func (g *GithubPayload) BranchName() string {
	return strings.TrimPrefix(g.Ref, "refs/heads/")
}

func (g *GithubPayload) EnvData() []string {
	commit := g.Commits[0]

	return []string{
		env("REPO", g.RepoName()),
		env("REPO_URL", g.Repository.Url),
		env("PRIVATE", fmt.Sprintf("%t", g.Repository.Private)),
		env("BRANCH", g.Ref),
		env("COMMIT", commit.Id),
		env("COMMIT_MESSAGE", commit.Message),
		env("COMMIT_TIME", commit.Timestamp),
		env("COMMIT_AUTHOR", commit.Author.Name),
		env("COMMIT_URL", commit.Url),
	}
}

func (g *GithubPayload) Trigger() bool {
	return true
}
