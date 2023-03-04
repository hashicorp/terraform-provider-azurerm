package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OfficePowerBIDataConnectorProperties struct {
	DataTypes OfficePowerBIConnectorDataTypes `json:"dataTypes"`
	TenantId  string                          `json:"tenantId"`
}
