package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OfficeDataConnectorDataTypes struct {
	Exchange   *DataConnectorDataTypeCommon `json:"exchange,omitempty"`
	SharePoint *DataConnectorDataTypeCommon `json:"sharePoint,omitempty"`
	Teams      *DataConnectorDataTypeCommon `json:"teams,omitempty"`
}
