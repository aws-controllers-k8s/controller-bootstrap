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
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate template files in an ACK service controller repository",
	RunE:  generateTemplates,
}

// todo: generateTemplates renders the template files and directories in an ACK service controller repository
func generateTemplates(cmd *cobra.Command, args []string) error {
	return nil
}
