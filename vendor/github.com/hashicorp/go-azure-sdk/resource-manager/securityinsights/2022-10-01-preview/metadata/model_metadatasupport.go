package metadata

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataSupport struct {
	Email *string     `json:"email,omitempty"`
	Link  *string     `json:"link,omitempty"`
	Name  *string     `json:"name,omitempty"`
	Tier  SupportTier `json:"tier"`
}
