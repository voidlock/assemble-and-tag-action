package assemble

import (
	"context"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/sethvargo/go-githubactions"
)

func TestTagger(t *testing.T) {
	mockedHTTPClient := mock.NewMockedHTTPClient()

	env := map[string]string{
		"GITHUB_REPOSITORY": "mock/action",
	}
	action := githubactions.New(githubactions.WithGetenv(func(k string) string {
		return env[k]
	}))
	m := mediator{
		gh:  github.NewClient(mockedHTTPClient),
		gtx: action.Context(),
	}

	m.CreateOrUpdateTag(context.Background(), "v1.0.0", commit *github.Commit)
}
