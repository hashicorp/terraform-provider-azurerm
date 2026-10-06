package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = URLRedirectAction{}

type URLRedirectAction struct {
	Parameters URLRedirectActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s URLRedirectAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = URLRedirectAction{}

func (s URLRedirectAction) MarshalJSON() ([]byte, error) {
	type wrapper URLRedirectAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLRedirectAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLRedirectAction: %+v", err)
	}

	decoded["name"] = "UrlRedirect"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLRedirectAction: %+v", err)
	}

	return encoded, nil
}
