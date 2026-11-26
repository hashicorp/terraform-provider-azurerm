package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessConnectorPollingPagingProperties struct {
	NextPageParaName                       *string `json:"nextPageParaName,omitempty"`
	NextPageTokenJsonPath                  *string `json:"nextPageTokenJsonPath,omitempty"`
	PageCountAttributePath                 *string `json:"pageCountAttributePath,omitempty"`
	PageSize                               *int64  `json:"pageSize,omitempty"`
	PageSizeParaName                       *string `json:"pageSizeParaName,omitempty"`
	PageTimeStampAttributePath             *string `json:"pageTimeStampAttributePath,omitempty"`
	PageTotalCountAttributePath            *string `json:"pageTotalCountAttributePath,omitempty"`
	PagingType                             string  `json:"pagingType"`
	SearchTheLatestTimeStampFromEventsList *string `json:"searchTheLatestTimeStampFromEventsList,omitempty"`
}
