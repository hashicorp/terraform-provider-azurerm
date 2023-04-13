package streamingpoliciesandstreaminglocators

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionScheme string

const (
	EncryptionSchemeCommonEncryptionCbcs EncryptionScheme = "CommonEncryptionCbcs"
	EncryptionSchemeCommonEncryptionCenc EncryptionScheme = "CommonEncryptionCenc"
	EncryptionSchemeEnvelopeEncryption   EncryptionScheme = "EnvelopeEncryption"
	EncryptionSchemeNoEncryption         EncryptionScheme = "NoEncryption"
)

func PossibleValuesForEncryptionScheme() []string {
	return []string{
		string(EncryptionSchemeCommonEncryptionCbcs),
		string(EncryptionSchemeCommonEncryptionCenc),
		string(EncryptionSchemeEnvelopeEncryption),
		string(EncryptionSchemeNoEncryption),
	}
}

func (s *EncryptionScheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionScheme(input string) (*EncryptionScheme, error) {
	vals := map[string]EncryptionScheme{
		"commonencryptioncbcs": EncryptionSchemeCommonEncryptionCbcs,
		"commonencryptioncenc": EncryptionSchemeCommonEncryptionCenc,
		"envelopeencryption":   EncryptionSchemeEnvelopeEncryption,
		"noencryption":         EncryptionSchemeNoEncryption,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionScheme(input)
	return &out, nil
}

type StreamingLocatorContentKeyType string

const (
	StreamingLocatorContentKeyTypeCommonEncryptionCbcs StreamingLocatorContentKeyType = "CommonEncryptionCbcs"
	StreamingLocatorContentKeyTypeCommonEncryptionCenc StreamingLocatorContentKeyType = "CommonEncryptionCenc"
	StreamingLocatorContentKeyTypeEnvelopeEncryption   StreamingLocatorContentKeyType = "EnvelopeEncryption"
)

func PossibleValuesForStreamingLocatorContentKeyType() []string {
	return []string{
		string(StreamingLocatorContentKeyTypeCommonEncryptionCbcs),
		string(StreamingLocatorContentKeyTypeCommonEncryptionCenc),
		string(StreamingLocatorContentKeyTypeEnvelopeEncryption),
	}
}

func (s *StreamingLocatorContentKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStreamingLocatorContentKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStreamingLocatorContentKeyType(input string) (*StreamingLocatorContentKeyType, error) {
	vals := map[string]StreamingLocatorContentKeyType{
		"commonencryptioncbcs": StreamingLocatorContentKeyTypeCommonEncryptionCbcs,
		"commonencryptioncenc": StreamingLocatorContentKeyTypeCommonEncryptionCenc,
		"envelopeencryption":   StreamingLocatorContentKeyTypeEnvelopeEncryption,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StreamingLocatorContentKeyType(input)
	return &out, nil
}

type StreamingPolicyStreamingProtocol string

const (
	StreamingPolicyStreamingProtocolDash            StreamingPolicyStreamingProtocol = "Dash"
	StreamingPolicyStreamingProtocolDownload        StreamingPolicyStreamingProtocol = "Download"
	StreamingPolicyStreamingProtocolHls             StreamingPolicyStreamingProtocol = "Hls"
	StreamingPolicyStreamingProtocolSmoothStreaming StreamingPolicyStreamingProtocol = "SmoothStreaming"
)

func PossibleValuesForStreamingPolicyStreamingProtocol() []string {
	return []string{
		string(StreamingPolicyStreamingProtocolDash),
		string(StreamingPolicyStreamingProtocolDownload),
		string(StreamingPolicyStreamingProtocolHls),
		string(StreamingPolicyStreamingProtocolSmoothStreaming),
	}
}

func (s *StreamingPolicyStreamingProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStreamingPolicyStreamingProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStreamingPolicyStreamingProtocol(input string) (*StreamingPolicyStreamingProtocol, error) {
	vals := map[string]StreamingPolicyStreamingProtocol{
		"dash":            StreamingPolicyStreamingProtocolDash,
		"download":        StreamingPolicyStreamingProtocolDownload,
		"hls":             StreamingPolicyStreamingProtocolHls,
		"smoothstreaming": StreamingPolicyStreamingProtocolSmoothStreaming,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StreamingPolicyStreamingProtocol(input)
	return &out, nil
}

type TrackPropertyCompareOperation string

const (
	TrackPropertyCompareOperationEqual   TrackPropertyCompareOperation = "Equal"
	TrackPropertyCompareOperationUnknown TrackPropertyCompareOperation = "Unknown"
)

func PossibleValuesForTrackPropertyCompareOperation() []string {
	return []string{
		string(TrackPropertyCompareOperationEqual),
		string(TrackPropertyCompareOperationUnknown),
	}
}

func (s *TrackPropertyCompareOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrackPropertyCompareOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrackPropertyCompareOperation(input string) (*TrackPropertyCompareOperation, error) {
	vals := map[string]TrackPropertyCompareOperation{
		"equal":   TrackPropertyCompareOperationEqual,
		"unknown": TrackPropertyCompareOperationUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrackPropertyCompareOperation(input)
	return &out, nil
}

type TrackPropertyType string

const (
	TrackPropertyTypeFourCC  TrackPropertyType = "FourCC"
	TrackPropertyTypeUnknown TrackPropertyType = "Unknown"
)

func PossibleValuesForTrackPropertyType() []string {
	return []string{
		string(TrackPropertyTypeFourCC),
		string(TrackPropertyTypeUnknown),
	}
}

func (s *TrackPropertyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrackPropertyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrackPropertyType(input string) (*TrackPropertyType, error) {
	vals := map[string]TrackPropertyType{
		"fourcc":  TrackPropertyTypeFourCC,
		"unknown": TrackPropertyTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrackPropertyType(input)
	return &out, nil
}
