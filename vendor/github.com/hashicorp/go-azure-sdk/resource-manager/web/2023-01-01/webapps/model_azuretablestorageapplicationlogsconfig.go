package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureTableStorageApplicationLogsConfig struct {
	Level  *LogLevel `json:"level,omitempty"`
	SasUrl string    `json:"sasUrl"`
}
