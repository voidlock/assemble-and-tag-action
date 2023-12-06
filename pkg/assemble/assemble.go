package assemble

import (
	"context"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/mod/semver"
)

type Command interface {
	Run(context.Context) error
}

type assembler struct {
	tagName string
	action  *githubactions.Action
	gtx     *githubactions.GitHubContext
	gh      *github.Client
	rel     *github.RepositoryRelease
}

func New(opts ...Option) (Command, error) {
	cmd := &assembler{}

	for _, opt := range opts {
		if err := opt(cmd); err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

func (cmd *assembler) Run(ctx context.Context) error {
	// if we aren't dealing with a release event return
	if cmd.rel == nil {
		cmd.action.Infof("Workflow not triggered by a `release` event, returning.")
		return nil
	}

	// if we have a draft or pre-release release return
	if cmd.rel.GetDraft() || cmd.rel.GetPrerelease() {
		cmd.action.Infof("Workflow triggered by a 'draft' or 'pre-release' event, returning.")
		return nil
	}

	cmd.action.Debugf("Found tag_Name %q", cmd.tagName)

	cmd.action.Infof("Building action release tree")
	builder := NewTreeBuilder(cmd.gtx, cmd.action, cmd.gh)
	tree, err := builder.CreateTree(ctx)
	if err != nil {
		return err
	}
	cmd.action.Infof("Created tree %q", github.Stringify(tree))

	committer := NewCommitter(cmd.gtx, cmd.gh)
	commit, err := committer.Commit(ctx, tree)
	if err != nil {
		return err
	}
	cmd.action.Infof("Created new commit %q", commit.GetSHA())

	tagger := NewTagger(cmd.gh, cmd.gtx)
	tag, err := tagger.CreateOrUpdateTag(ctx, cmd.tagName, commit)
	if err != nil {
		return err
	}
	cmd.action.Infof("Updated %q to match commit %q", tag.GetRef(), commit.GetSHA())

	major := semver.Major(cmd.tagName)
	tag, err = tagger.CreateOrUpdateTag(ctx, major, commit)
	if err != nil {
		return err
	}
	cmd.action.Infof("Updated %q to match commit %q", tag.GetRef(), commit.GetSHA())

	majorMinor := semver.MajorMinor(cmd.tagName)
	tag, err = tagger.CreateOrUpdateTag(ctx, majorMinor, commit)
	if err != nil {
		return err
	}
	cmd.action.Infof("Updated %q to match commit %q", tag.GetRef(), commit.GetSHA())

	cmd.action.Infof("Completed")
	return nil
}
