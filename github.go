package helpbot

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func initGithub(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
