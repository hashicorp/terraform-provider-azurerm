package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OsProfile struct {
	LogonType                 *LogonType                 `json:"logonType,omitempty"`
	SecretsManagementSettings *SecretsManagementSettings `json:"secretsManagementSettings,omitempty"`
}
