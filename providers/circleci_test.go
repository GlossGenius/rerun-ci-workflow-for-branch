package providers

import (
	"github.com/google/go-cmp/cmp"
	"github.com/grezar/go-circleci"
	"testing"
	"time"
)

func TestRemoveNils(t *testing.T) {
	pipeline1 := &circleci.Pipeline{}
	pipeline2 := &circleci.Pipeline{}
	pipelines := []*circleci.Pipeline{
		pipeline1, nil, pipeline2, nil,
	}
	expected := []*circleci.Pipeline{
		pipeline1, pipeline2,
	}
	actual := removeNils(pipelines)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("lists of pipelines differ.\nexpected: %v\nactual: %v\bdiff: %v", expected, actual, diff)
	}
}

func TestSortByCreateDate(t *testing.T) {
	now := time.Now()
	pipeline1 := &circleci.Pipeline{
		CreatedAt: now.Add(1 * time.Hour),
	}
	pipeline2 := &circleci.Pipeline{
		CreatedAt: now,
	}
	pipeline3 := &circleci.Pipeline{
		CreatedAt: now.Add(2 * time.Hour),
	}
	pipelines := []*circleci.Pipeline{
		pipeline1,
		pipeline2,
		pipeline3,
	}
	expected := []*circleci.Pipeline{
		pipeline3, pipeline1, pipeline2,
	}
	actual := sortByCreateDate(pipelines)
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("lists of pipelines not sorted as expected.\nexpected: %v\nactual: %v\bdiff: %v", expected, actual, diff)
	}
}
