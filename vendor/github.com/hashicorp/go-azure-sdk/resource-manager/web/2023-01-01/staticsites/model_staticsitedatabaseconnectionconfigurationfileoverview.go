package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteDatabaseConnectionConfigurationFileOverview struct {
	Contents *string `json:"contents,omitempty"`
	FileName *string `json:"fileName,omitempty"`
	Type     *string `json:"type,omitempty"`
}
