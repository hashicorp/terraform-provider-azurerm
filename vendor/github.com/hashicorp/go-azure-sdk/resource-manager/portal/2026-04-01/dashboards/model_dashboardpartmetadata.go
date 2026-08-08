package dashboards

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardPartMetadata interface {
	DashboardPartMetadata() BaseDashboardPartMetadataImpl
}

var _ DashboardPartMetadata = BaseDashboardPartMetadataImpl{}

type BaseDashboardPartMetadataImpl struct {
	Type DashboardPartMetadataType `json:"type"`
}

func (s BaseDashboardPartMetadataImpl) DashboardPartMetadata() BaseDashboardPartMetadataImpl {
	return s
}

var _ DashboardPartMetadata = RawDashboardPartMetadataImpl{}

// RawDashboardPartMetadataImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDashboardPartMetadataImpl struct {
	dashboardPartMetadata BaseDashboardPartMetadataImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawDashboardPartMetadataImpl) DashboardPartMetadata() BaseDashboardPartMetadataImpl {
	return s.dashboardPartMetadata
}

func UnmarshalDashboardPartMetadataImplementation(input []byte) (DashboardPartMetadata, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DashboardPartMetadata into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Extension/HubsExtension/PartType/MarkdownPart") {
		var out MarkdownPartMetadata
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MarkdownPartMetadata: %+v", err)
		}
		return out, nil
	}

	var parent BaseDashboardPartMetadataImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDashboardPartMetadataImpl: %+v", err)
	}

	return RawDashboardPartMetadataImpl{
		dashboardPartMetadata: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
