package rolemanagementpolicies

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleManagementPolicyProperties struct {
	Description           *string                     `json:"description,omitempty"`
	DisplayName           *string                     `json:"displayName,omitempty"`
	EffectiveRules        *[]RoleManagementPolicyRule `json:"effectiveRules,omitempty"`
	IsOrganizationDefault *bool                       `json:"isOrganizationDefault,omitempty"`
	LastModifiedBy        *Principal                  `json:"lastModifiedBy,omitempty"`
	LastModifiedDateTime  *string                     `json:"lastModifiedDateTime,omitempty"`
	PolicyProperties      *PolicyProperties           `json:"policyProperties,omitempty"`
	Rules                 *[]RoleManagementPolicyRule `json:"rules,omitempty"`
	Scope                 *string                     `json:"scope,omitempty"`
}

func (o *RoleManagementPolicyProperties) GetLastModifiedDateTimeAsTime() (*time.Time, error) {
	if o.LastModifiedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleManagementPolicyProperties) SetLastModifiedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDateTime = &formatted
}

var _ json.Unmarshaler = &RoleManagementPolicyProperties{}

func (s *RoleManagementPolicyProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Description           *string           `json:"description,omitempty"`
		DisplayName           *string           `json:"displayName,omitempty"`
		IsOrganizationDefault *bool             `json:"isOrganizationDefault,omitempty"`
		LastModifiedBy        *Principal        `json:"lastModifiedBy,omitempty"`
		LastModifiedDateTime  *string           `json:"lastModifiedDateTime,omitempty"`
		PolicyProperties      *PolicyProperties `json:"policyProperties,omitempty"`
		Scope                 *string           `json:"scope,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.IsOrganizationDefault = decoded.IsOrganizationDefault
	s.LastModifiedBy = decoded.LastModifiedBy
	s.LastModifiedDateTime = decoded.LastModifiedDateTime
	s.PolicyProperties = decoded.PolicyProperties
	s.Scope = decoded.Scope

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RoleManagementPolicyProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["effectiveRules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling EffectiveRules into list []json.RawMessage: %+v", err)
		}

		output := make([]RoleManagementPolicyRule, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalRoleManagementPolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'EffectiveRules' for 'RoleManagementPolicyProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.EffectiveRules = &output
	}

	if v, ok := temp["rules"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Rules into list []json.RawMessage: %+v", err)
		}

		output := make([]RoleManagementPolicyRule, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalRoleManagementPolicyRuleImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Rules' for 'RoleManagementPolicyProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Rules = &output
	}

	return nil
}
