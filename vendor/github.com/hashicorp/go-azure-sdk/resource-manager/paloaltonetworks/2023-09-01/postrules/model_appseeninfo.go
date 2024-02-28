package postrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppSeenInfo struct {
	Category      string `json:"category"`
	Risk          string `json:"risk"`
	StandardPorts string `json:"standardPorts"`
	SubCategory   string `json:"subCategory"`
	Tag           string `json:"tag"`
	Technology    string `json:"technology"`
	Title         string `json:"title"`
}
