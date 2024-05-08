package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgrade struct {
	Properties ClusterUpgradeProperties `json:"properties"`
}

var _ json.Unmarshaler = &ClusterUpgrade{}

func (s *ClusterUpgrade) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ClusterUpgrade into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalClusterUpgradePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'ClusterUpgrade': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
