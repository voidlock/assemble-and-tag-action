package assemble

import (
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

type Option func(*assembler) error

func WithTagName(tagName string) Option {
	return func(cmd *assembler) error {
		cmd.tagName = tagName
		return nil
	}
}

func WithAction(action *githubactions.Action) Option {
	return func(cmd *assembler) error {

		// We use this context enough that it is worth storing it
		gtx, err := action.Context()
		if err != nil {
			return err
		}

		cmd.action = action
		cmd.gtx = gtx

		// If we have a release event pre-parse the release info out of it. this is
		// easier than trying to work with the action Context.Event map[string]any
		if gtx.EventName == "release" {
			payload, err := os.ReadFile(cmd.gtx.EventPath)
			if err != nil {
				return err
			}
			event, err := github.ParseWebHook(cmd.gtx.EventName, payload)
			if err != nil {
				return err
			}
			releaseEvent := event.(github.ReleaseEvent)
			release := releaseEvent.GetRelease()

			cmd.rel = release
		}

		return nil
	}
}

func WithGithubClient(client *github.Client) Option {
	return func(cmd *assembler) error {
		cmd.gh = client
		return nil
	}
}
