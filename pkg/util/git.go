// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

//TODO: move this file to aws-controllers-k8s/pkg repository after initial implementation of the controller-bootstrap

// LoadRepository loads a repository from the local file system.
// TODO: load repository into a memory filesystem - needs go1.16
// migration or use something like https://github.com/spf13/afero
func LoadRepository(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

// HasTag checks if a tag exists in the local repository.
func HasTag(path string, tag string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--verify", fmt.Sprintf("refs/tags/%s", tag))
	return cmd.Run() == nil
}

// CloneRepository clones a git repository into a given directory.
//
// Calling his function is equivalent to executing `git clone $repositoryURL $path`
func CloneRepository(ctx context.Context, path, repositoryURL string) error {
	_, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL:      repositoryURL,
		Progress: nil,
		// Clone and fetch all tags
		Tags: git.AllTags,
	})
	return err
}

// FetchRepositoryTag fetches a single tag from the remote repository.
//
// Equivalent to: git -C $path fetch origin tag $tag
func FetchRepositoryTag(ctx context.Context, path string, tag string) error {
	cmd := exec.CommandContext(ctx, "git", "-C", path, "fetch", "origin", "tag", tag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

// getRepositoryTagRef returns the git reference (commit hash) of a given tag.
// NOTE: It is not possible to checkout a tag without knowing it's reference.
//
// Calling this function is equivalent to executing `git rev-list -n 1 $tagName`
func getRepositoryTagRef(repo *git.Repository, tagName string) (*plumbing.Reference, error) {
	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	for {
		tagRef, err := tagRefs.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error finding tag reference: %v", err)
		}
		if tagRef.Name().Short() == tagName {
			return tagRef, nil
		}
	}
	return nil, errors.New("tag reference not found")
}

// CheckoutRepositoryTag checkouts a repository tag by looking for the tag
// reference then calling the checkout function.
//
// Calling This function is equivalent to executing `git checkout tags/$tag`
func CheckoutRepositoryTag(path string, tag string) error {
	cmd := exec.Command("git", "-C", path, "checkout", fmt.Sprintf("tags/%s", tag), "-f")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}
