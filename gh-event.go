package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func IsGithubPayload(r *http.Request) bool {
	return r.Header.Get("X-Github-Event") == "push" &&
		(r.Header.Get("Content-Type") == "application/json" || r.Header.Get("Content-Type") == "application/x-www-form-urlencoded")
}

func ExtractGithubPayload(r *http.Request) (Payload, error) {
	payload := new(GithubPayload)

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/x-www-form-urlencoded":
		err := json.Unmarshal([]byte(r.PostFormValue("payload")), payload)
		return payload, err
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(payload)
		return payload, err
	default:
		return nil, fmt.Errorf("Uknown content type %s", contentType)
	}
}
