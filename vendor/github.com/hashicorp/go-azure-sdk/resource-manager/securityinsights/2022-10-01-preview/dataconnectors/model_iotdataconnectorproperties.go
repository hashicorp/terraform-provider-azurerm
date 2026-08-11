package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IoTDataConnectorProperties struct {
	DataTypes      *AlertsDataTypeOfDataConnector `json:"dataTypes,omitempty"`
	SubscriptionId *string                        `json:"subscriptionId,omitempty"`
}
