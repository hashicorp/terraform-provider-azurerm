package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcesMoveInfo struct {
	Resources           *[]string `json:"resources,omitempty"`
	TargetResourceGroup *string   `json:"targetResourceGroup,omitempty"`
}
