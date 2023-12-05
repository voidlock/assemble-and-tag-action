package assemble

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/sethvargo/go-githubactions"
)

func TestNewFromInputs(t *testing.T) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	if _, err := f.Write([]byte(`{"action": "published", "release": { "tag_name": "release_tag" }}`)); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	eventPayloadPath := f.Name()

	cases := []struct {
		name   string
		expTag string
		env    map[string]string
	}{
		{
			name:   "from action input",
			expTag: "input_tag",
			env: map[string]string{
				"INPUT_TAG_NAME": "input_tag",
				"GITHUB_TOKEN":   "fake",
			},
		},
		{
			name:   "from payload",
			expTag: "release_tag",
			env: map[string]string{
				"GITHUB_TOKEN":      "fake",
				"GITHUB_EVENT_NAME": "release",
				"GITHUB_EVENT_PATH": eventPayloadPath,
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actionLog := bytes.NewBuffer(nil)
			getenv := func(key string) string {
				return tc.env[key]
			}
			action := githubactions.New(
				githubactions.WithWriter(actionLog),
				githubactions.WithGetenv(getenv),
			)

			cmd, err := NewFromContext(action)
			if err != nil {
				t.Fatalf("faild with error %q", err)
			}

			cmd.Run(context.Background())

			if log := actionLog.String(); !strings.Contains(log, tc.expTag) {
				t.Fatalf("unexpected log found: %#v", log)
			}
		})
	}
}
