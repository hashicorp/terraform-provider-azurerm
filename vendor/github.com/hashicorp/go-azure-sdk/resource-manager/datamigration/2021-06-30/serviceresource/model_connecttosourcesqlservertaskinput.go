package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceSqlServerTaskInput struct {
	CheckPermissionsGroup     *ServerLevelPermissionsGroup `json:"checkPermissionsGroup,omitempty"`
	CollectAgentJobs          *bool                        `json:"collectAgentJobs,omitempty"`
	CollectDatabases          *bool                        `json:"collectDatabases,omitempty"`
	CollectLogins             *bool                        `json:"collectLogins,omitempty"`
	CollectTdeCertificateInfo *bool                        `json:"collectTdeCertificateInfo,omitempty"`
	SourceConnectionInfo      SqlConnectionInfo            `json:"sourceConnectionInfo"`
	ValidateSsisCatalogOnly   *bool                        `json:"validateSsisCatalogOnly,omitempty"`
}
