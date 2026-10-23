package goldengateconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionBaseProperties = OracleConnectionDetails{}

type OracleConnectionDetails struct {
	AuthenticationMode *string                        `json:"authenticationMode,omitempty"`
	ConnectionString   *string                        `json:"connectionString,omitempty"`
	DatabaseId         *string                        `json:"databaseId,omitempty"`
	PasswordSecretId   *string                        `json:"passwordSecretId,omitempty"`
	PrivateIP          *string                        `json:"privateIp,omitempty"`
	SessionMode        *SessionMode                   `json:"sessionMode,omitempty"`
	TechnologyType     OracleConnectionTechnologyType `json:"technologyType"`
	Username           string                         `json:"username"`
	WalletSecretId     *string                        `json:"walletSecretId,omitempty"`

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

func (s OracleConnectionDetails) ConnectionBaseProperties() BaseConnectionBasePropertiesImpl {
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

var _ json.Marshaler = OracleConnectionDetails{}

func (s OracleConnectionDetails) MarshalJSON() ([]byte, error) {
	type wrapper OracleConnectionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OracleConnectionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OracleConnectionDetails: %+v", err)
	}

	decoded["connectionType"] = "ORACLE"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OracleConnectionDetails: %+v", err)
	}

	return encoded, nil
}
