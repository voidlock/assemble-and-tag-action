package assemble

// import (
// 	"testing"
//
// 	"github.com/migueleliasweb/go-github-mock/src/mock"
// 	"github.com/sethvargo/go-githubactions"
// )
//
// func TestTagger(t *testing.T) {
// 	mockedHTTPClient := mock.NewMockedHTTPClient()
//
// 	env := map[string]string{
// 		"GITHUB_REPOSITORY": "mock/action",
// 	}
// 	action := githubactions.New(githubactions.WithGetenv(func(k string) string {
// 		return env[k]
// 	}))
//
// 	_, err := action.Context()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// m := mediator{
// 	// 	gh:  github.NewClient(mockedHTTPClient),
// 	// 	gtx: gtx,
// 	// }
//
// 	//m.CreateOrUpdateTag(context.Background(), "v1.0.0", commit*github.Commit)
// }
