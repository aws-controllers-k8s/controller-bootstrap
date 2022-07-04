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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var projectDescriptionFiles = []string{
	"CODE_OF_CONDUCT.md.tpl",
	"CONTRIBUTING.md.tpl",
	"GOVERNANCE.md.tpl",
	"LICENSE.tpl",
	"NOTICE.tpl",
	"READ_BEFORE_COMMIT.md.tpl",
	"SECURITY.md.tpl",
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update project description files in an existing ACK service controller repository",
	RunE:  updateController,
}

// updateController only updates the project description files for an
// existing service controller repository of the supplied service alias.
func updateController(cmd *cobra.Command, args []string) error {
	if optServiceAlias == "" {
		return ErrServiceAliasNotFound
	}
	if optOutputPath == "" {
		return ErrOutputPathNotFound
	}
	if !controllerExists() {
		return fmt.Errorf("the service controller repository for the supplied AWS service alias does not exist, please run the generate command for a new service controller")
	}

	// Loop over the template project description files and
	// render them in the existing service controller repository
	for _, tplPath := range projectDescriptionFiles {
		tplPath = filepath.Join(defaultTemplatesDir, tplPath)
		tmp, err := template.ParseFiles(tplPath)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err = tmp.Execute(&buf, nil); err != nil {
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
		if err = ioutil.WriteFile(outPath, buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}
