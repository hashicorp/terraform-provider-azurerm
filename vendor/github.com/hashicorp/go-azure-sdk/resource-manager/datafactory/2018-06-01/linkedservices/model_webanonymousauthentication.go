package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WebLinkedServiceTypeProperties = WebAnonymousAuthentication{}

type WebAnonymousAuthentication struct {

	// Fields inherited from WebLinkedServiceTypeProperties

	AuthenticationType WebAuthenticationType `json:"authenticationType"`
	Url                interface{}           `json:"url"`
}

func (s WebAnonymousAuthentication) WebLinkedServiceTypeProperties() BaseWebLinkedServiceTypePropertiesImpl {
	return BaseWebLinkedServiceTypePropertiesImpl{
		AuthenticationType: s.AuthenticationType,
		Url:                s.Url,
	}
}

var _ json.Marshaler = WebAnonymousAuthentication{}

func (s WebAnonymousAuthentication) MarshalJSON() ([]byte, error) {
	type wrapper WebAnonymousAuthentication
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebAnonymousAuthentication: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebAnonymousAuthentication: %+v", err)
	}

	decoded["authenticationType"] = "Anonymous"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebAnonymousAuthentication: %+v", err)
	}

	return encoded, nil
}
