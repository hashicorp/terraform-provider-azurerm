package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestriction = ContentKeyPolicyTokenRestriction{}

type ContentKeyPolicyTokenRestriction struct {
	AlternateVerificationKeys      *[]ContentKeyPolicyRestrictionTokenKey `json:"alternateVerificationKeys,omitempty"`
	Audience                       string                                 `json:"audience"`
	Issuer                         string                                 `json:"issuer"`
	OpenIdConnectDiscoveryDocument *string                                `json:"openIdConnectDiscoveryDocument,omitempty"`
	PrimaryVerificationKey         ContentKeyPolicyRestrictionTokenKey    `json:"primaryVerificationKey"`
	RequiredClaims                 *[]ContentKeyPolicyTokenClaim          `json:"requiredClaims,omitempty"`
	RestrictionTokenType           ContentKeyPolicyRestrictionTokenType   `json:"restrictionTokenType"`

	// Fields inherited from ContentKeyPolicyRestriction
}

var _ json.Marshaler = ContentKeyPolicyTokenRestriction{}

func (s ContentKeyPolicyTokenRestriction) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyTokenRestriction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyTokenRestriction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyTokenRestriction: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyTokenRestriction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyTokenRestriction: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ContentKeyPolicyTokenRestriction{}

func (s *ContentKeyPolicyTokenRestriction) UnmarshalJSON(bytes []byte) error {
	type alias ContentKeyPolicyTokenRestriction
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ContentKeyPolicyTokenRestriction: %+v", err)
	}

	s.Audience = decoded.Audience
	s.Issuer = decoded.Issuer
	s.OpenIdConnectDiscoveryDocument = decoded.OpenIdConnectDiscoveryDocument
	s.RequiredClaims = decoded.RequiredClaims
	s.RestrictionTokenType = decoded.RestrictionTokenType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ContentKeyPolicyTokenRestriction into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["alternateVerificationKeys"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling AlternateVerificationKeys into list []json.RawMessage: %+v", err)
		}

		output := make([]ContentKeyPolicyRestrictionTokenKey, 0)
		for i, val := range listTemp {
			impl, err := unmarshalContentKeyPolicyRestrictionTokenKeyImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'AlternateVerificationKeys' for 'ContentKeyPolicyTokenRestriction': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.AlternateVerificationKeys = &output
	}

	if v, ok := temp["primaryVerificationKey"]; ok {
		impl, err := unmarshalContentKeyPolicyRestrictionTokenKeyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PrimaryVerificationKey' for 'ContentKeyPolicyTokenRestriction': %+v", err)
		}
		s.PrimaryVerificationKey = impl
	}
	return nil
}
