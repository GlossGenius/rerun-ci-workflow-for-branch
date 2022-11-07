package providers

import "fmt"

const maxBranches = 1

type cantDisambiguateBranchesError struct {
	actualBranches int
}

func (c *cantDisambiguateBranchesError) Error() string {
	return fmt.Sprintf("too many branches found. found %d branches, expected <= %d", c.actualBranches, maxBranches)
}

func (c *cantDisambiguateBranchesError) Is(target error) bool {
	// TODO(michel): since I'm new to custom error handling, check if there's a way to do this without down casting
	tgt, ok := target.(*cantDisambiguateBranchesError)
	if !ok {
		return false
	}
	return tgt.actualBranches == c.actualBranches
}
