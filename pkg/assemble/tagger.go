package assemble

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

type Tagger interface {
	CreateOrUpdateTag(ctx context.Context, ref string, commit *github.Commit) (*github.Reference, error)
}

func NewTagger(gh *github.Client, gtx *githubactions.GitHubContext) Tagger {
	return mediator{gh, gtx}
}

type mediator struct {
	gh  *github.Client
	gtx *githubactions.GitHubContext
}

func (m mediator) CreateOrUpdateTag(ctx context.Context, tagName string, commit *github.Commit) (*github.Reference, error) {
	tag := "refs/tags/" + tagName
	owner, repo := m.gtx.Repo()

	ref, _, err := m.gh.Git.GetRef(ctx, owner, repo, tag)
	if err == nil {
		// attach the new commit sha to the existing reference
		ref.Object.SHA = commit.SHA

		newRef, _, err := m.gh.Git.UpdateRef(ctx, owner, repo, ref, true)
		return newRef, err
	}

	// create a new ref for the tag
	newRef := &github.Reference{Ref: github.String(tag), Object: &github.GitObject{SHA: commit.SHA}}
	ref, _, err = m.gh.Git.CreateRef(ctx, owner, repo, newRef)
	return ref, err
}
