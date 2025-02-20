package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CalloutPoliciesList struct {
	NextLink *string          `json:"nextLink,omitempty"`
	Value    *[]CalloutPolicy `json:"value,omitempty"`
}
