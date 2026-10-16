package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsDeploymentStorage struct {
	Authentication *FunctionsDeploymentStorageAuthentication `json:"authentication,omitempty"`
	Type           *FunctionsDeploymentStorageType           `json:"type,omitempty"`
	Value          *string                                   `json:"value,omitempty"`
}
