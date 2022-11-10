package providers

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v48/github"
	"testing"
)

func TestFilterForPrefix(t *testing.T) {
	release1 := &github.Branch{Name: github.String("release/1")}
	release2 := &github.Branch{Name: github.String("release/2")}
	hotfix := &github.Branch{Name: github.String("hotfix/123")}
	integration := &github.Branch{Name: github.String("integration")}
	branches := []*github.Branch{
		release1,
		release2,
		hotfix,
		integration,
	}
	type testEntry struct {
		prefix   string
		expected []*github.Branch
	}
	testEntries := []testEntry{
		{
			prefix:   "no-overlap",
			expected: []*github.Branch{},
		},
		{
			prefix: "release/",
			expected: []*github.Branch{
				release1,
				release2,
			},
		},
		{
			prefix: "release/1",
			expected: []*github.Branch{
				release1,
			},
		},
		{
			prefix:   "",
			expected: branches,
		},
	}
	for _, entry := range testEntries {
		got := filterForPrefix(branches, entry.prefix)
		if diff := cmp.Diff(got, entry.expected); diff != "" {
			t.Errorf("expected: %v\n, got: %v\n, diff: %v\n", entry.expected, got, diff)
		}
	}
}

func TestGetBranchWithPrefix(t *testing.T) {
	release1 := &github.Branch{Name: github.String("release/1")}
	release2 := &github.Branch{Name: github.String("release/2")}
	hotfix := &github.Branch{Name: github.String("hotfix/123")}
	integration := &github.Branch{Name: github.String("integration")}
	branches := []*github.Branch{
		release1,
		release2,
		hotfix,
		integration,
	}
	type testEntry struct {
		prefix         string
		expectedBranch *github.Branch
		expectedError  error
	}
	testEntries := []testEntry{
		{
			prefix:         "no-overlap",
			expectedBranch: nil,
			expectedError:  nil,
		},
		{
			prefix:         "release/",
			expectedBranch: nil,
			expectedError:  &cantDisambiguateBranchesError{2},
		},
		{
			prefix:         "release/1",
			expectedBranch: release1,
			expectedError:  nil,
		},
		{
			prefix:         "",
			expectedError:  &cantDisambiguateBranchesError{len(branches)},
			expectedBranch: nil,
		},
	}
	for _, entry := range testEntries {
		gotBranch, gotError := getBranchWithPrefix(branches, entry.prefix)
		if diff := cmp.Diff(entry.expectedBranch, gotBranch); diff != "" {
			t.Errorf("branches differ:\nexpected: %v\ngot: %v\ndiff: %s", entry.expectedBranch, gotBranch, diff)
		}
		if !errors.Is(gotError, entry.expectedError) {
			t.Errorf("errors differ:\nexpected: %s\ngot: %s\n", entry.expectedError, gotError)
		}
	}
}
