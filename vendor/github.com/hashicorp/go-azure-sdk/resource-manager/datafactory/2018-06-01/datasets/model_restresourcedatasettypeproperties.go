package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestResourceDatasetTypeProperties struct {
	AdditionalHeaders *map[string]string `json:"additionalHeaders,omitempty"`
	PaginationRules   *map[string]string `json:"paginationRules,omitempty"`
	RelativeURL       *string            `json:"relativeUrl,omitempty"`
	RequestBody       *string            `json:"requestBody,omitempty"`
	RequestMethod     *string            `json:"requestMethod,omitempty"`
}
