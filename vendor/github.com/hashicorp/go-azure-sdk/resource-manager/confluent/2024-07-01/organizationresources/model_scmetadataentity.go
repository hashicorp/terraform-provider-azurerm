package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SCMetadataEntity struct {
	CreatedTimestamp *string `json:"createdTimestamp,omitempty"`
	DeletedTimestamp *string `json:"deletedTimestamp,omitempty"`
	ResourceName     *string `json:"resourceName,omitempty"`
	Self             *string `json:"self,omitempty"`
	UpdatedTimestamp *string `json:"updatedTimestamp,omitempty"`
}
