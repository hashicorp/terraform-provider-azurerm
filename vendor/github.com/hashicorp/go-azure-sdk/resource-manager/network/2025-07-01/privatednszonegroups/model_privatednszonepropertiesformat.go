package privatednszonegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateDnsZonePropertiesFormat struct {
	PrivateDnsZoneId *string      `json:"privateDnsZoneId,omitempty"`
	RecordSets       *[]RecordSet `json:"recordSets,omitempty"`
}
