package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterResourceProperties struct {
	Fqdn                             *string                          `json:"fqdn,omitempty"`
	InfraResourceGroup               *string                          `json:"infraResourceGroup,omitempty"`
	MaintenanceScheduleConfiguration MaintenanceScheduleConfiguration `json:"maintenanceScheduleConfiguration"`
	ManagedEnvironmentId             *string                          `json:"managedEnvironmentId,omitempty"`
	MarketplaceResource              *MarketplaceResource             `json:"marketplaceResource,omitempty"`
	NetworkProfile                   *NetworkProfile                  `json:"networkProfile,omitempty"`
	PowerState                       *PowerState                      `json:"powerState,omitempty"`
	ProvisioningState                *ProvisioningState               `json:"provisioningState,omitempty"`
	ServiceId                        *string                          `json:"serviceId,omitempty"`
	Version                          *int64                           `json:"version,omitempty"`
	VnetAddons                       *ServiceVNetAddons               `json:"vnetAddons,omitempty"`
	ZoneRedundant                    *bool                            `json:"zoneRedundant,omitempty"`
}

var _ json.Unmarshaler = &ClusterResourceProperties{}

func (s *ClusterResourceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Fqdn                 *string              `json:"fqdn,omitempty"`
		InfraResourceGroup   *string              `json:"infraResourceGroup,omitempty"`
		ManagedEnvironmentId *string              `json:"managedEnvironmentId,omitempty"`
		MarketplaceResource  *MarketplaceResource `json:"marketplaceResource,omitempty"`
		NetworkProfile       *NetworkProfile      `json:"networkProfile,omitempty"`
		PowerState           *PowerState          `json:"powerState,omitempty"`
		ProvisioningState    *ProvisioningState   `json:"provisioningState,omitempty"`
		ServiceId            *string              `json:"serviceId,omitempty"`
		Version              *int64               `json:"version,omitempty"`
		VnetAddons           *ServiceVNetAddons   `json:"vnetAddons,omitempty"`
		ZoneRedundant        *bool                `json:"zoneRedundant,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Fqdn = decoded.Fqdn
	s.InfraResourceGroup = decoded.InfraResourceGroup
	s.ManagedEnvironmentId = decoded.ManagedEnvironmentId
	s.MarketplaceResource = decoded.MarketplaceResource
	s.NetworkProfile = decoded.NetworkProfile
	s.PowerState = decoded.PowerState
	s.ProvisioningState = decoded.ProvisioningState
	s.ServiceId = decoded.ServiceId
	s.Version = decoded.Version
	s.VnetAddons = decoded.VnetAddons
	s.ZoneRedundant = decoded.ZoneRedundant

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ClusterResourceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["maintenanceScheduleConfiguration"]; ok {
		impl, err := UnmarshalMaintenanceScheduleConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'MaintenanceScheduleConfiguration' for 'ClusterResourceProperties': %+v", err)
		}
		s.MaintenanceScheduleConfiguration = impl
	}

	return nil
}
