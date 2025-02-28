package networkstatus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkStatusContractByLocation struct {
	Location      *string                `json:"location,omitempty"`
	NetworkStatus *NetworkStatusContract `json:"networkStatus,omitempty"`
}
