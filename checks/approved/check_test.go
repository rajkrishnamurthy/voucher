package approved

import (
	"context"
	"testing"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/repository"
	r "github.com/Shopify/voucher/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApprovedCheck(t *testing.T) {
	ctx := context.Background()
	imageData, err := voucher.NewImageData("gcr.io/voucher-test-project/apps/staging/voucher-internal@sha256:73d506a23331fce5cb6f49bfb4c27450d2ef4878efce89f03a46b27372a88430")
	require.NoErrorf(t, err, "failed to get ImageData: %s", err)
	buildDetail := r.BuildDetail{RepositoryURL: "https://github.com/Shopify/voucher-internal", Commit: "efgh6543"}
	commitURL := "https://github.com/Shopify/voucher-internal/commit/efgh6543"

	cases := []struct {
		name                 string
		defaultBranchCommits []r.CommitRef
		isSigned             bool
		status               string
		pullRequest          r.PullRequest
		shouldPass           bool
	}{
		{
			name:                 "Should pass",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}},
			shouldPass:           true,
		},
		{
			name:                 "Not built off default branch",
			defaultBranchCommits: []r.CommitRef{{URL: "otherCommit"}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}},
			shouldPass:           false,
		},
		{
			name:                 "Commit not signed",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             false,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}},
			shouldPass:           false,
		},
		{
			name:                 "Commit not a merge commit",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "SUCCESS",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: "otherURL"}},
			shouldPass:           false,
		},
		{
			name:                 "CI check not successful",
			defaultBranchCommits: []r.CommitRef{{URL: commitURL}},
			isSigned:             true,
			status:               "FAILURE",
			pullRequest:          r.PullRequest{IsMerged: true, MergeCommit: r.CommitRef{URL: commitURL}},
			shouldPass:           false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			commit := r.Commit{
				URL:                    commitURL,
				Status:                 testCase.status,
				IsSigned:               testCase.isSigned,
				AssociatedPullRequests: []r.PullRequest{testCase.pullRequest},
			}
			defaultBranch := r.Branch{Name: "production", CommitRefs: testCase.defaultBranchCommits}

			metadataClient := new(voucher.MockMetadataClient)
			metadataClient.On("GetBuildDetail", ctx, imageData).Return(buildDetail, nil)

			repositoryClient := new(repository.MockClient)
			repositoryClient.On("GetCommit", ctx, buildDetail).Return(commit, nil)
			repositoryClient.On("GetDefaultBranch", ctx, buildDetail).Return(defaultBranch, nil)

			orgCheck := new(check)
			orgCheck.SetMetadataClient(metadataClient)
			orgCheck.SetRepositoryClient(repositoryClient)

			status, err := orgCheck.Check(ctx, imageData)

			assert.NoErrorf(t, err, "check failed with error: %s", err)
			assert.Equal(t, testCase.shouldPass, status)
		})
	}
}
