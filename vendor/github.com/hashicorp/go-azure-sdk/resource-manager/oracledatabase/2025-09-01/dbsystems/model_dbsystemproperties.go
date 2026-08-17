package dbsystems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DbSystemBaseProperties = DbSystemProperties{}

type DbSystemProperties struct {
	AdminPassword   *string                     `json:"adminPassword,omitempty"`
	DatabaseEdition DbSystemDatabaseEditionType `json:"databaseEdition"`
	DbVersion       string                      `json:"dbVersion"`
	PdbName         *string                     `json:"pdbName,omitempty"`

	// Fields inherited from DbSystemBaseProperties

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

func (s DbSystemProperties) DbSystemBaseProperties() BaseDbSystemBasePropertiesImpl {
	return BaseDbSystemBasePropertiesImpl{
		ClusterName:                  s.ClusterName,
		ComputeCount:                 s.ComputeCount,
		ComputeModel:                 s.ComputeModel,
		DataStorageSizeInGbs:         s.DataStorageSizeInGbs,
		DbSystemOptions:              s.DbSystemOptions,
		DiskRedundancy:               s.DiskRedundancy,
		DisplayName:                  s.DisplayName,
		Domain:                       s.Domain,
		GridImageOcid:                s.GridImageOcid,
		Hostname:                     s.Hostname,
		InitialDataStorageSizeInGb:   s.InitialDataStorageSizeInGb,
		LicenseModel:                 s.LicenseModel,
		LifecycleDetails:             s.LifecycleDetails,
		LifecycleState:               s.LifecycleState,
		ListenerPort:                 s.ListenerPort,
		MemorySizeInGbs:              s.MemorySizeInGbs,
		NetworkAnchorId:              s.NetworkAnchorId,
		NodeCount:                    s.NodeCount,
		OciURL:                       s.OciURL,
		Ocid:                         s.Ocid,
		ProvisioningState:            s.ProvisioningState,
		ResourceAnchorId:             s.ResourceAnchorId,
		ScanDnsName:                  s.ScanDnsName,
		ScanIPs:                      s.ScanIPs,
		Shape:                        s.Shape,
		Source:                       s.Source,
		SshPublicKeys:                s.SshPublicKeys,
		StorageVolumePerformanceMode: s.StorageVolumePerformanceMode,
		TimeZone:                     s.TimeZone,
		Version:                      s.Version,
	}
}

var _ json.Marshaler = DbSystemProperties{}

func (s DbSystemProperties) MarshalJSON() ([]byte, error) {
	type wrapper DbSystemProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DbSystemProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DbSystemProperties: %+v", err)
	}

	decoded["source"] = "None"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DbSystemProperties: %+v", err)
	}

	return encoded, nil
}
