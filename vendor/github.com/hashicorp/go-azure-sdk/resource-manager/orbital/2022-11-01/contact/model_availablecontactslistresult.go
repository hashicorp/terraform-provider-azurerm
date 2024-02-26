package contact

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableContactsListResult struct {
	NextLink *string              `json:"nextLink,omitempty"`
	Value    *[]AvailableContacts `json:"value,omitempty"`
}
