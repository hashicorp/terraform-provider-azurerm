package defenderforstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnUploadFilters struct {
	ExcludeBlobsLargerThan *interface{} `json:"excludeBlobsLargerThan,omitempty"`
	ExcludeBlobsWithPrefix *[]string    `json:"excludeBlobsWithPrefix,omitempty"`
	ExcludeBlobsWithSuffix *[]string    `json:"excludeBlobsWithSuffix,omitempty"`
}
