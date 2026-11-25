package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = URLRedirectActionParameters{}

type URLRedirectActionParameters struct {
	CustomFragment      *string              `json:"customFragment,omitempty"`
	CustomHostname      *string              `json:"customHostname,omitempty"`
	CustomPath          *string              `json:"customPath,omitempty"`
	CustomQueryString   *string              `json:"customQueryString,omitempty"`
	DestinationProtocol *DestinationProtocol `json:"destinationProtocol,omitempty"`
	RedirectType        RedirectType         `json:"redirectType"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s URLRedirectActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = URLRedirectActionParameters{}

func (s URLRedirectActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper URLRedirectActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling URLRedirectActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling URLRedirectActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleUrlRedirectActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling URLRedirectActionParameters: %+v", err)
	}

	return encoded, nil
}
