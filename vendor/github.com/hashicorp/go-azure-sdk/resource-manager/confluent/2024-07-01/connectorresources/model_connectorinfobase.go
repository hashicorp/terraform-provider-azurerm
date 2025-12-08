package connectorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectorInfoBase struct {
	ConnectorClass *ConnectorClass  `json:"connectorClass,omitempty"`
	ConnectorId    *string          `json:"connectorId,omitempty"`
	ConnectorName  *string          `json:"connectorName,omitempty"`
	ConnectorState *ConnectorStatus `json:"connectorState,omitempty"`
	ConnectorType  *ConnectorType   `json:"connectorType,omitempty"`
}
