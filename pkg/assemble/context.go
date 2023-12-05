package assemble

import (
	"errors"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

const (
	inputTagName   = "tag_name"
	envGithubToken = "GITHUB_TOKEN"
)

var (
	errNoTagName     = errors.New("no tag_name was found or provided")
	errNoGithubToken = errors.New("no GITHUB_TOKEN was found or provided")
)

func NewFromContext(action *githubactions.Action) (Command, error) {
	tagName, err := getTagName(action)
	if err != nil {
		return nil, err
	}

	token, err := getGithubToken(action)
	if err != nil {
		return nil, err
	}

	action.AddMask(token)
	client := github.NewClient(nil).WithAuthToken(token)

	return New(
		WithTagName(tagName),
		WithAction(action),
		WithGithubClient(client),
	)
}

func getTagName(action *githubactions.Action) (string, error) {

	if tagName := action.GetInput("tag_name"); tagName != "" {
		return tagName, nil
	}

	gtx, err := action.Context()
	if err != nil {
		return "", err
	}
	if gtx.EventName == "release" {
		return gtx.Event["release"].(map[string]any)["tag_name"].(string), nil
	}

	return "", errNoTagName
}

func getGithubToken(action *githubactions.Action) (string, error) {
	token := action.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", errNoGithubToken
	}

	return token, nil
}
