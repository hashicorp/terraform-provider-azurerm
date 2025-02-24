package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISPackageLocationTypeProperties struct {
	AccessCredential              *SSISAccessCredential `json:"accessCredential,omitempty"`
	ChildPackages                 *[]SSISChildPackage   `json:"childPackages,omitempty"`
	ConfigurationAccessCredential *SSISAccessCredential `json:"configurationAccessCredential,omitempty"`
	ConfigurationPath             *string               `json:"configurationPath,omitempty"`
	PackageContent                *string               `json:"packageContent,omitempty"`
	PackageLastModifiedDate       *string               `json:"packageLastModifiedDate,omitempty"`
	PackageName                   *string               `json:"packageName,omitempty"`
	PackagePassword               SecretBase            `json:"packagePassword"`
}

var _ json.Unmarshaler = &SSISPackageLocationTypeProperties{}

func (s *SSISPackageLocationTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AccessCredential              *SSISAccessCredential `json:"accessCredential,omitempty"`
		ChildPackages                 *[]SSISChildPackage   `json:"childPackages,omitempty"`
		ConfigurationAccessCredential *SSISAccessCredential `json:"configurationAccessCredential,omitempty"`
		ConfigurationPath             *string               `json:"configurationPath,omitempty"`
		PackageContent                *string               `json:"packageContent,omitempty"`
		PackageLastModifiedDate       *string               `json:"packageLastModifiedDate,omitempty"`
		PackageName                   *string               `json:"packageName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AccessCredential = decoded.AccessCredential
	s.ChildPackages = decoded.ChildPackages
	s.ConfigurationAccessCredential = decoded.ConfigurationAccessCredential
	s.ConfigurationPath = decoded.ConfigurationPath
	s.PackageContent = decoded.PackageContent
	s.PackageLastModifiedDate = decoded.PackageLastModifiedDate
	s.PackageName = decoded.PackageName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SSISPackageLocationTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["packagePassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PackagePassword' for 'SSISPackageLocationTypeProperties': %+v", err)
		}
		s.PackagePassword = impl
	}

	return nil
}
