package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CorsRules struct {
	CorsRules *[]CorsRule `json:"corsRules,omitempty"`
}
