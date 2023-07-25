package cluster

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPropertiesUpdateParameters struct {
	AddOnFeatures                        *[]AddOnFeatures                      `json:"addOnFeatures,omitempty"`
	ApplicationTypeVersionsCleanupPolicy *ApplicationTypeVersionsCleanupPolicy `json:"applicationTypeVersionsCleanupPolicy,omitempty"`
	Certificate                          *CertificateDescription               `json:"certificate,omitempty"`
	CertificateCommonNames               *ServerCertificateCommonNames         `json:"certificateCommonNames,omitempty"`
	ClientCertificateCommonNames         *[]ClientCertificateCommonName        `json:"clientCertificateCommonNames,omitempty"`
	ClientCertificateThumbprints         *[]ClientCertificateThumbprint        `json:"clientCertificateThumbprints,omitempty"`
	ClusterCodeVersion                   *string                               `json:"clusterCodeVersion,omitempty"`
	EventStoreServiceEnabled             *bool                                 `json:"eventStoreServiceEnabled,omitempty"`
	FabricSettings                       *[]SettingsSectionDescription         `json:"fabricSettings,omitempty"`
	InfrastructureServiceManager         *bool                                 `json:"infrastructureServiceManager,omitempty"`
	NodeTypes                            *[]NodeTypeDescription                `json:"nodeTypes,omitempty"`
	Notifications                        *[]Notification                       `json:"notifications,omitempty"`
	ReliabilityLevel                     *ReliabilityLevel                     `json:"reliabilityLevel,omitempty"`
	ReverseProxyCertificate              *CertificateDescription               `json:"reverseProxyCertificate,omitempty"`
	SfZonalUpgradeMode                   *SfZonalUpgradeMode                   `json:"sfZonalUpgradeMode,omitempty"`
	UpgradeDescription                   *ClusterUpgradePolicy                 `json:"upgradeDescription,omitempty"`
	UpgradeMode                          *UpgradeMode                          `json:"upgradeMode,omitempty"`
	UpgradePauseEndTimestampUtc          *string                               `json:"upgradePauseEndTimestampUtc,omitempty"`
	UpgradePauseStartTimestampUtc        *string                               `json:"upgradePauseStartTimestampUtc,omitempty"`
	UpgradeWave                          *ClusterUpgradeCadence                `json:"upgradeWave,omitempty"`
	VMSSZonalUpgradeMode                 *VMSSZonalUpgradeMode                 `json:"vmssZonalUpgradeMode,omitempty"`
	WaveUpgradePaused                    *bool                                 `json:"waveUpgradePaused,omitempty"`
}

func (o *ClusterPropertiesUpdateParameters) GetUpgradePauseEndTimestampUtcAsTime() (*time.Time, error) {
	if o.UpgradePauseEndTimestampUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpgradePauseEndTimestampUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterPropertiesUpdateParameters) SetUpgradePauseEndTimestampUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpgradePauseEndTimestampUtc = &formatted
}

func (o *ClusterPropertiesUpdateParameters) GetUpgradePauseStartTimestampUtcAsTime() (*time.Time, error) {
	if o.UpgradePauseStartTimestampUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpgradePauseStartTimestampUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterPropertiesUpdateParameters) SetUpgradePauseStartTimestampUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpgradePauseStartTimestampUtc = &formatted
}
