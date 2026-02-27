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

package command

import (
	"context"
	"controller-bootstrap/pkg/util"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

//TODO: move this file to aws-controllers-k8s/pkg repository after initial implementation of the controller-bootstrap

const (
	sdkRepoURL             = "https://github.com/aws/aws-sdk-go-v2"
	defaultGitCloneTimeout = 180 * time.Second
	defaultGitFetchTimeout = 30 * time.Second
)

// controllerExists returns True if the service controller already exists, false otherwise
func controllerExists() bool {
	file := filepath.Join(optOutputPath, "cmd", "controller", "main.go")
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func contextWithSigterm(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	signalCh := make(chan os.Signal, 1)

	// recreate the context.CancelFunc
	cancelFunc := func() {
		signal.Stop(signalCh)
		cancel()
	}

	// notify on SIGINT or SIGTERM
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancelFunc
}

// ensureDir makes sure that a supplied directory exists and
// returns whether the directory already existed.
func ensureDir(fp string) (bool, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return false, os.MkdirAll(fp, os.ModePerm)
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("expected %s to be a directory", fp)
	}
	if !isDirWriteable(fp) {
		return true, fmt.Errorf("%s is not a writeable directory", fp)
	}

	return true, nil
}

// isDirWriteable returns true if the supplied directory path is writeable,
// false otherwise
func isDirWriteable(fp string) bool {
	testPath := filepath.Join(fp, "test")
	f, err := os.Create(testPath)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testPath)
	return true
}

// ensureSDKRepo ensures that we have a git clone'd copy of the aws-sdk-go
// repository, which we use model JSON files from. Upon successful return of
// this function, the sdkDir global variable will be set to the directory where
// the aws-sdk-go is found. It will also optionally fetch all the remote tags
// and checkout the given tag.
// EnsureRepo ensures that we have a git clone'd copy of the aws-sdk-go
// repository, which we use model JSON files from. Upon successful return of
// this function, the sdkDir global variable will be set to the directory where
// the aws-sdk-go is found. It will also optionally fetch all the remote tags
// and checkout the given tag.
func ensureSDKRepo(
	ctx context.Context,
	cacheDir string,
	// A boolean instructing ensureSDKRepo whether to fetch the remote tags from
	// the upstream repository
	fetchTags bool,
) error {
	var err error
	srcPath := filepath.Join(cacheDir, "src")
	if err = os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	// Clone repository if it doesn't exist
	sdkDir := filepath.Join(srcPath, "aws-sdk-go-v2")
	if _, err = os.Stat(sdkDir); os.IsNotExist(err) {
		cloneStart := time.Now()
		ctx, cancel := context.WithTimeout(ctx, defaultGitCloneTimeout)
		defer cancel()
		err = util.CloneRepository(ctx, sdkDir, sdkRepoURL)
		if err != nil {
			// See https://github.com/aws-controllers-k8s/community/issues/1642
			if errors.Is(err, context.DeadlineExceeded) {
				err = fmt.Errorf("%w: take too long to clone aws sdk repo, "+
					"please consider manually 'git clone %s' to cache dir %s", err, sdkRepoURL, sdkDir)
			}
			return fmt.Errorf("cannot clone repository: %v", err)
		}
		util.Tracef("git clone: %s\n", time.Since(cloneStart))
	} else {
		util.Tracef("git clone: skipped (cached)\n")
	}

	// Resolve the SDK version before fetching so we only fetch the tag we need
	sdkVersion := getSDKVersion()
	if err != nil {
		return err
	}
	sdkVersion = ensureSemverPrefix(sdkVersion)

	// Fetch the tag only if it doesn't already exist locally
	if fetchTags && !util.HasTag(sdkDir, sdkVersion) {
		fetchStart := time.Now()
		ctx, cancel := context.WithTimeout(ctx, defaultGitFetchTimeout)
		defer cancel()
		err = util.FetchRepositoryTag(ctx, sdkDir, sdkVersion)
		if err != nil {
			return fmt.Errorf("cannot fetch tag %s: %v", sdkVersion, err)
		}
		util.Tracef("git fetch tag %s: %s\n", sdkVersion, time.Since(fetchStart))
	} else if fetchTags {
		util.Tracef("git fetch: skipped (tag %s exists locally)\n", sdkVersion)
	}

	err = util.CheckoutRepositoryTag(sdkDir, sdkVersion)
	if err != nil {
		return fmt.Errorf("cannot checkout tag: %v", err)
	}
	return nil
}

// ensureSemverPrefix takes a semver string and tries to append the 'v'
// prefix if it's missing.
func ensureSemverPrefix(s string) string {
	// trim all leading 'v' runes (characters)
	s = strings.TrimLeft(s, "v")
	return fmt.Sprintf("v%s", s)
}

// getSDKVersion returns the github.com/aws/aws-sdk-go-v2 version to use
// from the --aws-sdk-go-version flag.
func getSDKVersion() string {
	return optAWSSDKGoVersion
}
