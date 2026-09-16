package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdditionalCapabilities struct {
	EnableFips1403Encryption *bool `json:"enableFips1403Encryption,omitempty"`
	HibernationEnabled       *bool `json:"hibernationEnabled,omitempty"`
	UltraSSDEnabled          *bool `json:"ultraSSDEnabled,omitempty"`
}
