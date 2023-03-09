package capacitypools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolProperties struct {
	CoolAccess              *bool           `json:"coolAccess,omitempty"`
	EncryptionType          *EncryptionType `json:"encryptionType,omitempty"`
	PoolId                  *string         `json:"poolId,omitempty"`
	ProvisioningState       *string         `json:"provisioningState,omitempty"`
	QosType                 *QosType        `json:"qosType,omitempty"`
	ServiceLevel            ServiceLevel    `json:"serviceLevel"`
	Size                    int64           `json:"size"`
	TotalThroughputMibps    *float64        `json:"totalThroughputMibps,omitempty"`
	UtilizedThroughputMibps *float64        `json:"utilizedThroughputMibps,omitempty"`
}
