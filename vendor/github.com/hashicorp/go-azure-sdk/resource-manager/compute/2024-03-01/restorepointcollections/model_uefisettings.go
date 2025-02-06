package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UefiSettings struct {
	SecureBootEnabled *bool `json:"secureBootEnabled,omitempty"`
	VTpmEnabled       *bool `json:"vTpmEnabled,omitempty"`
}
