package contentkeypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyFairPlayRentalAndLeaseKeyType string

const (
	ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeDualExpiry          ContentKeyPolicyFairPlayRentalAndLeaseKeyType = "DualExpiry"
	ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentLimited   ContentKeyPolicyFairPlayRentalAndLeaseKeyType = "PersistentLimited"
	ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentUnlimited ContentKeyPolicyFairPlayRentalAndLeaseKeyType = "PersistentUnlimited"
	ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUndefined           ContentKeyPolicyFairPlayRentalAndLeaseKeyType = "Undefined"
	ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUnknown             ContentKeyPolicyFairPlayRentalAndLeaseKeyType = "Unknown"
)

func PossibleValuesForContentKeyPolicyFairPlayRentalAndLeaseKeyType() []string {
	return []string{
		string(ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeDualExpiry),
		string(ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentLimited),
		string(ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentUnlimited),
		string(ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUndefined),
		string(ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUnknown),
	}
}

func (s *ContentKeyPolicyFairPlayRentalAndLeaseKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentKeyPolicyFairPlayRentalAndLeaseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentKeyPolicyFairPlayRentalAndLeaseKeyType(input string) (*ContentKeyPolicyFairPlayRentalAndLeaseKeyType, error) {
	vals := map[string]ContentKeyPolicyFairPlayRentalAndLeaseKeyType{
		"dualexpiry":          ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeDualExpiry,
		"persistentlimited":   ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentLimited,
		"persistentunlimited": ContentKeyPolicyFairPlayRentalAndLeaseKeyTypePersistentUnlimited,
		"undefined":           ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUndefined,
		"unknown":             ContentKeyPolicyFairPlayRentalAndLeaseKeyTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentKeyPolicyFairPlayRentalAndLeaseKeyType(input)
	return &out, nil
}

type ContentKeyPolicyPlayReadyContentType string

const (
	ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload  ContentKeyPolicyPlayReadyContentType = "UltraVioletDownload"
	ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming ContentKeyPolicyPlayReadyContentType = "UltraVioletStreaming"
	ContentKeyPolicyPlayReadyContentTypeUnknown              ContentKeyPolicyPlayReadyContentType = "Unknown"
	ContentKeyPolicyPlayReadyContentTypeUnspecified          ContentKeyPolicyPlayReadyContentType = "Unspecified"
)

func PossibleValuesForContentKeyPolicyPlayReadyContentType() []string {
	return []string{
		string(ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload),
		string(ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming),
		string(ContentKeyPolicyPlayReadyContentTypeUnknown),
		string(ContentKeyPolicyPlayReadyContentTypeUnspecified),
	}
}

func (s *ContentKeyPolicyPlayReadyContentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentKeyPolicyPlayReadyContentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentKeyPolicyPlayReadyContentType(input string) (*ContentKeyPolicyPlayReadyContentType, error) {
	vals := map[string]ContentKeyPolicyPlayReadyContentType{
		"ultravioletdownload":  ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload,
		"ultravioletstreaming": ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming,
		"unknown":              ContentKeyPolicyPlayReadyContentTypeUnknown,
		"unspecified":          ContentKeyPolicyPlayReadyContentTypeUnspecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentKeyPolicyPlayReadyContentType(input)
	return &out, nil
}

type ContentKeyPolicyPlayReadyLicenseType string

const (
	ContentKeyPolicyPlayReadyLicenseTypeNonPersistent ContentKeyPolicyPlayReadyLicenseType = "NonPersistent"
	ContentKeyPolicyPlayReadyLicenseTypePersistent    ContentKeyPolicyPlayReadyLicenseType = "Persistent"
	ContentKeyPolicyPlayReadyLicenseTypeUnknown       ContentKeyPolicyPlayReadyLicenseType = "Unknown"
)

func PossibleValuesForContentKeyPolicyPlayReadyLicenseType() []string {
	return []string{
		string(ContentKeyPolicyPlayReadyLicenseTypeNonPersistent),
		string(ContentKeyPolicyPlayReadyLicenseTypePersistent),
		string(ContentKeyPolicyPlayReadyLicenseTypeUnknown),
	}
}

func (s *ContentKeyPolicyPlayReadyLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentKeyPolicyPlayReadyLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentKeyPolicyPlayReadyLicenseType(input string) (*ContentKeyPolicyPlayReadyLicenseType, error) {
	vals := map[string]ContentKeyPolicyPlayReadyLicenseType{
		"nonpersistent": ContentKeyPolicyPlayReadyLicenseTypeNonPersistent,
		"persistent":    ContentKeyPolicyPlayReadyLicenseTypePersistent,
		"unknown":       ContentKeyPolicyPlayReadyLicenseTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentKeyPolicyPlayReadyLicenseType(input)
	return &out, nil
}

type ContentKeyPolicyPlayReadyUnknownOutputPassingOption string

const (
	ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowed                      ContentKeyPolicyPlayReadyUnknownOutputPassingOption = "Allowed"
	ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowedWithVideoConstriction ContentKeyPolicyPlayReadyUnknownOutputPassingOption = "AllowedWithVideoConstriction"
	ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed                   ContentKeyPolicyPlayReadyUnknownOutputPassingOption = "NotAllowed"
	ContentKeyPolicyPlayReadyUnknownOutputPassingOptionUnknown                      ContentKeyPolicyPlayReadyUnknownOutputPassingOption = "Unknown"
)

func PossibleValuesForContentKeyPolicyPlayReadyUnknownOutputPassingOption() []string {
	return []string{
		string(ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowed),
		string(ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowedWithVideoConstriction),
		string(ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed),
		string(ContentKeyPolicyPlayReadyUnknownOutputPassingOptionUnknown),
	}
}

func (s *ContentKeyPolicyPlayReadyUnknownOutputPassingOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentKeyPolicyPlayReadyUnknownOutputPassingOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentKeyPolicyPlayReadyUnknownOutputPassingOption(input string) (*ContentKeyPolicyPlayReadyUnknownOutputPassingOption, error) {
	vals := map[string]ContentKeyPolicyPlayReadyUnknownOutputPassingOption{
		"allowed":                      ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowed,
		"allowedwithvideoconstriction": ContentKeyPolicyPlayReadyUnknownOutputPassingOptionAllowedWithVideoConstriction,
		"notallowed":                   ContentKeyPolicyPlayReadyUnknownOutputPassingOptionNotAllowed,
		"unknown":                      ContentKeyPolicyPlayReadyUnknownOutputPassingOptionUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentKeyPolicyPlayReadyUnknownOutputPassingOption(input)
	return &out, nil
}

type ContentKeyPolicyRestrictionTokenType string

const (
	ContentKeyPolicyRestrictionTokenTypeJwt     ContentKeyPolicyRestrictionTokenType = "Jwt"
	ContentKeyPolicyRestrictionTokenTypeSwt     ContentKeyPolicyRestrictionTokenType = "Swt"
	ContentKeyPolicyRestrictionTokenTypeUnknown ContentKeyPolicyRestrictionTokenType = "Unknown"
)

func PossibleValuesForContentKeyPolicyRestrictionTokenType() []string {
	return []string{
		string(ContentKeyPolicyRestrictionTokenTypeJwt),
		string(ContentKeyPolicyRestrictionTokenTypeSwt),
		string(ContentKeyPolicyRestrictionTokenTypeUnknown),
	}
}

func (s *ContentKeyPolicyRestrictionTokenType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContentKeyPolicyRestrictionTokenType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContentKeyPolicyRestrictionTokenType(input string) (*ContentKeyPolicyRestrictionTokenType, error) {
	vals := map[string]ContentKeyPolicyRestrictionTokenType{
		"jwt":     ContentKeyPolicyRestrictionTokenTypeJwt,
		"swt":     ContentKeyPolicyRestrictionTokenTypeSwt,
		"unknown": ContentKeyPolicyRestrictionTokenTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContentKeyPolicyRestrictionTokenType(input)
	return &out, nil
}
