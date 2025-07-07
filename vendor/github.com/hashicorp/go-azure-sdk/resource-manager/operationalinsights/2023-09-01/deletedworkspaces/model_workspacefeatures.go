package deletedworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceFeatures struct {
	ClusterResourceId                           *string `json:"clusterResourceId,omitempty"`
	DisableLocalAuth                            *bool   `json:"disableLocalAuth,omitempty"`
	EnableDataExport                            *bool   `json:"enableDataExport,omitempty"`
	EnableLogAccessUsingOnlyResourcePermissions *bool   `json:"enableLogAccessUsingOnlyResourcePermissions,omitempty"`
	ImmediatePurgeDataOn30Days                  *bool   `json:"immediatePurgeDataOn30Days,omitempty"`
	UnifiedSentinelBillingOnly                  *bool   `json:"unifiedSentinelBillingOnly,omitempty"`
}
