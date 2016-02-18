package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/// Bitbucket types
type BitbucketUser struct {
	Name     string `json:"display_name"`
	Username string
}

type BitbucketRepo struct {
	Name     string
	Fullname string `json:"full_name"`
	Private  bool   `json:"is_private"`
}

type BitbucketPush struct {
	Changes []BitbucketChange
}

type BitbucketChange struct {
	New     *BitbucketRef
	Old     *BitbucketRef
	Created bool
	Forced  bool
	Closed  bool
}

type BitbucketRef struct {
	Type   string
	Name   string
	Target BitbucketCommit
	Links  BitbucketLinks
}

type BitbucketCommit struct {
	Type      string
	Hash      string
	Author    BitbucketUser
	Message   string
	Timestamp string
	Links     BitbucketLinks
}

type BitbucketLinks struct {
	Html BitbucketLink
}

type BitbucketLink struct {
	Href string
}

type BitbucketPayload struct {
	Actor      BitbucketUser
	Repository BitbucketRepo
	Push       *BitbucketPush
}

func (b *BitbucketPayload) Trigger() bool {
	return b.Push != nil &&
		len(b.Push.Changes) > 0 &&
		b.Push.Changes[0].New != nil &&
		b.Push.Changes[0].New.Type == "branch" &&
		b.Push.Changes[0].New.Target.Type == "commit"
}

func (b *BitbucketPayload) RepoName() string {
	return b.Repository.Fullname
}

func (b *BitbucketPayload) BranchName() string {
	return b.Push.Changes[0].New.Name
}

func (b *BitbucketPayload) EnvData() []string {
	commit := b.Push.Changes[0].New.Target

	return []string{
		env("REPO", b.RepoName()),
		env("REPO_URL", fmt.Sprintf("https://bitbucket.org/%s", b.RepoName())),
		env("PRIVATE", fmt.Sprintf("%t", b.Repository.Private)),
		env("BRANCH", b.BranchName()),
		env("COMMIT", commit.Hash),
		env("COMMIT_MESSAGE", commit.Message),
		env("COMMIT_TIME", commit.Timestamp),
		env("COMMIT_AUTHOR", commit.Author.Name),
		env("COMMIT_URL", commit.Links.Html.Href),
	}
}

func IsBitbucketPayload(r *http.Request) bool {
	return r.Header.Get("X-Event-Key") == "repo:push" && r.Header.Get("Content-Type") == "application/json"
}

func ExtractBitbucketPayload(r *http.Request) (Payload, error) {
	payload := new(BitbucketPayload)
	err := json.NewDecoder(r.Body).Decode(payload)
	return payload, err
}
