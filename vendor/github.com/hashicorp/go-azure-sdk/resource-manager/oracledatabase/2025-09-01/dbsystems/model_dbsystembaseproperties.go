package dbsystems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbSystemBaseProperties interface {
	DbSystemBaseProperties() BaseDbSystemBasePropertiesImpl
}

var _ DbSystemBaseProperties = BaseDbSystemBasePropertiesImpl{}

type BaseDbSystemBasePropertiesImpl struct {
	ClusterName                  *string                         `json:"clusterName,omitempty"`
	ComputeCount                 *int64                          `json:"computeCount,omitempty"`
	ComputeModel                 *ComputeModel                   `json:"computeModel,omitempty"`
	DataStorageSizeInGbs         *int64                          `json:"dataStorageSizeInGbs,omitempty"`
	DbSystemOptions              *DbSystemOptions                `json:"dbSystemOptions,omitempty"`
	DiskRedundancy               *DiskRedundancyType             `json:"diskRedundancy,omitempty"`
	DisplayName                  *string                         `json:"displayName,omitempty"`
	Domain                       *string                         `json:"domain,omitempty"`
	GridImageOcid                *string                         `json:"gridImageOcid,omitempty"`
	Hostname                     string                          `json:"hostname"`
	InitialDataStorageSizeInGb   *int64                          `json:"initialDataStorageSizeInGb,omitempty"`
	LicenseModel                 *LicenseModel                   `json:"licenseModel,omitempty"`
	LifecycleDetails             *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState               *DbSystemLifecycleState         `json:"lifecycleState,omitempty"`
	ListenerPort                 *int64                          `json:"listenerPort,omitempty"`
	MemorySizeInGbs              *int64                          `json:"memorySizeInGbs,omitempty"`
	NetworkAnchorId              string                          `json:"networkAnchorId"`
	NodeCount                    *int64                          `json:"nodeCount,omitempty"`
	OciURL                       *string                         `json:"ociUrl,omitempty"`
	Ocid                         *string                         `json:"ocid,omitempty"`
	ProvisioningState            *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId             string                          `json:"resourceAnchorId"`
	ScanDnsName                  *string                         `json:"scanDnsName,omitempty"`
	ScanIPs                      *[]string                       `json:"scanIps,omitempty"`
	Shape                        string                          `json:"shape"`
	Source                       DbSystemSourceType              `json:"source"`
	SshPublicKeys                []string                        `json:"sshPublicKeys"`
	StorageVolumePerformanceMode *StorageVolumePerformanceMode   `json:"storageVolumePerformanceMode,omitempty"`
	TimeZone                     *string                         `json:"timeZone,omitempty"`
	Version                      *string                         `json:"version,omitempty"`
}

func (s BaseDbSystemBasePropertiesImpl) DbSystemBaseProperties() BaseDbSystemBasePropertiesImpl {
	return s
}

var _ DbSystemBaseProperties = RawDbSystemBasePropertiesImpl{}

// RawDbSystemBasePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDbSystemBasePropertiesImpl struct {
	dbSystemBaseProperties BaseDbSystemBasePropertiesImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawDbSystemBasePropertiesImpl) DbSystemBaseProperties() BaseDbSystemBasePropertiesImpl {
	return s.dbSystemBaseProperties
}

func UnmarshalDbSystemBasePropertiesImplementation(input []byte) (DbSystemBaseProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DbSystemBaseProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["source"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "None") {
		var out DbSystemProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DbSystemProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseDbSystemBasePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDbSystemBasePropertiesImpl: %+v", err)
	}

	return RawDbSystemBasePropertiesImpl{
		dbSystemBaseProperties: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
