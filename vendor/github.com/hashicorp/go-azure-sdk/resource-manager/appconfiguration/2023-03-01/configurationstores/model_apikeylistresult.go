package configurationstores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKeyListResult struct {
	NextLink *string   `json:"nextLink,omitempty"`
	Value    *[]ApiKey `json:"value,omitempty"`
}
