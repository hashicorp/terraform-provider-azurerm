package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSTIDataConnectorProperties struct {
	DataTypes MSTIDataConnectorDataTypes `json:"dataTypes"`
	TenantId  string                     `json:"tenantId"`
}
