package agents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KnowledgeGraphConfiguration struct {
	Identity         *string   `json:"identity,omitempty"`
	ManagedResources *[]string `json:"managedResources,omitempty"`
}
