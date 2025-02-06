package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TriggerContext = AdhocBasedTriggerContext{}

type AdhocBasedTriggerContext struct {
	TaggingCriteria AdhocBasedTaggingCriteria `json:"taggingCriteria"`

	// Fields inherited from TriggerContext

	ObjectType string `json:"objectType"`
}

func (s AdhocBasedTriggerContext) TriggerContext() BaseTriggerContextImpl {
	return BaseTriggerContextImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = AdhocBasedTriggerContext{}

func (s AdhocBasedTriggerContext) MarshalJSON() ([]byte, error) {
	type wrapper AdhocBasedTriggerContext
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AdhocBasedTriggerContext: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AdhocBasedTriggerContext: %+v", err)
	}

	decoded["objectType"] = "AdhocBasedTriggerContext"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AdhocBasedTriggerContext: %+v", err)
	}

	return encoded, nil
}
