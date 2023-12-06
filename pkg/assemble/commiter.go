package assemble

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

type Committer interface {
	Commit(ctx context.Context, tree *github.Tree) (*github.Commit, error)
}

type committer struct {
	gtx *githubactions.GitHubContext
	gh  *github.Client
}

func NewCommitter(gtx *githubactions.GitHubContext, gh *github.Client) Committer {
	return committer{
		gtx: gtx,
		gh:  gh,
	}
}

func (c committer) Commit(ctx context.Context, tree *github.Tree) (*github.Commit, error) {
	owner, repo := c.gtx.Repo()
	parent := &github.Commit{
		SHA: github.String(c.gtx.SHA),
	}

	commit := &github.Commit{
		Message: github.String("Automatic compilation"),
		Tree:    tree,
		Parents: []*github.Commit{parent},
	}

	opts := &github.CreateCommitOptions{}

	result, _, err := c.gh.Git.CreateCommit(ctx, owner, repo, commit, opts)
	return result, err
}
