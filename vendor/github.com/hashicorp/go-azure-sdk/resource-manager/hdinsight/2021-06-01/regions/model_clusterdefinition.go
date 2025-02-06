package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterDefinition struct {
	Blueprint        *string            `json:"blueprint,omitempty"`
	ComponentVersion *map[string]string `json:"componentVersion,omitempty"`
	Configurations   *interface{}       `json:"configurations,omitempty"`
	Kind             *string            `json:"kind,omitempty"`
}
