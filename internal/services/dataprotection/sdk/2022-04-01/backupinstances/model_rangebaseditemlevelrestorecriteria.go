package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = RangeBasedItemLevelRestoreCriteria{}

type RangeBasedItemLevelRestoreCriteria struct {
	MaxMatchingValue *string `json:"maxMatchingValue,omitempty"`
	MinMatchingValue *string `json:"minMatchingValue,omitempty"`

	// Fields inherited from ItemLevelRestoreCriteria
}

var _ json.Marshaler = RangeBasedItemLevelRestoreCriteria{}

func (s RangeBasedItemLevelRestoreCriteria) MarshalJSON() ([]byte, error) {
	type wrapper RangeBasedItemLevelRestoreCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RangeBasedItemLevelRestoreCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RangeBasedItemLevelRestoreCriteria: %+v", err)
	}
	decoded["objectType"] = "RangeBasedItemLevelRestoreCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RangeBasedItemLevelRestoreCriteria: %+v", err)
	}

	return encoded, nil
}
