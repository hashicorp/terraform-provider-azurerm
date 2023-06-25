package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = VMwareDetails{}

type VMwareDetails struct {
	AgentCount                         *string                                           `json:"agentCount,omitempty"`
	AgentExpiryDate                    *string                                           `json:"agentExpiryDate,omitempty"`
	AgentVersion                       *string                                           `json:"agentVersion,omitempty"`
	AgentVersionDetails                *VersionDetails                                   `json:"agentVersionDetails,omitempty"`
	AvailableMemoryInBytes             *int64                                            `json:"availableMemoryInBytes,omitempty"`
	AvailableSpaceInBytes              *int64                                            `json:"availableSpaceInBytes,omitempty"`
	CpuLoad                            *string                                           `json:"cpuLoad,omitempty"`
	CpuLoadStatus                      *string                                           `json:"cpuLoadStatus,omitempty"`
	CsServiceStatus                    *string                                           `json:"csServiceStatus,omitempty"`
	DatabaseServerLoad                 *string                                           `json:"databaseServerLoad,omitempty"`
	DatabaseServerLoadStatus           *string                                           `json:"databaseServerLoadStatus,omitempty"`
	HostName                           *string                                           `json:"hostName,omitempty"`
	IPAddress                          *string                                           `json:"ipAddress,omitempty"`
	LastHeartbeat                      *string                                           `json:"lastHeartbeat,omitempty"`
	MasterTargetServers                *[]MasterTargetServer                             `json:"masterTargetServers,omitempty"`
	MemoryUsageStatus                  *string                                           `json:"memoryUsageStatus,omitempty"`
	ProcessServerCount                 *string                                           `json:"processServerCount,omitempty"`
	ProcessServers                     *[]ProcessServer                                  `json:"processServers,omitempty"`
	ProtectedServers                   *string                                           `json:"protectedServers,omitempty"`
	PsTemplateVersion                  *string                                           `json:"psTemplateVersion,omitempty"`
	ReplicationPairCount               *string                                           `json:"replicationPairCount,omitempty"`
	RunAsAccounts                      *[]RunAsAccount                                   `json:"runAsAccounts,omitempty"`
	SpaceUsageStatus                   *string                                           `json:"spaceUsageStatus,omitempty"`
	SslCertExpiryDate                  *string                                           `json:"sslCertExpiryDate,omitempty"`
	SslCertExpiryRemainingDays         *int64                                            `json:"sslCertExpiryRemainingDays,omitempty"`
	SwitchProviderBlockingErrorDetails *[]InMageFabricSwitchProviderBlockingErrorDetails `json:"switchProviderBlockingErrorDetails,omitempty"`
	SystemLoad                         *string                                           `json:"systemLoad,omitempty"`
	SystemLoadStatus                   *string                                           `json:"systemLoadStatus,omitempty"`
	TotalMemoryInBytes                 *int64                                            `json:"totalMemoryInBytes,omitempty"`
	TotalSpaceInBytes                  *int64                                            `json:"totalSpaceInBytes,omitempty"`
	VersionStatus                      *string                                           `json:"versionStatus,omitempty"`
	WebLoad                            *string                                           `json:"webLoad,omitempty"`
	WebLoadStatus                      *string                                           `json:"webLoadStatus,omitempty"`

	// Fields inherited from FabricSpecificDetails
}

var _ json.Marshaler = VMwareDetails{}

func (s VMwareDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMwareDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareDetails: %+v", err)
	}
	decoded["instanceType"] = "VMware"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareDetails: %+v", err)
	}

	return encoded, nil
}
