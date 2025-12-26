package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTableDatasetTypeProperties struct {
	Index int64        `json:"index"`
	Path  *interface{} `json:"path,omitempty"`
}
