package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportedStacksCollection struct {
	NextLink *string                   `json:"nextLink,omitempty"`
	Value    *[]SupportedStackResource `json:"value,omitempty"`
}
