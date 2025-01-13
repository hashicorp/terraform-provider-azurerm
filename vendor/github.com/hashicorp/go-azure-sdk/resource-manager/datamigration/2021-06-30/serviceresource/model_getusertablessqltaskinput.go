package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUserTablesSqlTaskInput struct {
	ConnectionInfo    SqlConnectionInfo `json:"connectionInfo"`
	SelectedDatabases []string          `json:"selectedDatabases"`
}
