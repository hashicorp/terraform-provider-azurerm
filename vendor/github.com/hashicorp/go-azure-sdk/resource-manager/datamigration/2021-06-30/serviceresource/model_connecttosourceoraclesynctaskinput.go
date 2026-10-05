package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceOracleSyncTaskInput struct {
	SourceConnectionInfo OracleConnectionInfo `json:"sourceConnectionInfo"`
}
