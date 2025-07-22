package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetStatus struct {
	Datasets *[]AssetStatusDataset `json:"datasets,omitempty"`
	Errors   *[]AssetStatusError   `json:"errors,omitempty"`
	Events   *[]AssetStatusEvent   `json:"events,omitempty"`
	Version  *int64                `json:"version,omitempty"`
}
