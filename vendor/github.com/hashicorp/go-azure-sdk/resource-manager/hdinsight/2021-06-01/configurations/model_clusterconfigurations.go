package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterConfigurations struct {
	Configurations *map[string]map[string]string `json:"configurations,omitempty"`
}
