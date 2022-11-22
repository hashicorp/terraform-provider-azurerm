package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OSDiskDetails struct {
	OsType  *string `json:"osType,omitempty"`
	OsVhdId *string `json:"osVhdId,omitempty"`
	VhdName *string `json:"vhdName,omitempty"`
}
