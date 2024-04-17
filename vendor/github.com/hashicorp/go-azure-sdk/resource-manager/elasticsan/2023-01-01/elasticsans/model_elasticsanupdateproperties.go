package elasticsans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticSanUpdateProperties struct {
	BaseSizeTiB             *int64               `json:"baseSizeTiB,omitempty"`
	ExtendedCapacitySizeTiB *int64               `json:"extendedCapacitySizeTiB,omitempty"`
	PublicNetworkAccess     *PublicNetworkAccess `json:"publicNetworkAccess,omitempty"`
}
