package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HDInsightOnDemandLinkedServiceTypeProperties struct {
	AdditionalLinkedServiceNames *[]LinkedServiceReference `json:"additionalLinkedServiceNames,omitempty"`
	ClusterNamePrefix            *interface{}              `json:"clusterNamePrefix,omitempty"`
	ClusterPassword              SecretBase                `json:"clusterPassword"`
	ClusterResourceGroup         interface{}               `json:"clusterResourceGroup"`
	ClusterSize                  int64                     `json:"clusterSize"`
	ClusterSshPassword           SecretBase                `json:"clusterSshPassword"`
	ClusterSshUserName           *interface{}              `json:"clusterSshUserName,omitempty"`
	ClusterType                  *interface{}              `json:"clusterType,omitempty"`
	ClusterUserName              *interface{}              `json:"clusterUserName,omitempty"`
	CoreConfiguration            *interface{}              `json:"coreConfiguration,omitempty"`
	Credential                   *CredentialReference      `json:"credential,omitempty"`
	DataNodeSize                 *interface{}              `json:"dataNodeSize,omitempty"`
	EncryptedCredential          *string                   `json:"encryptedCredential,omitempty"`
	HBaseConfiguration           *interface{}              `json:"hBaseConfiguration,omitempty"`
	HcatalogLinkedServiceName    *LinkedServiceReference   `json:"hcatalogLinkedServiceName,omitempty"`
	HdfsConfiguration            *interface{}              `json:"hdfsConfiguration,omitempty"`
	HeadNodeSize                 *interface{}              `json:"headNodeSize,omitempty"`
	HiveConfiguration            *interface{}              `json:"hiveConfiguration,omitempty"`
	HostSubscriptionId           interface{}               `json:"hostSubscriptionId"`
	LinkedServiceName            LinkedServiceReference    `json:"linkedServiceName"`
	MapReduceConfiguration       *interface{}              `json:"mapReduceConfiguration,omitempty"`
	OozieConfiguration           *interface{}              `json:"oozieConfiguration,omitempty"`
	ScriptActions                *[]ScriptAction           `json:"scriptActions,omitempty"`
	ServicePrincipalId           *interface{}              `json:"servicePrincipalId,omitempty"`
	ServicePrincipalKey          SecretBase                `json:"servicePrincipalKey"`
	SparkVersion                 *interface{}              `json:"sparkVersion,omitempty"`
	StormConfiguration           *interface{}              `json:"stormConfiguration,omitempty"`
	SubnetName                   *interface{}              `json:"subnetName,omitempty"`
	Tenant                       interface{}               `json:"tenant"`
	TimeToLive                   interface{}               `json:"timeToLive"`
	Version                      interface{}               `json:"version"`
	VirtualNetworkId             *interface{}              `json:"virtualNetworkId,omitempty"`
	YarnConfiguration            *interface{}              `json:"yarnConfiguration,omitempty"`
	ZookeeperNodeSize            *interface{}              `json:"zookeeperNodeSize,omitempty"`
}

var _ json.Unmarshaler = &HDInsightOnDemandLinkedServiceTypeProperties{}

