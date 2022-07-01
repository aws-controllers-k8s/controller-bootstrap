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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type templateVars struct {
	*metaVars
	AWSSDKGoVersion    string
	RuntimeVersion     string
	TestInfraCommitSHA string
}

var templateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate template files in an ACK service controller repository",
	RunE:  generateController,
}

// generateController creates the initial directories and files for a service controller
// repository by rendering go template files.
// TODO: When a controller is already existing, then this method only updates the project
// description files.
func generateController(cmd *cobra.Command, args []string) error {
	if err := validateArgs(); err != nil {
		return err
	}
	ctx, cancel := contextWithSigterm(context.Background())
	defer cancel()
	if err := ensureSDKRepo(ctx, defaultCacheACKDir, optRefreshCache); err != nil {
		return err
	}

	svcVars, err := getServiceResources()
	if err != nil {
		return err
	}
	tplVars := &templateVars{
		metaVars:           svcVars,
		AWSSDKGoVersion:    optAWSSDKGoVersion,
		RuntimeVersion:     optRuntimeVersion,
		TestInfraCommitSHA: optTestInfraCommitSHA,
	}

	var tplPaths []string
	tplPaths, err = setTemplatePaths(tplPaths)
	if err != nil {
		return err
	}

	err = renderTemplateFiles(tplPaths, tplVars)
	if err != nil {
		return err
	}
	return nil
}

func validateArgs() error {
	if optServiceAlias == "" {
		return fmt.Errorf("please specify the AWS service alias for the service controller to generate")
	}
	if optRuntimeVersion == "" {
		return fmt.Errorf("please specify the aws-controllers-k8s/runtime version to generate the service controller")
	}
	if optAWSSDKGoVersion == "" {
		return fmt.Errorf("please specify the aws-sdk-go version to generate the service controller")
	}
	if optOutputPath == "" {
		return fmt.Errorf("please specify the output path to generate the service controller")
	}
	if optTestInfraCommitSHA == "" {
		return fmt.Errorf("please specify the aws-controllers-k8s/test-infra commit SHA to generate the service controller")
	}
	return nil
}

// Append the template files inside the templates directory to tplPaths.
func setTemplatePaths(tplPaths []string) ([]string, error) {
	err := filepath.Walk(defaultTemplatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			tplPaths = append(tplPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tplPaths, nil
}

// Loop over the template files from the templates directory and
// render them in the output path
func renderTemplateFiles(tplPaths []string, tplVars *templateVars) error {
	for _, tplPath := range tplPaths {
		tmp, err := template.ParseFiles(tplPath)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err = tmp.Execute(&buf, tplVars); err != nil {
			return err
		}
		file := strings.TrimPrefix(tplPath, defaultTemplatesDir)
		file = strings.TrimSuffix(file, ".tpl")
		outPath := filepath.Join(optOutputPath, file)
		if optDryRun {
			fmt.Printf("============================= %s ============================= \n", outPath)
			fmt.Println(strings.TrimSpace(buf.String()))
			continue
		}
		outDir := filepath.Dir(outPath)
		if _, err = ensureDir(outDir); err != nil {
			return err
		}
		if err = ioutil.WriteFile(outPath, buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}
