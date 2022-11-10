package providers

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
)

type Github struct {
	client       *github.Client
	branchPrefix string
	projectSlug  string
	org          string
}

func NewGithub(ctx context.Context, apiKey, branchPrefix, projectSlug, org string) *Github {
	return &Github{
		client:       getGithubClient(ctx, apiKey),
		branchPrefix: branchPrefix,
		projectSlug:  projectSlug,
		org:          org,
	}
}

// getGithubClient creates a new client based on a GitHub access token
func getGithubClient(ctx context.Context, accessToken string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	return github.NewClient(tokenClient)
}

// getAllBranches fetches all branches in all pages for the github context
func (g *Github) getAllBranches(ctx context.Context) ([]*github.Branch, error) {
	branches := []*github.Branch{}
	listOptions := github.ListOptions{}
	for {
		branchesInPage, response, err := g.client.Repositories.ListBranches(
			ctx, g.org, g.projectSlug, &github.BranchListOptions{ListOptions: listOptions},
		)
		if err != nil {
			return []*github.Branch{}, fmt.Errorf("listing branches: %w", err)
		}
		branches = append(branches, branchesInPage...)
		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage
	}
	return branches, nil
}

// GetBranch finds the branch with the prefix stored in the instance of the Github struct
func (g *Github) GetBranch(ctx context.Context) (*github.Branch, error) {
	branches, err := g.getAllBranches(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all branches: %w", err)
	}
	branch, err := getBranchWithPrefix(branches, g.branchPrefix)
	if err != nil {
		return nil, fmt.Errorf("getting branch: %w", err)
	}
	log.Println(fmt.Sprintf("found github branch with name: %s", branch.GetName()))
	return branch, nil
}

// getBranchWithPrefix finds the branch for a given prefix, raises an error if it can't disambiguate
func getBranchWithPrefix(branches []*github.Branch, prefix string) (*github.Branch, error) {
	branches = filterForPrefix(branches, prefix)
	switch len(branches) {
	case 0:
		break
	case 1:
		return branches[0], nil
	default:
		return nil, &cantDisambiguateBranchesError{actualBranches: len(branches)}
	}
	return nil, nil
}

// filterForPrefix finds the branches that have the prefix specified
func filterForPrefix(branches []*github.Branch, prefix string) []*github.Branch {
	result := []*github.Branch{}
	for _, branch := range branches {
		if branch != nil && strings.HasPrefix(*branch.Name, prefix) {
			result = append(result, branch)
		}
	}
	return result
}
