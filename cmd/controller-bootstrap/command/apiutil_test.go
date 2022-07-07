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
	"go/build"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var (
	goPath               = build.Default.GOPATH
	defaultRootDirectory = filepath.Join(goPath, "src/github.com/aws-controllers-k8s")
	serviceAPIVersion    = "0000-00-00"
)

// integration/ end-to-end test - this will run the command line
// self-contain the unit test - make a copy of the api-2.json file - 2 assertions
func Test_modelAPI(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	apiPath := filepath.Join(defaultRootDirectory, "controller-bootstrap", "pkg", "testdata", "models", "apis")
	tests := []struct {
		ServiceModelName    string
		ServiceID           string
		ServiceAbbreviation string
		ServiceFullName     string
		CRDNames            []string
	}{
		{
			ServiceModelName:    "eks",
			ServiceID:           "EKS",
			ServiceAbbreviation: "Amazon EKS",
			ServiceFullName:     "Amazon Elastic Kubernetes Service",
			CRDNames:            []string{"Addon", "Cluster", "FargateProfile", "Nodegroup"},
		},
		{
			ServiceModelName:    "rds",
			ServiceID:           "RDS",
			ServiceAbbreviation: "Amazon RDS",
			ServiceFullName:     "Amazon Relational Database Service",
			CRDNames: []string{"CustomAvailabilityZone", "DBCluster", "DBClusterEndpoint",
				"DBClusterParameterGroup", "DBClusterSnapshot", "DBInstance", "DBInstanceReadReplica",
				"DBParameterGroup", "DBProxy", "DBSecurityGroup", "DBSnapshot",
				"DBSubnetGroup", "EventSubscription", "GlobalCluster", "OptionGroup"},
		},
	}
	h := newSDKHelper()
	for _, test := range tests {
		apiFile := filepath.Join(apiPath, test.ServiceModelName, serviceAPIVersion, "api-2.json")
		svcVars, err := h.modelAPI(apiFile)
		require.NoError(err)
		assert.Equal(test.ServiceID, svcVars.ServiceID)
		assert.Equal(test.ServiceAbbreviation, svcVars.ServiceAbbreviation)
		assert.Equal(test.ServiceFullName, svcVars.ServiceFullName)
		assert.Equal(test.CRDNames, svcVars.CRDNames)
	}
}
