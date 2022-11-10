package main

import (
	"context"
	"fmt"
	"github.com/GlossGenius/rerun-ci-workflow-for-branch/internal/helpers"
	"github.com/GlossGenius/rerun-ci-workflow-for-branch/internal/providers"
	"log"
	"os"
)

// TODO(michel): I could be easily convinced to used structured logging if we wanted that

var (
	// NB this token is documented here: https://docs.github.com/en/actions/security-guides/automatic-token-authentication
	githubAccessToken   = os.Getenv("GITHUB_TOKEN")
	githubOrg           = os.Getenv("GITHUB_ORG_SLUG")
	githubProjectSlug   = os.Getenv("GITHUB_REPOSITORY_OWNER")
	branchPrefix        = os.Getenv("BRANCH_PREFIX")
	circleCIToken       = os.Getenv("CIRCLE_CI_TOKEN")
	workflowName        = os.Getenv("CIRCLE_CI_WORKFLOW_NAME")
	circleCIProjectSlug = os.Getenv("CIRCLE_CI_PROJECT_SLUG")
	dryRun              = os.Getenv("DRY_RUN")
)

func run(ctx context.Context, github *providers.Github, circleci *providers.CircleCI) {
	branch, err := github.GetBranch(ctx)
	if err != nil {
		panic(fmt.Errorf("getting branch %w: ", err))
	}
	// do nothing since we have nothing to rerun
	if branch == nil {
		log.Printf("WARNING: branch with prefix %s does not exist, exiting.", branchPrefix)
		return
	}
	if err := circleci.TriggerWorkflow(ctx, branch.GetName(), workflowName); err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	dryRunParsed := helpers.MustParseBool(dryRun)
	circleci := providers.NewCircleCI(circleCIToken, circleCIProjectSlug, dryRunParsed)
	github := providers.NewGithub(ctx, githubAccessToken, branchPrefix, githubProjectSlug, githubOrg)
	run(ctx, github, circleci)
}
