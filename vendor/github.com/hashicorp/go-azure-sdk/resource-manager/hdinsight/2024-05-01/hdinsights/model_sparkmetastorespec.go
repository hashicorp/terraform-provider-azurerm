package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SparkMetastoreSpec struct {
	DbConnectionAuthenticationMode *DbConnectionAuthenticationMode `json:"dbConnectionAuthenticationMode,omitempty"`
	DbName                         string                          `json:"dbName"`
	DbPasswordSecretName           *string                         `json:"dbPasswordSecretName,omitempty"`
	DbServerHost                   string                          `json:"dbServerHost"`
	DbUserName                     *string                         `json:"dbUserName,omitempty"`
	KeyVaultId                     *string                         `json:"keyVaultId,omitempty"`
	ThriftUrl                      *string                         `json:"thriftUrl,omitempty"`
}
