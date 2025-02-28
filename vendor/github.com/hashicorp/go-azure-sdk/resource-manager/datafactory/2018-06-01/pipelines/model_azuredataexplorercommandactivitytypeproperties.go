package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDataExplorerCommandActivityTypeProperties struct {
	Command        string  `json:"command"`
	CommandTimeout *string `json:"commandTimeout,omitempty"`
}
