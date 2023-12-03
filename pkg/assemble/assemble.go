package assemble

import (
	"context"
	"errors"

	"github.com/sethvargo/go-githubactions"
)

var errNoTagName = errors.New("no tag_name was found or provided")

type Command interface {
	Run(context.Context) error
}

type assembler struct {
	action  *githubactions.Action
	tagName string
}

func NewFromInputs(action *githubactions.Action) (Command, error) {
	tagName, err := getTagName(action)
	return &assembler{action: action, tagName: tagName}, err
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

func (a *assembler) Run(_ context.Context) error {
	a.action.Infof("Found tag_Name %q", a.tagName)
	return nil
}
