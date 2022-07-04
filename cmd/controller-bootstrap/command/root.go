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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	appName      = "controller-bootstrap"
	appShortDesc = "A bootstrap tool to initialize an ACK service controller repository"
)

var (
	optRuntimeVersion     string
	optAWSSDKGoVersion    string
	optTestInfraCommitSHA string
	optModelName          string
	optRefreshCache       bool
	optServiceAlias       string
	optDryRun             bool
	optOutputPath         string
	sdkDir                string
	defaultCacheACKDir    string
	defaultTemplatesDir   string
)

// rootCmd represents the base command when called without any subcommands
// placeholder for cobra description
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appShortDesc,
}

func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("unable to determine $HOME: %s\n", err)
		os.Exit(1)
	}
	defaultCacheACKDir = filepath.Join(hd, ".cache", "aws-controllers-k8s")

	cd, err := os.Getwd()
	if err != nil {
		fmt.Printf("unable to determine current working directory: %s\n", err)
		os.Exit(1)
	}
	defaultTemplatesDir = filepath.Join(cd, "templates")

	templateCmd.PersistentFlags().StringVar(
		&optRuntimeVersion, "ack-runtime-version", "", "Version of aws-controllers-k8s/runtime",
	)
	templateCmd.PersistentFlags().StringVar(
		&optAWSSDKGoVersion, "aws-sdk-go-version", "", "Version of github.com/aws/aws-sdk-go used to infer service metadata and resources",
	)
	templateCmd.PersistentFlags().StringVar(
		&optTestInfraCommitSHA, "test-infra-commit-sha", "", "Commit SHA of aws-controllers-k8s/test-infra",
	)
	templateCmd.PersistentFlags().StringVar(
		&optModelName, "model-name", "", "Optional: service model name of the corresponding service alias",
	)
	templateCmd.PersistentFlags().BoolVar(
		&optRefreshCache, "refresh-cache", true, "Optional: if true, and aws-sdk-go repo is already cloned, will git pull the latest aws-sdk-go commit",
	)
	rootCmd.PersistentFlags().StringVar(
		&optServiceAlias, "aws-service-alias", "", "AWS service alias",
	)
	rootCmd.PersistentFlags().BoolVar(
		&optDryRun, "dry-run", false, "Optional: if true, output files to stdout",
	)
	rootCmd.PersistentFlags().StringVar(
		&optOutputPath, "output-path", "", "Path to ACK service controller directory to bootstrap",
	)
	templateCmd.MarkPersistentFlagRequired("ack-runtime-version")
	templateCmd.MarkPersistentFlagRequired("aws-sdk-go-version")
	templateCmd.MarkPersistentFlagRequired("test-infra-commit-sha")
	rootCmd.MarkPersistentFlagRequired("aws-service-alias")
	rootCmd.MarkPersistentFlagRequired("output-path")
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(updateCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
