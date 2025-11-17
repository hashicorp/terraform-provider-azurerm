package linkedservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureSynapseArtifactsLinkedServiceTypeProperties struct {
	Authentication      *interface{} `json:"authentication,omitempty"`
	Endpoint            interface{}  `json:"endpoint"`
	WorkspaceResourceId *interface{} `json:"workspaceResourceId,omitempty"`
}
