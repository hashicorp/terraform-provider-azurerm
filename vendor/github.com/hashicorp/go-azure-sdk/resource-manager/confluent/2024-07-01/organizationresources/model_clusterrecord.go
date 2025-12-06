package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterRecord struct {
	DisplayName *string              `json:"display_name,omitempty"`
	Id          *string              `json:"id,omitempty"`
	Kind        *string              `json:"kind,omitempty"`
	Metadata    *MetadataEntity      `json:"metadata,omitempty"`
	Spec        *ClusterSpecEntity   `json:"spec,omitempty"`
	Status      *ClusterStatusEntity `json:"status,omitempty"`
}
