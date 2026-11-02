package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentExtensionSpec struct {
	Name       string                             `json:"name"`
	Properties *DeploymentExtensionSpecProperties `json:"properties,omitempty"`
}
