package storageinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightProperties struct {
	Containers     *[]string             `json:"containers,omitempty"`
	Status         *StorageInsightStatus `json:"status,omitempty"`
	StorageAccount StorageAccount        `json:"storageAccount"`
	Tables         *[]string             `json:"tables,omitempty"`
}
