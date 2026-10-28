/*
 * SPDX-License-Identifier: GPL-3.0
 * ChunCord integration maintained by chunma-san.
 */
package main

import "testing"

func TestLatestRevision(t *testing.T) {
	testCases := []struct {
		name     string
		release  GithubRelease
		expected string
	}{
		{
			name:     "current workflow release: clean title with target_commitish",
			release:  GithubRelease{Name: "ChunCord", TagName: "chuncord-0123456789ab", TargetCommitish: "0123456789abcdef0123456789abcdef01234567"},
			expected: "0123456789abcdef0123456789abcdef01234567",
		},
		{
			name:     "commitish is a branch name: fall back to title suffix",
			release:  GithubRelease{Name: "Equicord 0123456", TargetCommitish: "custom"},
			expected: "0123456",
		},
		{
			name:     "legacy release without commitish",
			release:  GithubRelease{Name: "Equicord 0123456"},
			expected: "0123456",
		},
		{
			name:     "short commitish is accepted",
			release:  GithubRelease{Name: "ChunCord", TargetCommitish: "0123456"},
			expected: "0123456",
		},
		{
			name:     "uppercase commitish is normalized",
			release:  GithubRelease{Name: "ChunCord", TargetCommitish: "0123456789ABCDEF0123456789ABCDEF01234567"},
			expected: "0123456789abcdef0123456789abcdef01234567",
		},
		{
			name:     "clean title without any revision source stays unchanged",
			release:  GithubRelease{Name: "ChunCord"},
			expected: "ChunCord",
		},
	}
	for _, testCase := range testCases {
		if actual := latestRevision(&testCase.release); actual != testCase.expected {
			t.Errorf("%s: latestRevision() = %q, want %q", testCase.name, actual, testCase.expected)
		}
	}
}
