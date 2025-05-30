// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
)

func TestValidateNetAppVolumeGroupSAPHanaVolumes(t *testing.T) {
	cases := []struct {
		Name        string
		VolumesData []volumegroups.VolumeGroupVolumeProperties
		Errors      int
	}{
		{
			Name: "ValidateCorrectSettingsAllVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
				{ // shared
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaShared))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaShared)),
					},
				},
				{ // data-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaDataBackup))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaDataBackup)),
					},
				},
				{ // log-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLogBackup))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaLogBackup)),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateMinimumVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateRequiredVolumeSpecs",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // shared
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaShared))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaShared)),
					},
				},
				{ // data-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaDataBackup))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaDataBackup)),
					},
				},
				{ // log-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLogBackup))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaLogBackup)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateLessThanMinimumVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateMultiProtocolFails",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(true),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1", "NFSv3"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 6,
		},
		{
			Name: "ValidateNoProtocolFails",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateInvalidProtocolList",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"NFSv4.1", "InvalidProtocol"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 4,
		},
		{
			Name: "ValidateInvalidProtocol",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"InvalidProtocol"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 3,
		},
		{
			Name: "ValidateCIFSInvalidProtocolForSAPHana",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ExportPolicy: &volumegroups.VolumePropertiesExportPolicy{
							Rules: &[]volumegroups.ExportPolicyRule{
								{
									Nfsv3:  pointer.To(false),
									Nfsv41: pointer.To(true),
								},
							},
						},
						ProtocolTypes:           pointer.To([]string{"CIFS"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateNoNfsVersionThreeOnDataLogAndSharedVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv3"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv3"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
				{ // shared
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaShared))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv3"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaShared)),
					},
				},
			},
			Errors: 3,
		},
		{
			Name: "ValidateNoPPGBackupVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
				{ // data-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaDataBackup))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaDataBackup)),
					},
				},
				{ // log-backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLogBackup))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLogBackup)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateRequiredPpgForNonBackupVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
				{ // shared
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaShared))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:  pointer.To([]string{"NFSv4.1"}),
						SecurityStyle:  pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName: pointer.To(string(VolumeSpecNameSAPHanaShared)),
					},
				},
			},
			Errors: 3,
		},
		{
			Name: "ValidateVolumeSpecCantRepeat",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
				{ // shared
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaShared))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateEndpointDstNotEnabledOnLogVolume",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
						DataProtection: &volumegroups.VolumePropertiesDataProtection{
							Replication: &volumegroups.ReplicationObject{
								EndpointType: pointer.To(volumegroups.EndpointTypeDst),
							},
						},
					},
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateSnapshotPolicyNotEnabledOnEndpointDstVolume",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaData))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaData)),
						DataProtection: &volumegroups.VolumePropertiesDataProtection{
							Replication: &volumegroups.ReplicationObject{
								EndpointType: pointer.To(volumegroups.EndpointTypeDst),
							},
							Snapshot: &volumegroups.VolumeSnapshotProperties{
								SnapshotPolicyId: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/volume1/snapshotPolicies/snapshotPolicy1"),
							},
						},
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameSAPHanaLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameSAPHanaLog)),
					},
				},
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeGroupSAPHanaVolumes(pointer.To(tc.VolumesData))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateNetAppVolumeGroupSAPHanaVolumes to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
