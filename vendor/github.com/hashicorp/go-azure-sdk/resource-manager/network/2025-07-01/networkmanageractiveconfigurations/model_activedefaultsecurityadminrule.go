package networkmanageractiveconfigurations

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ActiveBaseSecurityAdminRule = ActiveDefaultSecurityAdminRule{}

type ActiveDefaultSecurityAdminRule struct {
	Properties *DefaultAdminPropertiesFormat `json:"properties,omitempty"`

	// Fields inherited from ActiveBaseSecurityAdminRule

	CommitTime                    *string                            `json:"commitTime,omitempty"`
	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	Kind                          EffectiveAdminRuleKind             `json:"kind"`
	Region                        *string                            `json:"region,omitempty"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

func (s ActiveDefaultSecurityAdminRule) ActiveBaseSecurityAdminRule() BaseActiveBaseSecurityAdminRuleImpl {
	return BaseActiveBaseSecurityAdminRuleImpl{
		CommitTime:                    s.CommitTime,
		ConfigurationDescription:      s.ConfigurationDescription,
		Id:                            s.Id,
		Kind:                          s.Kind,
		Region:                        s.Region,
		RuleCollectionAppliesToGroups: s.RuleCollectionAppliesToGroups,
		RuleCollectionDescription:     s.RuleCollectionDescription,
		RuleGroups:                    s.RuleGroups,
	}
}

func (o *ActiveDefaultSecurityAdminRule) GetCommitTimeAsTime() (*time.Time, error) {
	if o.CommitTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CommitTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ActiveDefaultSecurityAdminRule) SetCommitTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CommitTime = &formatted
}

var _ json.Marshaler = ActiveDefaultSecurityAdminRule{}

func (s ActiveDefaultSecurityAdminRule) MarshalJSON() ([]byte, error) {
	type wrapper ActiveDefaultSecurityAdminRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ActiveDefaultSecurityAdminRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ActiveDefaultSecurityAdminRule: %+v", err)
	}

	decoded["kind"] = "Default"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ActiveDefaultSecurityAdminRule: %+v", err)
	}

	return encoded, nil
}
