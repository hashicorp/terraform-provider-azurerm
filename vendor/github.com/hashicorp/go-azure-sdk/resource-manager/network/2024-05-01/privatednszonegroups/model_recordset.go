package privatednszonegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSet struct {
	Fqdn              *string            `json:"fqdn,omitempty"`
	IPAddresses       *[]string          `json:"ipAddresses,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RecordSetName     *string            `json:"recordSetName,omitempty"`
	RecordType        *string            `json:"recordType,omitempty"`
	Ttl               *int64             `json:"ttl,omitempty"`
}
