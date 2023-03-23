package dashboard

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardParts struct {
	Metadata *interface{}           `json:"metadata,omitempty"`
	Position DashboardPartsPosition `json:"position"`
}
