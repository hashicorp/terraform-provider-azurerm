package bastionshareablelink

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionShareableLinkListRequest struct {
	VMs *[]BastionShareableLink `json:"vms,omitempty"`
}
