package eventhubsclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterSku struct {
	Capacity *int64         `json:"capacity,omitempty"`
	Name     ClusterSkuName `json:"name"`
}
