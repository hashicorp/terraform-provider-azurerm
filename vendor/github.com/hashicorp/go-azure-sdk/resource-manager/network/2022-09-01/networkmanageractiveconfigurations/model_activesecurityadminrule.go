package networkmanageractiveconfigurations

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ActiveBaseSecurityAdminRule = ActiveSecurityAdminRule{}

type ActiveSecurityAdminRule struct {
	Properties *AdminPropertiesFormat `json:"properties,omitempty"`

	// Fields inherited from ActiveBaseSecurityAdminRule
	CommitTime                    *string                            `json:"commitTime,omitempty"`
	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	Region                        *string                            `json:"region,omitempty"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

func (o *ActiveSecurityAdminRule) GetCommitTimeAsTime() (*time.Time, error) {
	if o.CommitTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CommitTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ActiveSecurityAdminRule) SetCommitTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CommitTime = &formatted
}

var _ json.Marshaler = ActiveSecurityAdminRule{}

func (s ActiveSecurityAdminRule) MarshalJSON() ([]byte, error) {
	type wrapper ActiveSecurityAdminRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ActiveSecurityAdminRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ActiveSecurityAdminRule: %+v", err)
	}
	decoded["kind"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ActiveSecurityAdminRule: %+v", err)
	}

	return encoded, nil
}
