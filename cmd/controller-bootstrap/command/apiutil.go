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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"
)

type metaVars struct {
	ServiceID           string
	ServicePackageName  string
	ServiceModelName    string
	ServiceAbbreviation string
	ServiceFullName     string
	CRDNames            []string
}

// SDKHelper is a helper struct to work with the aws-sdk-go models and
// API model loader
// TODO: move SDKHelper struct and its corresponding methods to aws-controllers-k8s/pkg repository
type SDKHelper struct {
	loader *awssdkmodel.Loader
	// Default is set by `FirstAPIVersion`
	apiVersion string
}

var (
	ErrInvalidVersionDirectory = errors.New(
		"expected to find only directories in api model directory but found non-directory",
	)
	ErrNoValidVersionDirectory = errors.New(
		"no valid version directories found",
	)
	ErrServiceAPIFileNotFound = errors.New(
		"unable to find the supplied service's api-2.json file, please re-try specifying the service model name",
	)
)

// getServiceResources infers aws-sdk-go to fetch the service metadata and custom resource names
func getServiceResources() (*metaVars, error) {
	serviceModelName := strings.ToLower(optModelName)
	if optModelName == "" {
		serviceModelName = strings.ToLower(optServiceAlias)
	}
	h := newSDKHelper()
	modelPath, err := h.findModelPath(serviceModelName)
	if err != nil {
		return nil, ErrServiceAPIFileNotFound
	}
	svcVars, err := h.modelAPI(modelPath)
	if err != nil {
		return nil, err
	}
	return svcVars, nil
}

// newSDKHelper returns a new SDKHelper struct
func newSDKHelper() *SDKHelper {
	return &SDKHelper{
		loader: &awssdkmodel.Loader{
			BaseImport:            sdkDir,
			IgnoreUnsupportedAPIs: true,
		},
	}
}

// findModelPath returns the path to the supplied service's api-2.json file
func (h *SDKHelper) findModelPath(
	serviceModelName string,
) (string, error) {
	if h.apiVersion == "" {
		apiVersion, err := h.firstAPIVersion(serviceModelName)
		if err != nil {
			return "", err
		}
		h.apiVersion = apiVersion
	}
	versionPath := filepath.Join(
		sdkDir, "models", "apis", serviceModelName, h.apiVersion,
	)
	modelPath := filepath.Join(versionPath, "api-2.json")
	return modelPath, nil
}

// firstAPIVersion returns the first found API version for a service API.
// (e.h. "2012-10-03")
func (h *SDKHelper) firstAPIVersion(serviceModelName string) (string, error) {
	versions, err := h.getAPIVersions(serviceModelName)
	if err != nil {
		return "", err
	}
	sort.Strings(versions)
	return versions[0], nil
}

// getAPIVersions returns the list of API Versions found in a service directory.
func (h *SDKHelper) getAPIVersions(serviceModelName string) ([]string, error) {
	apiPath := filepath.Join(sdkDir, "models", "apis", serviceModelName)
	versionDirs, err := ioutil.ReadDir(apiPath)
	if err != nil {
		return nil, err
	}
	versions := []string{}
	for _, f := range versionDirs {
		version := f.Name()
		fp := filepath.Join(apiPath, version)
		fi, err := os.Lstat(fp)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("found %s: %v", version, ErrInvalidVersionDirectory)
		}
		versions = append(versions, version)
	}
	if len(versions) == 0 {
		return nil, ErrNoValidVersionDirectory
	}
	return versions, nil
}

// modelAPI returns the populated metaVars struct with the service metadata
// and custom resource names extracted from the aws-sdk-go model API object
func (h *SDKHelper) modelAPI(modelPath string) (*metaVars, error) {
	// loads the API model file(s) and returns the map of API package
	apis, err := h.loader.Load([]string{modelPath})
	if err != nil {
		return nil, err
	}
	// apis is a map, keyed by the service package name, of pointers
	// to aws-sdk-go model API objects
	for _, api := range apis {
		_ = api.ServicePackageDoc()
		svcVars := serviceMetaVars(api)
		return svcVars, nil
	}
	return nil, err
}

// serviceMetaVars returns a metaVars struct populated with metadata
// and custom resource names for the supplied AWS service
func serviceMetaVars(api *awssdkmodel.API) *metaVars {
	return &metaVars{
		ServicePackageName:  strings.ToLower(optServiceAlias),
		ServiceID:           api.Metadata.ServiceID,
		ServiceModelName:    strings.ToLower(optModelName),
		ServiceAbbreviation: api.Metadata.ServiceAbbreviation,
		ServiceFullName:     api.Metadata.ServiceFullName,
		CRDNames:            getCRDNames(api),
	}
}

// getCRDNames returns the CustomResource names present in the api.
// CustomResource names are created by dropping the prefix "Create" from
// all the operation names that start with prefix "Create".
// Operations with prefix "CreateBatch" are ignored.
func getCRDNames(api *awssdkmodel.API) []string {
	var crdNames []string
	pluralize := pluralize.NewClient()
	for _, opName := range api.OperationNames() {
		if strings.HasPrefix(opName, "CreateBatch") {
			continue
		}
		if strings.HasPrefix(opName, "Create") {
			resName := strings.TrimPrefix(opName, "Create")
			if pluralize.IsSingular(resName) {
				crdNames = append(crdNames, resName)
			}
		}
	}
	return crdNames
}
