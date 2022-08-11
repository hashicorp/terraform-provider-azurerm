package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleRule struct {
	AzureQueue *QueueScaleRule  `json:"azureQueue,omitempty"`
	Custom     *CustomScaleRule `json:"custom,omitempty"`
	Http       *HttpScaleRule   `json:"http,omitempty"`
	Name       *string          `json:"name,omitempty"`
}
