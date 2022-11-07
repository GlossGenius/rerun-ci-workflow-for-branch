# Rerun ci workflow for branch

This tool is used to re-trigger a circle ci workflow of a circle ci pipeline based on a prefix of a git branch. The 
original inspiration for this work comes from [this](https://glossgenius.slack.com/archives/C034J6ZLJJ3/p1667303969324539)
Slack conversation where a manual action is required to redeploy the release branch to staging.

## Configuration

The tool is configured via environment variables and configure the following functionality:

```text
BRANCH_PREFIX: the prefix of the branch we wish to rerun a workflow for
DRY_RUN: if true, will prevent the workflow from being rerun
CIRCLE_CI_PROJECT_SLUG: the project slug circle ci uses to access the workflow. EX: gh/GlossGenius/core-api
CIRCLE_CI_WORKFLOW_NAME: the name of the workflow to rerun
GITHUB_ORG_SLUG: the github organization who owns the GITHUB_PROJECT_SLUG. EX: GlossGenius
GITHUB_PROJECT_SLUG: the github project where the branch with BRANCH_PREFIX exists. EX: core-api


# external API access
CIRCLE_CI_TOKEN: the token used to access the circle ci api
GITHUB_TOKEN: the token used to access the github API
```

## Example

Based on the Slack conversation above, we would run this tool with:

```text
BRANCH_PREFIX=release/
CIRCLE_CI_WORKFLOW_NAME=deploy_staging
CIRCLE_CI_TOKEN=<token>
CIRCLE_CI_PROJECT_SLUG=gh/GlossGenius/core-api
GITHUB_ORG_SLUG=GlossGenius
GITHUB_PROJECT_SLUG=core-api
GITHUB_TOKEN=<token>
```

to rerun the release pipeline in question. If combine this with a GitHub action that triggers on merge of a `hotfix/*` 
branch, then we can successfully automate the manual action taken after a hotfix branch merges to master.

## Limitations

One call out that should be made: the circleci api currently lets you rerun workflows that are already running.
In the code, we protect against this by checking the workflow's status and doing a noop if it's currently running.
However, there is a race condition possible here. When a user triggers a rerun of a workflow, and before the circleci 
API persists the new workflow, another user can rerun the same workflow thereby allowing the same workflow being run 
multiple times even if we don't want this behaviour to happen.