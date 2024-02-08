package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceAssociationLinksListResult struct {
	NextLink *string                   `json:"nextLink,omitempty"`
	Value    *[]ServiceAssociationLink `json:"value,omitempty"`
}
