package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestResourceDatasetTypeProperties struct {
	AdditionalHeaders *map[string]interface{} `json:"additionalHeaders,omitempty"`
	PaginationRules   *map[string]interface{} `json:"paginationRules,omitempty"`
	RelativeURL       *interface{}            `json:"relativeUrl,omitempty"`
	RequestBody       *interface{}            `json:"requestBody,omitempty"`
	RequestMethod     *interface{}            `json:"requestMethod,omitempty"`
}
