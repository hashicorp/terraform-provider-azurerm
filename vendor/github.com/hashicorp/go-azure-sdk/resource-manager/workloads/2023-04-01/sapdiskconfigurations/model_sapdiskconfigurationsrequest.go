package sapdiskconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPDiskConfigurationsRequest struct {
	AppLocation    string             `json:"appLocation"`
	DatabaseType   SAPDatabaseType    `json:"databaseType"`
	DbVMSku        string             `json:"dbVmSku"`
	DeploymentType SAPDeploymentType  `json:"deploymentType"`
	Environment    SAPEnvironmentType `json:"environment"`
	SapProduct     SAPProductType     `json:"sapProduct"`
}
