// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

func SearchDatasourceStorageConnectionString(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	// TODO: Does it work if it's unknown during plan??
	if value == "" {
		errors = append(errors, fmt.Errorf("property `%s` must not be empty", k))
		return warnings, errors
	}

	isResourceId := strings.Contains(value, "ResourceId=")
	isAccountKey := strings.Contains(value, "DefaultEndpointsProtocol=") && strings.Contains(value, "AccountName=") && strings.Contains(value, "AccountKey=")
	isSAS := strings.Contains(value, "BlobEndpoint=") && strings.Contains(value, "SharedAccessSignature=")
	isContainerSAS := strings.Contains(value, "ContainerSharedAccessUri=")

	if !isResourceId && !isAccountKey && !isSAS && !isContainerSAS {
		errors = append(errors, fmt.Errorf("property `%s` must be in one of the following formats: "+
			"`ResourceId=[resource ID string];` or "+
			"`DefaultEndpointsProtocol=https;AccountName=[your storage account];AccountKey=[your account key];` or "+
			"`BlobEndpoint=[your blob endpoint URI];SharedAccessSignature=[your account shared access token];` or "+
			"`ContainerSharedAccessUri=[your container shared access token with read and list permissions]`", k))
	}

	return warnings, errors
}
