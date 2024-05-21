package replicationrecoveryservicesproviders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryServicesProviderProperties struct {
	AllowedScenarios                       *[]string                `json:"allowedScenarios,omitempty"`
	AuthenticationIdentityDetails          *IdentityProviderDetails `json:"authenticationIdentityDetails,omitempty"`
	BiosId                                 *string                  `json:"biosId,omitempty"`
	ConnectionStatus                       *string                  `json:"connectionStatus,omitempty"`
	DataPlaneAuthenticationIdentityDetails *IdentityProviderDetails `json:"dataPlaneAuthenticationIdentityDetails,omitempty"`
	DraIdentifier                          *string                  `json:"draIdentifier,omitempty"`
	FabricFriendlyName                     *string                  `json:"fabricFriendlyName,omitempty"`
	FabricType                             *string                  `json:"fabricType,omitempty"`
	FriendlyName                           *string                  `json:"friendlyName,omitempty"`
	HealthErrorDetails                     *[]HealthError           `json:"healthErrorDetails,omitempty"`
	LastHeartBeat                          *string                  `json:"lastHeartBeat,omitempty"`
	MachineId                              *string                  `json:"machineId,omitempty"`
	MachineName                            *string                  `json:"machineName,omitempty"`
	ProtectedItemCount                     *int64                   `json:"protectedItemCount,omitempty"`
	ProviderVersion                        *string                  `json:"providerVersion,omitempty"`
	ProviderVersionDetails                 *VersionDetails          `json:"providerVersionDetails,omitempty"`
	ProviderVersionExpiryDate              *string                  `json:"providerVersionExpiryDate,omitempty"`
	ProviderVersionState                   *string                  `json:"providerVersionState,omitempty"`
	ResourceAccessIdentityDetails          *IdentityProviderDetails `json:"resourceAccessIdentityDetails,omitempty"`
	ServerVersion                          *string                  `json:"serverVersion,omitempty"`
}

func (o *RecoveryServicesProviderProperties) GetLastHeartBeatAsTime() (*time.Time, error) {
	if o.LastHeartBeat == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartBeat, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryServicesProviderProperties) SetLastHeartBeatAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartBeat = &formatted
}

func (o *RecoveryServicesProviderProperties) GetProviderVersionExpiryDateAsTime() (*time.Time, error) {
	if o.ProviderVersionExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProviderVersionExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RecoveryServicesProviderProperties) SetProviderVersionExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProviderVersionExpiryDate = &formatted
}
