package assemble

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/v57/github"
	"github.com/sethvargo/go-githubactions"
)

const (
	EntryFileMode       = "100644"
	EntryExecutableMode = "100755"
	EntryType           = "blob"
)

type TreeBuilder interface {
	CreateTree(ctx context.Context) (*github.Tree, error)
}

type builder struct {
	gtx    *githubactions.GitHubContext
	action *githubactions.Action
	gh     *github.Client
}

func NewTreeBuilder(gtx *githubactions.GitHubContext, action *githubactions.Action, gh *github.Client) TreeBuilder {
	return builder{
		gtx:    gtx,
		action: action,
		gh:     gh,
	}
}

func (b builder) CreateTree(ctx context.Context) (*github.Tree, error) {
	owner, repo := b.gtx.Repo()

	entries, err := createEntries(b.action)
	if err != nil {
		return nil, err
	}

	tree, _, err := b.gh.Git.CreateTree(ctx, owner, repo, "", entries)
	return tree, err
}

func createEntries(action *githubactions.Action) ([]*github.TreeEntry, error) {
	entries := []*github.TreeEntry{}

	binDir := filepath.Join("bin")
	files, err := os.ReadDir(binDir)
	if err != nil {
		return entries, err
	}

	for _, file := range files {
		path := filepath.Join(binDir, file.Name())
		entry, rerr := createEntry(path, EntryExecutableMode)
		if rerr != nil {
			return entries, err
		}

		action.Debugf("appending %q to bare tree", path)
		entries = append(entries, entry)
	}

	yaml, err := createEntry("action.yaml", EntryFileMode)
	if err != nil {
		return entries, err
	}
	action.Debugf("appending %q to bare tree", *yaml.Path)
	entries = append(entries, yaml)

	shim, err := createEntry(filepath.Join("shim", "invoke-binary.js"), EntryFileMode)
	if err != nil {
		return entries, err
	}
	action.Debugf("appending %q to bare tree", *shim.Path)
	entries = append(entries, shim)

	return entries, nil
}

func createEntry(path, mode string) (*github.TreeEntry, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}

	return &github.TreeEntry{
		Path:    github.String(path),
		Mode:    github.String(mode),
		Type:    github.String(EntryType),
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
