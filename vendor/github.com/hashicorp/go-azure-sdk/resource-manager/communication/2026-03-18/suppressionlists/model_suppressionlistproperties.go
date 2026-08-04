package suppressionlists

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionListProperties struct {
	CreatedTimeStamp     *string `json:"createdTimeStamp,omitempty"`
	DataLocation         *string `json:"dataLocation,omitempty"`
	LastUpdatedTimeStamp *string `json:"lastUpdatedTimeStamp,omitempty"`
	ListName             *string `json:"listName,omitempty"`
}
