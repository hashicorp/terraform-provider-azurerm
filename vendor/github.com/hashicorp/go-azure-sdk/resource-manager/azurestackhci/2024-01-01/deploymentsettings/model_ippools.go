package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPPools struct {
	EndingAddress   *string `json:"endingAddress,omitempty"`
	StartingAddress *string `json:"startingAddress,omitempty"`
}
