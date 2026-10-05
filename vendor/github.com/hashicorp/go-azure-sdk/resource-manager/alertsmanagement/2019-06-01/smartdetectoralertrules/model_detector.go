package smartdetectoralertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Detector struct {
	Description            *string                 `json:"description,omitempty"`
	Id                     string                  `json:"id"`
	ImagePaths             *[]string               `json:"imagePaths,omitempty"`
	Name                   *string                 `json:"name,omitempty"`
	Parameters             *map[string]interface{} `json:"parameters,omitempty"`
	SupportedResourceTypes *[]string               `json:"supportedResourceTypes,omitempty"`
}