func (s *HDInsightOnDemandLinkedServiceTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AdditionalLinkedServiceNames *[]LinkedServiceReference `json:"additionalLinkedServiceNames,omitempty"`
		ClusterNamePrefix            *interface{}              `json:"clusterNamePrefix,omitempty"`
		ClusterResourceGroup         interface{}               `json:"clusterResourceGroup"`
		ClusterSize                  int64                     `json:"clusterSize"`
		ClusterSshUserName           *interface{}              `json:"clusterSshUserName,omitempty"`
		ClusterType                  *interface{}              `json:"clusterType,omitempty"`
		ClusterUserName              *interface{}              `json:"clusterUserName,omitempty"`
		CoreConfiguration            *interface{}              `json:"coreConfiguration,omitempty"`
		Credential                   *CredentialReference      `json:"credential,omitempty"`
		DataNodeSize                 *interface{}              `json:"dataNodeSize,omitempty"`
		EncryptedCredential          *string                   `json:"encryptedCredential,omitempty"`
		HBaseConfiguration           *interface{}              `json:"hBaseConfiguration,omitempty"`
		HcatalogLinkedServiceName    *LinkedServiceReference   `json:"hcatalogLinkedServiceName,omitempty"`
		HdfsConfiguration            *interface{}              `json:"hdfsConfiguration,omitempty"`
		HeadNodeSize                 *interface{}              `json:"headNodeSize,omitempty"`
		HiveConfiguration            *interface{}              `json:"hiveConfiguration,omitempty"`
		HostSubscriptionId           interface{}               `json:"hostSubscriptionId"`
		LinkedServiceName            LinkedServiceReference    `json:"linkedServiceName"`
		MapReduceConfiguration       *interface{}              `json:"mapReduceConfiguration,omitempty"`
		OozieConfiguration           *interface{}              `json:"oozieConfiguration,omitempty"`
		ScriptActions                *[]ScriptAction           `json:"scriptActions,omitempty"`
		ServicePrincipalId           *interface{}              `json:"servicePrincipalId,omitempty"`
		SparkVersion                 *interface{}              `json:"sparkVersion,omitempty"`
		StormConfiguration           *interface{}              `json:"stormConfiguration,omitempty"`
		SubnetName                   *interface{}              `json:"subnetName,omitempty"`
		Tenant                       interface{}               `json:"tenant"`
		TimeToLive                   interface{}               `json:"timeToLive"`
		Version                      interface{}               `json:"version"`
		VirtualNetworkId             *interface{}              `json:"virtualNetworkId,omitempty"`
		YarnConfiguration            *interface{}              `json:"yarnConfiguration,omitempty"`
		ZookeeperNodeSize            *interface{}              `json:"zookeeperNodeSize,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AdditionalLinkedServiceNames = decoded.AdditionalLinkedServiceNames
	s.ClusterNamePrefix = decoded.ClusterNamePrefix
	s.ClusterResourceGroup = decoded.ClusterResourceGroup
	s.ClusterSize = decoded.ClusterSize
	s.ClusterSshUserName = decoded.ClusterSshUserName
	s.ClusterType = decoded.ClusterType
	s.ClusterUserName = decoded.ClusterUserName
	s.CoreConfiguration = decoded.CoreConfiguration
	s.Credential = decoded.Credential
	s.DataNodeSize = decoded.DataNodeSize
	s.EncryptedCredential = decoded.EncryptedCredential
	s.HBaseConfiguration = decoded.HBaseConfiguration
	s.HcatalogLinkedServiceName = decoded.HcatalogLinkedServiceName
	s.HdfsConfiguration = decoded.HdfsConfiguration
	s.HeadNodeSize = decoded.HeadNodeSize
	s.HiveConfiguration = decoded.HiveConfiguration
	s.HostSubscriptionId = decoded.HostSubscriptionId
	s.LinkedServiceName = decoded.LinkedServiceName
	s.MapReduceConfiguration = decoded.MapReduceConfiguration
	s.OozieConfiguration = decoded.OozieConfiguration
	s.ScriptActions = decoded.ScriptActions
	s.ServicePrincipalId = decoded.ServicePrincipalId
	s.SparkVersion = decoded.SparkVersion
	s.StormConfiguration = decoded.StormConfiguration
	s.SubnetName = decoded.SubnetName
	s.Tenant = decoded.Tenant
	s.TimeToLive = decoded.TimeToLive
	s.Version = decoded.Version
	s.VirtualNetworkId = decoded.VirtualNetworkId
	s.YarnConfiguration = decoded.YarnConfiguration
	s.ZookeeperNodeSize = decoded.ZookeeperNodeSize

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling HDInsightOnDemandLinkedServiceTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["clusterPassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClusterPassword' for 'HDInsightOnDemandLinkedServiceTypeProperties': %+v", err)
		}
		s.ClusterPassword = impl
	}

	if v, ok := temp["clusterSshPassword"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ClusterSshPassword' for 'HDInsightOnDemandLinkedServiceTypeProperties': %+v", err)
		}
		s.ClusterSshPassword = impl
	}

	if v, ok := temp["servicePrincipalKey"]; ok {
		impl, err := UnmarshalSecretBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ServicePrincipalKey' for 'HDInsightOnDemandLinkedServiceTypeProperties': %+v", err)
		}
		s.ServicePrincipalKey = impl
	}

	return nil
}
