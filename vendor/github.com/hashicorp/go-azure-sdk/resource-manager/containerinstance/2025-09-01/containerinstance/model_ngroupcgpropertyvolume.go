package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NGroupCGPropertyVolume struct {
	AzureFile *AzureFileVolume `json:"azureFile,omitempty"`
	Name      string           `json:"name"`
}
