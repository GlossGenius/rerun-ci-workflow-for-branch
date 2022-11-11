# Rerun ci workflow for branch

This tool is used to re-trigger a circle ci workflow of a circle ci pipeline based on a prefix of a git branch. The 
original inspiration for this work comes from [this](https://glossgenius.slack.com/archives/C034J6ZLJJ3/p1667303969324539)
Slack conversation where a manual action is required to redeploy the release branch to staging.
## Inputs

## `branch_prefix`
**Required.** The prefix of the branch we wish to rerun a workflow for

## `circle_ci_project_slug`
**Required.** The project slug circle ci uses to access the workflow. EX: gh/GlossGenius/core-api

## `circle_ci_workflow_name` 
**Required.** The name of the circle ci workflow to rerun.

## `github_project_slug`
**Required.** The project to search for a `branch_prefix` in. ex: `core-api`

## `github_repository_owner`
**Required.** The GitHub organization who owns the `github_project_slug` project. EX: GlossGenius

## `dry_run`
Optional. Defaults to false. If true, will prevent rerunning of any workflow on circle ci

# external API access environment variables

The project requires additional configuration be passed in through environment variables. Those are:

`CIRCLE_CI_TOKEN`: the token used to access the circle ci api

`GITHUB_TOKEN`: the token used to access the github API

## Example

Based on the Slack conversation above, we would run this tool with:

jobs:
  # deploy-release-to-staging re-triggers the circleci workflow that deploys the release candidate
  deploy-release-to-staging:
    if: github.event.pull_request.merged == true && startsWith(github.head_ref, 'hotfix/')
    runs-on: ubuntu-latest
    name: deploy-release-to-staging
    steps:
      - name: rerun deploy_staging
        uses: GlossGenius/rerun-ci-workflow-for-branch@v1
        id: rerun_workflow
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          CIRCLE_CI_TOKEN: ${{ secrets.CIRCLECI_API_TOKEN }}
        with:
          branch_prefix: release/
          circle_ci_workflow_name: deploy_staging
          circle_ci_project_slug: gh/${{ github.repository }}
          github_repository_owner: ${{ github.repository_owner }}
          github_project_slug: ${{ github.event.pull_request.base.repo.name }}
          # TODO: delete this when done testing
          dry_run: true

## Limitations

One call out that should be made: the circleci api currently lets you rerun workflows that are already running.
In the code, we protect against this by checking the workflow's status and doing a noop if it's currently running.
However, there is a race condition possible here. When a user triggers a rerun of a workflow, and before the circleci 
API persists the new workflow, another user can rerun the same workflow thereby allowing the same workflow being run 
multiple times even if we don't want this behaviour to happen.