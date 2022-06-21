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

	"github.com/spf13/cobra"
)

const (
	appName      = "controller-bootstrap"
	appShortDesc = "A bootstrap tool to initialize an ACK service controller repository"
)

var (
	optServiceAlias       string
	optRuntimeVersion     string
	optAWSSDKGoVersion    string
	optDryRun             bool
	optExistingController bool
	optOutputPath         string
	optModelName          string
)

// rootCmd represents the base command when called without any subcommands
// placeholder for cobra description
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: appShortDesc,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&optServiceAlias, "aws-service-alias", "s", "", "Supplied AWS service alias",
	)
	rootCmd.PersistentFlags().StringVarP(
		&optRuntimeVersion, "ack-runtime-version", "r", "", "Version of aws-controllers-k8s/runtime",
	)
	rootCmd.PersistentFlags().StringVarP(
		&optAWSSDKGoVersion, "aws-sdk-go-version", "v", "", "Version of github.com/aws/aws-sdk-go used to infer service metadata and resources",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&optDryRun, "dry-run", "d", false, "Optional: if true, output files to stdout",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&optExistingController, "existing-service-controller", "e", false, "Optional: if true, update the existing service controller",
	)
	rootCmd.PersistentFlags().StringVarP(
		&optOutputPath, "output", "o", "", "Path to ACK service controller directory to bootstrap",
	)
	rootCmd.PersistentFlags().StringVarP(
		&optModelName, "service-model-name", "m", "", "Optional: service model name of the supplied service alias",
	)
	rootCmd.MarkPersistentFlagRequired("aws-service-alias")
	rootCmd.MarkPersistentFlagRequired("ack-runtime-version")
	rootCmd.MarkPersistentFlagRequired("aws-sdk-go-version")
	rootCmd.MarkPersistentFlagRequired("output")
	rootCmd.AddCommand(templateCmd)
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
