package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SapOpenHubTableDatasetTypeProperties struct {
	BaseRequestId          *int64 `json:"baseRequestId,omitempty"`
	ExcludeLastRequest     *bool  `json:"excludeLastRequest,omitempty"`
	OpenHubDestinationName string `json:"openHubDestinationName"`
}
