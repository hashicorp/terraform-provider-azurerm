package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterServiceConfig struct {
	Component string              `json:"component"`
	Files     []ClusterConfigFile `json:"files"`
}
