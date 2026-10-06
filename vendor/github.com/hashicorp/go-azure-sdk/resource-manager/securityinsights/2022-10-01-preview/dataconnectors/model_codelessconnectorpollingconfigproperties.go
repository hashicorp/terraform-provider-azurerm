package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessConnectorPollingConfigProperties struct {
	Auth     CodelessConnectorPollingAuthProperties      `json:"auth"`
	IsActive *bool                                       `json:"isActive,omitempty"`
	Paging   *CodelessConnectorPollingPagingProperties   `json:"paging,omitempty"`
	Request  CodelessConnectorPollingRequestProperties   `json:"request"`
	Response *CodelessConnectorPollingResponseProperties `json:"response,omitempty"`
}
