package assemble

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/mod/semver"
)

const (
	EntryMode = "100644"
	EntryType = "blob"
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
	cmd.action.Debugf("Found tag_Name %q", cmd.tagName)

	commit, err := cmd.createCommit(ctx)
	if err != nil {
		return err
	}
	cmd.action.Debugf("Created new commit %q", commit)

	tagger := NewTagger(cmd.gh, cmd.gtx)
	_, err = tagger.CreateOrUpdateTag(ctx, cmd.tagName, commit)
	if err != nil {
		return err
	}

	// if we aren't dealing with a release event return
	if cmd.rel == nil {
		return nil
	}

	// if we have a draft or pre-release release return
	if *cmd.rel.Draft || *cmd.rel.Prerelease {
		return nil
	}

	major := semver.Major(cmd.tagName)
	majorMinor := semver.MajorMinor(cmd.tagName)
	if _, err = tagger.CreateOrUpdateTag(ctx, major, commit); err != nil {
		return err
	}
	if _, err = tagger.CreateOrUpdateTag(ctx, majorMinor, commit); err != nil {
		return err
	}

	return nil
}

func (a *assembler) createCommit(ctx context.Context) (*github.Commit, error) {
	owner, repo := a.gtx.Repo()

	entries, err := createEntries(a.action)
	if err != nil {
		return nil, err
	}

	tree, _, err := a.gh.Git.CreateTree(ctx, owner, repo, "", entries)
	if err != nil {
		return nil, err
	}

	parent := &github.Commit{
		SHA: github.String(a.gtx.SHA),
	}

	commit := &github.Commit{
		Message: github.String("Automatic compilation"),
		Tree:    tree,
		Parents: []*github.Commit{parent},
	}

	opts := &github.CreateCommitOptions{}

	result, _, err := a.gh.Git.CreateCommit(ctx, owner, repo, commit, opts)
	return result, err
}

func (a *assembler) updateTag(ctx context.Context, commit *github.Commit) (*github.Reference, error) {
	owner, repo := a.gtx.Repo()

	ref := &github.Reference{
		Ref: github.String(a.tagName),
	}

	result, _, err := a.gh.Git.UpdateRef(ctx, owner, repo, ref, true)
	return result, err
}

func createEntries(action *githubactions.Action) ([]*github.TreeEntry, error) {
	entries := []*github.TreeEntry{}

	cwd, err := os.Getwd()
	if err != nil {
		return entries, err
	}

	files, err := os.ReadDir(filepath.Join(cwd, "bin"))
	if err != nil {
		return entries, err
	}

	for _, file := range files {
		path := file.Name()
		entry, err := createEntry(path)
		if err != nil {
			return entries, err
		}

		action.Debugf("appending %q to bare tree", path)
		entries = append(entries, entry)
	}

	yaml, err := createEntry(filepath.Join(cwd, "action.yaml"))
	if err != nil {
		return entries, err
	}
	action.Debugf("appending %q to bare tree", *yaml.Path)
	entries = append(entries, yaml)

	shim, err := createEntry(filepath.Join(cwd, "shim", "invoke-binary.js"))
	if err != nil {
		return entries, err
	}
	action.Debugf("appending %q to bare tree", *shim.Path)
	entries = append(entries, shim)

	return entries, nil
}

func createEntry(path string) (*github.TreeEntry, error) {
	mode := "100644"
	blob := "blob"
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return &github.TreeEntry{
		Path:    github.String(path),
		Mode:    github.String(mode),
		Type:    github.String(blob),
		Content: github.String(content),
	}, nil
}

func readFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := string(bytes[:])

	return content, nil
}
