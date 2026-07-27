package vaultcertificates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceCertificateDetails interface {
	ResourceCertificateDetails() BaseResourceCertificateDetailsImpl
}

var _ ResourceCertificateDetails = BaseResourceCertificateDetailsImpl{}

type BaseResourceCertificateDetailsImpl struct {
	AuthType     string  `json:"authType"`
	Certificate  *string `json:"certificate,omitempty"`
	FriendlyName *string `json:"friendlyName,omitempty"`
	Issuer       *string `json:"issuer,omitempty"`
	ResourceId   *int64  `json:"resourceId,omitempty"`
	Subject      *string `json:"subject,omitempty"`
	Thumbprint   *string `json:"thumbprint,omitempty"`
	ValidFrom    *string `json:"validFrom,omitempty"`
	ValidTo      *string `json:"validTo,omitempty"`
}

func (s BaseResourceCertificateDetailsImpl) ResourceCertificateDetails() BaseResourceCertificateDetailsImpl {
	return s
}

var _ ResourceCertificateDetails = RawResourceCertificateDetailsImpl{}

// RawResourceCertificateDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawResourceCertificateDetailsImpl struct {
	resourceCertificateDetails BaseResourceCertificateDetailsImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawResourceCertificateDetailsImpl) ResourceCertificateDetails() BaseResourceCertificateDetailsImpl {
	return s.resourceCertificateDetails
}

func (s RawResourceCertificateDetailsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalResourceCertificateDetailsImplementation(input []byte) (ResourceCertificateDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ResourceCertificateDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["authType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureActiveDirectory") {
		var out ResourceCertificateAndAadDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResourceCertificateAndAadDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AccessControlService") {
		var out ResourceCertificateAndAcsDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResourceCertificateAndAcsDetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseResourceCertificateDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseResourceCertificateDetailsImpl: %+v", err)
	}

	return RawResourceCertificateDetailsImpl{
		resourceCertificateDetails: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
