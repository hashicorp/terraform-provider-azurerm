package factories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryProperties struct {
	CreateTime           *string                                  `json:"createTime,omitempty"`
	Encryption           *EncryptionConfiguration                 `json:"encryption,omitempty"`
	GlobalParameters     *map[string]GlobalParameterSpecification `json:"globalParameters,omitempty"`
	ProvisioningState    *string                                  `json:"provisioningState,omitempty"`
	PublicNetworkAccess  *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
	PurviewConfiguration *PurviewConfiguration                    `json:"purviewConfiguration,omitempty"`
	RepoConfiguration    FactoryRepoConfiguration                 `json:"repoConfiguration"`
	Version              *string                                  `json:"version,omitempty"`
}

func (o *FactoryProperties) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *FactoryProperties) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

var _ json.Unmarshaler = &FactoryProperties{}

func (s *FactoryProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CreateTime           *string                                  `json:"createTime,omitempty"`
		Encryption           *EncryptionConfiguration                 `json:"encryption,omitempty"`
		GlobalParameters     *map[string]GlobalParameterSpecification `json:"globalParameters,omitempty"`
		ProvisioningState    *string                                  `json:"provisioningState,omitempty"`
		PublicNetworkAccess  *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
		PurviewConfiguration *PurviewConfiguration                    `json:"purviewConfiguration,omitempty"`
		Version              *string                                  `json:"version,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CreateTime = decoded.CreateTime
	s.Encryption = decoded.Encryption
	s.GlobalParameters = decoded.GlobalParameters
	s.ProvisioningState = decoded.ProvisioningState
	s.PublicNetworkAccess = decoded.PublicNetworkAccess
	s.PurviewConfiguration = decoded.PurviewConfiguration
	s.Version = decoded.Version

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FactoryProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["repoConfiguration"]; ok {
		impl, err := UnmarshalFactoryRepoConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RepoConfiguration' for 'FactoryProperties': %+v", err)
		}
		s.RepoConfiguration = impl
	}

	return nil
}
