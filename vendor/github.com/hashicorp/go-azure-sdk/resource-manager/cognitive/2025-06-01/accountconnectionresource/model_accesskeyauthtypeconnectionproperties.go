package accountconnectionresource

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ConnectionPropertiesV2 = AccessKeyAuthTypeConnectionProperties{}

type AccessKeyAuthTypeConnectionProperties struct {
	Credentials *ConnectionAccessKey `json:"credentials,omitempty"`

	// Fields inherited from ConnectionPropertiesV2

	AuthType                    ConnectionAuthType    `json:"authType"`
	Category                    *ConnectionCategory   `json:"category,omitempty"`
	CreatedByWorkspaceArmId     *string               `json:"createdByWorkspaceArmId,omitempty"`
	Error                       *string               `json:"error,omitempty"`
	ExpiryTime                  *string               `json:"expiryTime,omitempty"`
	Group                       *ConnectionGroup      `json:"group,omitempty"`
	IsSharedToAll               *bool                 `json:"isSharedToAll,omitempty"`
	Metadata                    *map[string]string    `json:"metadata,omitempty"`
	PeRequirement               *ManagedPERequirement `json:"peRequirement,omitempty"`
	PeStatus                    *ManagedPEStatus      `json:"peStatus,omitempty"`
	SharedUserList              *[]string             `json:"sharedUserList,omitempty"`
	Target                      *string               `json:"target,omitempty"`
	UseWorkspaceManagedIdentity *bool                 `json:"useWorkspaceManagedIdentity,omitempty"`
}

func (s AccessKeyAuthTypeConnectionProperties) ConnectionPropertiesV2() BaseConnectionPropertiesV2Impl {
	return BaseConnectionPropertiesV2Impl{
		AuthType:                    s.AuthType,
		Category:                    s.Category,
		CreatedByWorkspaceArmId:     s.CreatedByWorkspaceArmId,
		Error:                       s.Error,
		ExpiryTime:                  s.ExpiryTime,
		Group:                       s.Group,
		IsSharedToAll:               s.IsSharedToAll,
		Metadata:                    s.Metadata,
		PeRequirement:               s.PeRequirement,
		PeStatus:                    s.PeStatus,
		SharedUserList:              s.SharedUserList,
		Target:                      s.Target,
		UseWorkspaceManagedIdentity: s.UseWorkspaceManagedIdentity,
	}
}

func (o *AccessKeyAuthTypeConnectionProperties) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AccessKeyAuthTypeConnectionProperties) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

var _ json.Marshaler = AccessKeyAuthTypeConnectionProperties{}

func (s AccessKeyAuthTypeConnectionProperties) MarshalJSON() ([]byte, error) {
	type wrapper AccessKeyAuthTypeConnectionProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AccessKeyAuthTypeConnectionProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AccessKeyAuthTypeConnectionProperties: %+v", err)
	}

	decoded["authType"] = "AccessKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AccessKeyAuthTypeConnectionProperties: %+v", err)
	}

	return encoded, nil
}
