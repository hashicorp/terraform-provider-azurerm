package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToTargetSqlMISyncTaskInput struct {
	AzureApp             AzureActiveDirectoryApp `json:"azureApp"`
	TargetConnectionInfo MiSqlConnectionInfo     `json:"targetConnectionInfo"`
}
