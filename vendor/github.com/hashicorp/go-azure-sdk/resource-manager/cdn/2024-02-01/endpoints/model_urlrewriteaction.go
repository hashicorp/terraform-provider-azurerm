package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = URLRewriteAction{}

type URLRewriteAction struct {
	Parameters URLRewriteActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s URLRewriteAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = URLRewriteAction{}

func (s URLRewriteAction) MarshalJSON() ([]byte, error) {
	type wrapper URLRewriteAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLRewriteAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLRewriteAction: %+v", err)
	}

	decoded["name"] = "UrlRewrite"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLRewriteAction: %+v", err)
	}

	return encoded, nil
}
