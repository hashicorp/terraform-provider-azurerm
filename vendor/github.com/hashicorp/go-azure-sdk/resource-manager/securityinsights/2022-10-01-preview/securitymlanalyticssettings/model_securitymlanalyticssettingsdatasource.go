package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingsDataSource struct {
	ConnectorId *string   `json:"connectorId,omitempty"`
	DataTypes   *[]string `json:"dataTypes,omitempty"`
}
