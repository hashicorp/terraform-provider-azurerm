package goldengateconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionBaseProperties = MicrosoftFabricConnectionDetails{}

type MicrosoftFabricConnectionDetails struct {
	ClientId             string                                  `json:"clientId"`
	ClientSecretSecretId *string                                 `json:"clientSecretSecretId,omitempty"`
	Endpoint             *string                                 `json:"endpoint,omitempty"`
	TechnologyType       MicrosoftFabricConnectionTechnologyType `json:"technologyType"`
	TenantId             string                                  `json:"tenantId"`

	// Fields inherited from ConnectionBaseProperties

	CompartmentId     *string                         `json:"compartmentId,omitempty"`
	ConnectionType    ConnectionType                  `json:"connectionType"`
	DisplayName       string                          `json:"displayName"`
	DoesUseSecretIds  *bool                           `json:"doesUseSecretIds,omitempty"`
	KeyId             *string                         `json:"keyId,omitempty"`
	LifecycleDetails  *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState    *ConnectionLifecycleState       `json:"lifecycleState,omitempty"`
	NetworkAnchorId   string                          `json:"networkAnchorId"`
	Ocid              *string                         `json:"ocid,omitempty"`
	ProvisioningState *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId  string                          `json:"resourceAnchorId"`
	RoutingMethod     *RoutingMethod                  `json:"routingMethod,omitempty"`
	TimeCreated       *string                         `json:"timeCreated,omitempty"`
	TimeUpdated       *string                         `json:"timeUpdated,omitempty"`
	VaultId           *string                         `json:"vaultId,omitempty"`
}

func (s MicrosoftFabricConnectionDetails) ConnectionBaseProperties() BaseConnectionBasePropertiesImpl {
	return BaseConnectionBasePropertiesImpl{
		CompartmentId:     s.CompartmentId,
		ConnectionType:    s.ConnectionType,
		DisplayName:       s.DisplayName,
		DoesUseSecretIds:  s.DoesUseSecretIds,
		KeyId:             s.KeyId,
		LifecycleDetails:  s.LifecycleDetails,
		LifecycleState:    s.LifecycleState,
		NetworkAnchorId:   s.NetworkAnchorId,
		Ocid:              s.Ocid,
		ProvisioningState: s.ProvisioningState,
		ResourceAnchorId:  s.ResourceAnchorId,
		RoutingMethod:     s.RoutingMethod,
		TimeCreated:       s.TimeCreated,
		TimeUpdated:       s.TimeUpdated,
		VaultId:           s.VaultId,
	}
}

var _ json.Marshaler = MicrosoftFabricConnectionDetails{}

func (s MicrosoftFabricConnectionDetails) MarshalJSON() ([]byte, error) {
	type wrapper MicrosoftFabricConnectionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MicrosoftFabricConnectionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MicrosoftFabricConnectionDetails: %+v", err)
	}

	decoded["connectionType"] = "MICROSOFT_FABRIC"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MicrosoftFabricConnectionDetails: %+v", err)
	}

	return encoded, nil
}
