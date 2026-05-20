package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBlobStorageApplicationLogsConfig struct {
	Level           *LogLevel `json:"level,omitempty"`
	RetentionInDays *int64    `json:"retentionInDays,omitempty"`
	SasURL          *string   `json:"sasUrl,omitempty"`
}
