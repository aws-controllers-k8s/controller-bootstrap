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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
	basePath string
	// Default is set by `latestAPIVersion`
	apiVersion string
}

// API holds all the shapes defined in the <service>.json
// api model file provided by aws-sdk-go-v2
type API struct {
	Shapes map[string]Shape `json:"shapes"`
}

// Shape contains the definition of a resource, field,
// operation, service, etc.
type Shape struct {
	Type       string
	Traits     map[string]interface{}
	MemberRefs map[string]*ShapeRef `json:"members"`
	MemberRef  *ShapeRef            `json:"member"`
	KeyRef     ShapeRef             `json:"key"`
	ValueRef   ShapeRef             `json:"value"`
	InputRef   ShapeRef             `json:"input"`
	OutputRef  ShapeRef             `json:"output"`
	ErrorRefs  []ShapeRef           `json:"errors"`
}

// ShapeRef defines the usage of a shape within the API
type ShapeRef struct {
	ShapeName string `json:"target"`
	Traits    map[string]interface{}
}

var (
	ErrInvalidVersionDirectory = errors.New(
		"expected to find only directories in api model directory but found non-directory",
	)
	ErrNoValidVersionDirectory = errors.New(
		"no valid version directories found",
	)
	ErrServiceAPIFileNotFound = errors.New(
		"unable to find the supplied service's api model file, please re-try specifying the service model name",
	)
)

// getServiceResources infers aws-sdk-go to fetch the service metadata and custom resource names
func getServiceResources() (*metaVars, error) {
	serviceModelName := strings.ToLower(optModelName)
	if optModelName == "" {
		serviceModelName = strings.ToLower(optServiceAlias)
	}
	h := newSDKHelper()
	modelPath, err := h.ModelAndDocsPath(serviceModelName)
	if err != nil {
		return nil, err
	}

	svcVars, err := loadAPI(modelPath)
	if err != nil {
		return nil, err
	}

	return svcVars, nil
}

// newSDKHelper returns a new SDKHelper struct
func newSDKHelper() *SDKHelper {
	return &SDKHelper{
		basePath: sdkDir,
	}
}

// getCRDNames returns the CustomResource names present in the api.
// CustomResource names are created by dropping the prefix "Create" from
// all the operation names that start with prefix "Create".
// Operations with prefix "CreateBatch" are ignored.
func getCRDNames(operations []string) []string {
	var crdNames []string
	pluralize := pluralize.NewClient()
	sort.Strings(operations)

	for _, opName := range operations {
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

// Loads v2 API model and parse Service Metadata and derives CRD names from API model operations.
func loadAPI(modelPath string) (*metaVars, error) {
	file, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var customAPI API
	err = json.Unmarshal(file, &customAPI)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling file: %v", err)
	}

	svcVars := metaVars{}
	operations := make([]string, 0)
	for shapeName, shape := range customAPI.Shapes {
		switch shape.Type {
		case "service":
			serviceId, ok := shape.Traits["aws.api#service"].(map[string]interface{})["sdkId"]
			if !ok {
				return nil, errors.New("service id not found")
			}

			serviceTitle, ok := shape.Traits["smithy.api#title"].(string)
			if !ok {
				return nil, errors.New("Service title not found")
			}

			svcVars.ServiceID = serviceId.(string)
			svcVars.ServiceFullName = serviceTitle
			svcVars.ServiceAbbreviation = serviceTitle

		case "operation":
			name, err := removeShapeNamePrefix(shapeName)
			if err != nil {
				return nil, err
			}

			operations = append(operations, name)

		default:
			continue
		}
	}

	svcVars.CRDNames = getCRDNames(operations)

	return &svcVars, nil
}

// removeShapeNamePrefix removes the prefix from the shapeName.
// The prefix format of a shape in v2 is com.amazonaws.<serviceAlias>#shapeName
func removeShapeNamePrefix(name string) (string, error) {
	temp := strings.Split(name, "#")
	if len(temp) != 2 {
		return "", fmt.Errorf("%s shape name is not formatted correctly, expected format: <url>:<shapeName>", name)
	}
	newName := temp[1]

	return newName, nil
}

// extractServiceAlias extracts the service alias from a shapeName
// (see removeShapeNamePrefix)
func extractServiceAlias(name string) string {
	temp := strings.Split(name, ".")
	anotherTemp := strings.Split(temp[len(temp)-1], "#")
	if len(anotherTemp) != 2 {
		return ""
	}
	alias := anotherTemp[0]
	return alias
}

// ModelAndDocsPath returns two string paths to the supplied service's API and
// doc JSON files
func (h *SDKHelper) ModelAndDocsPath(serviceModelName string) (string, error) {
	modelPath := filepath.Join(
		h.basePath,
		"codegen",
		"sdk-codegen",
		"aws-models",
		fmt.Sprintf("%s.json", serviceModelName),
	)

	if _, err := os.Stat(modelPath); err == nil {
		return modelPath, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return "", ErrServiceAPIFileNotFound
	} else {
		return "", err
	}
}
