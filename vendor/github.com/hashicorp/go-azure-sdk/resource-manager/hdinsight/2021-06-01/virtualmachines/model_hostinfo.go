package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostInfo struct {
	EffectiveDiskEncryptionKeyUrl *string `json:"effectiveDiskEncryptionKeyUrl,omitempty"`
	Fqdn                          *string `json:"fqdn,omitempty"`
	Name                          *string `json:"name,omitempty"`
}
