package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppsAndRoles struct {
	Applications      *[]Application      `json:"applications,omitempty"`
	BizTalkServers    *[]BizTalkServer    `json:"bizTalkServers,omitempty"`
	ExchangeServers   *[]ExchangeServer   `json:"exchangeServers,omitempty"`
	Features          *[]Feature          `json:"features,omitempty"`
	OtherDatabases    *[]OtherDatabase    `json:"otherDatabases,omitempty"`
	SharePointServers *[]SharePointServer `json:"sharePointServers,omitempty"`
	SqlServers        *[]SQLServer        `json:"sqlServers,omitempty"`
	SystemCenters     *[]SystemCenter     `json:"systemCenters,omitempty"`
	WebApplications   *[]WebApplication   `json:"webApplications,omitempty"`
}
