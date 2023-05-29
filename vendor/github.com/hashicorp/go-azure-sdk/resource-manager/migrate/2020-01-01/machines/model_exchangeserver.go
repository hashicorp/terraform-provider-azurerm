package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExchangeServer struct {
	Edition     *string `json:"edition,omitempty"`
	ProductName *string `json:"productName,omitempty"`
	Roles       *string `json:"roles,omitempty"`
	ServicePack *string `json:"servicePack,omitempty"`
	Version     *string `json:"version,omitempty"`
}
