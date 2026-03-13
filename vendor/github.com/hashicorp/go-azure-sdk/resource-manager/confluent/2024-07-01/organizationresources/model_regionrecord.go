package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionRecord struct {
	Id         *string           `json:"id,omitempty"`
	Kind       *string           `json:"kind,omitempty"`
	Properties *RegionProperties `json:"properties,omitempty"`
}
