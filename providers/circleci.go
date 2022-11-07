package providers

import (
	"context"
	"fmt"
	"github.com/grezar/go-circleci"
	"log"
	"sort"
)

const runningStatus = "running"

type CircleCI struct {
	client      *circleci.Client
	projectSlug string
	dryRun      bool
}

func NewCircleCI(token string, projectSlug string, dryRun bool) *CircleCI {
	return &CircleCI{
		client:      mustGetCircleCIClient(token),
		projectSlug: projectSlug,
		dryRun:      dryRun,
	}
}

// TriggerWorkflow gets the most recent pipeline for a branch, then reruns the specified workflow
func (c *CircleCI) TriggerWorkflow(ctx context.Context, branch string, workflowName string) error {
	pipeline, err := c.getMostRecentPipelineForBranch(ctx, branch)
	if err != nil {
		return fmt.Errorf("getting pipelines for branch: %w", err)
	}
	if pipeline == nil {
		log.Printf("WARNING: no pipelines found for branch with prefix %s - no work performed", branch)
		return nil
	}
	workflow, err := c.getWorkflow(ctx, pipeline, workflowName)
	if err != nil {
		return fmt.Errorf("finding workflow: %w", err)
	}
	// TODO(michel): these no-exec cases are a bit loud, they could be abstracted to another function to clean up
	if workflow == nil {
		log.Printf("WARNING: could not find the workflow %s provided - no work done", workflowName)
		return nil
	}
	if workflow.Status == runningStatus {
		log.Printf("INFO: workflow %#v is already running - no work done", workflow)
		return nil
	}
	if c.dryRun {
		log.Printf("INFO: found workflow: %#v but skipping rerun due to dry-run flag being set", workflow)
		return nil
	}
	// todo(michel): good candidate for structured logging
	log.Printf("INFO: retriggered workflow: %s for pipeline %s for branch %s", workflow.ID, pipeline.ID, branch)
	if err := c.client.Workflows.Rerun(ctx, workflow.ID, circleci.WorkflowRerunOptions{}); err != nil {
		return fmt.Errorf("rerunning workflow: %w", err)
	}
	return nil
}

// getWorkflow finds the workflow specified in from the current pipeline
func (c *CircleCI) getWorkflow(ctx context.Context, pipeline *circleci.Pipeline, workflowName string) (*circleci.Workflow, error) {
	// TODO(michel): this currently just grabs the first page of workflows, which at the time of writing is 20 elements.
	//  If we wish to run this for pipelines with >20 workflows, we'd have to fetch all pages
	workflows, err := c.client.Pipelines.ListWorkflows(ctx, pipeline.ID, circleci.PipelineListWorkflowsOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing workflows: %w", err)
	}
	// TODO(michel): are workflow names unique? if not, then this could produce ambiguous results
	//  could not find specs on the circleci website
	for _, workflow := range workflows.Items {
		if workflow.Name == workflowName {
			return workflow, nil
		}
	}
	return nil, nil
}

// getAllPipelinesForBranch given the circleci context and branch of interest, fetch all pipelines
func (c *CircleCI) getAllPipelinesForBranch(ctx context.Context, branch string) ([]*circleci.Pipeline, error) {
	// NB(michel): since the circle CI API offers limited functionality to filter results or change the page size,
	//  we can only grab all results and filter the ones we want. Although today it seems the api is sorting results by date,
	//  there's no supporting documentation for this, and we can't send a query param to ensure that the result is
	//  sorted the way we want it to be when we receive the payload.
	pipelines := []*circleci.Pipeline{}
	pipelineOptions := circleci.ProjectListPipelinesOptions{Branch: &branch}
	for {
		pipelinesPage, err := c.client.Projects.ListPipelines(ctx, c.projectSlug, pipelineOptions)
		if err != nil {
			return []*circleci.Pipeline{}, fmt.Errorf("listing pipelines: %w", err)
		}
		pipelines = append(pipelines, pipelinesPage.Items...)
		if pipelinesPage.NextPageToken == "" {
			break
		}
		pipelineOptions = circleci.ProjectListPipelinesOptions{Branch: &branch, PageToken: &pipelinesPage.NextPageToken}
	}
	return pipelines, nil
}

// getMostRecentPipelineForBranch fetches the most recently created pipeline for a given git branch
func (c *CircleCI) getMostRecentPipelineForBranch(ctx context.Context, branch string) (*circleci.Pipeline, error) {
	pipelines, err := c.getAllPipelinesForBranch(ctx, branch)
	pipelines = sortByCreateDate(removeNils(pipelines))
	if len(pipelines) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting all pipelines: %w", err)
	}
	return pipelines[0], nil
}

// mustGetCircleCIClient creates a circleci client based on token, panics if it can't
func mustGetCircleCIClient(token string) *circleci.Client {
	config := circleci.DefaultConfig()
	config.Token = token

	client, err := circleci.NewClient(config)
	if err != nil {
		panic(fmt.Errorf("creating new circleci client: %w", err))
	}
	return client
}

// removeNils filters a slice of pipelines eliminating its nil values
func removeNils(pipelines []*circleci.Pipeline) []*circleci.Pipeline {
	output := []*circleci.Pipeline{}
	for _, p := range pipelines {
		if p != nil {
			output = append(output, p)
		}
	}
	return output
}

// sortByCreateDate sorts the slice of pipelines by their creation date
func sortByCreateDate(pipeline []*circleci.Pipeline) []*circleci.Pipeline {
	output := make([]*circleci.Pipeline, len(pipeline))
	copy(output, pipeline)
	sort.Slice(output, func(i, j int) bool {
		return output[i].CreatedAt.After(output[j].CreatedAt)
	})
	return output
}
