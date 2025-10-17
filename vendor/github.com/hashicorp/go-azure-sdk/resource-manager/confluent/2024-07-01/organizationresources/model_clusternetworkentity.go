package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterNetworkEntity struct {
	Environment  *string `json:"environment,omitempty"`
	Id           *string `json:"id,omitempty"`
	Related      *string `json:"related,omitempty"`
	ResourceName *string `json:"resource_name,omitempty"`
}
