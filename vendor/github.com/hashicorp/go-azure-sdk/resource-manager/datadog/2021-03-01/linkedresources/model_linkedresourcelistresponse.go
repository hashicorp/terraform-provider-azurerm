package linkedresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedResourceListResponse struct {
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]LinkedResource `json:"value,omitempty"`
}
