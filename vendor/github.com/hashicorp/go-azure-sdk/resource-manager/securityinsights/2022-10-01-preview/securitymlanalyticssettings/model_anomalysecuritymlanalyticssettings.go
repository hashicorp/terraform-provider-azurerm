package securitymlanalyticssettings

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecurityMLAnalyticsSetting = AnomalySecurityMLAnalyticsSettings{}

type AnomalySecurityMLAnalyticsSettings struct {
	Properties *AnomalySecurityMLAnalyticsSettingsProperties `json:"properties,omitempty"`

	// Fields inherited from SecurityMLAnalyticsSetting

	Etag       *string                         `json:"etag,omitempty"`
	Id         *string                         `json:"id,omitempty"`
	Kind       SecurityMLAnalyticsSettingsKind `json:"kind"`
	Name       *string                         `json:"name,omitempty"`
	SystemData *systemdata.SystemData          `json:"systemData,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}

func (s AnomalySecurityMLAnalyticsSettings) SecurityMLAnalyticsSetting() BaseSecurityMLAnalyticsSettingImpl {
	return BaseSecurityMLAnalyticsSettingImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = AnomalySecurityMLAnalyticsSettings{}

func (s AnomalySecurityMLAnalyticsSettings) MarshalJSON() ([]byte, error) {
	type wrapper AnomalySecurityMLAnalyticsSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AnomalySecurityMLAnalyticsSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AnomalySecurityMLAnalyticsSettings: %+v", err)
	}

	decoded["kind"] = "Anomaly"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AnomalySecurityMLAnalyticsSettings: %+v", err)
	}

	return encoded, nil
}
