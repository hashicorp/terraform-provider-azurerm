package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlinkHiveCatalogOption struct {
	MetastoreDbConnectionAuthenticationMode *MetastoreDbConnectionAuthenticationMode `json:"metastoreDbConnectionAuthenticationMode,omitempty"`
	MetastoreDbConnectionPasswordSecret     *string                                  `json:"metastoreDbConnectionPasswordSecret,omitempty"`
	MetastoreDbConnectionURL                string                                   `json:"metastoreDbConnectionURL"`
	MetastoreDbConnectionUserName           *string                                  `json:"metastoreDbConnectionUserName,omitempty"`
}
