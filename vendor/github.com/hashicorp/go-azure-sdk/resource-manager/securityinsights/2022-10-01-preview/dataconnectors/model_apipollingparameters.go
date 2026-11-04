package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiPollingParameters struct {
	ConnectorUiConfig *CodelessUiConnectorConfigProperties      `json:"connectorUiConfig,omitempty"`
	PollingConfig     *CodelessConnectorPollingConfigProperties `json:"pollingConfig,omitempty"`
}
