// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
)

func TestValidateNetAppVolumeGroupOracleVolumes(t *testing.T) {
	cases := []struct {
		Name        string
		VolumesData []volumegroups.VolumeGroupVolumeProperties
		Errors      int
	}{
		{
			Name: "ValidateCorrectSettingsAllVolumesPPG",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
				{ // data2
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData2))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData2)),
					},
				},
				{ // data3
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData3))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData3)),
					},
				},
				{ // data4
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData4))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData4)),
					},
				},
				{ // data5
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData5))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData5)),
					},
				},
				{ // data6
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData6))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData6)),
					},
				},
				{ // data7
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData7))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData7)),
					},
				},
				{ // data8
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData8))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData8)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
					},
				},
				{ // binary
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBinary))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleBinary)),
					},
				},
				{ // mirror
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleMirror))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleMirror)),
					},
				},
				{ // backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBackup))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleBackup)),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectSettingsAllVolumesAvailabilityZone",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData1)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data2
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData2))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData2)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data3
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData3))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData3)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data4
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData4))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData4)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data5
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData5))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData5)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data6
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData6))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData6)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data7
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData7))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData7)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // data8
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData8))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData8)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleLog)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // binary
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBinary))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleBinary)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // mirror
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleMirror))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleMirror)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBackup))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleBackup)),
					},
					Zones: pointer.To([]string{"1"}),
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectSettingsAllVolumesProximityPlacementGroup",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data2
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData2))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData2)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data3
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData3))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData3)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data4
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData4))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData4)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data5
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData5))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData5)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data6
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData6))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData6)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data7
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData7))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData7)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data8
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData8))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData8)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // binary
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBinary))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleBinary)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // mirror
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleMirror))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleMirror)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBackup))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleBackup)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectSettingsAllVolumesProximityPlacementGroupFailsWhenDifferent",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
					},
				},
				{ // data2
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData2))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData2)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg2"),
					},
				},
				{ // data3
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData3))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData3)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg3"),
					},
				},
				{ // data4
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData4))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData4)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg4"),
					},
				},
				{ // data5
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData5))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData5)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg5"),
					},
				},
				{ // data6
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData6))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData6)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg6"),
					},
				},
				{ // data7
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData7))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData7)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg7"),
					},
				},
				{ // data8
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData8))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData8)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg8"),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg9"),
					},
				},
				{ // binary
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBinary))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleBinary)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg10"),
					},
				},
				{ // mirror
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleMirror))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleMirror)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg11"),
					},
				},
				{ // backup
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBackup))),
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
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleBackup)),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg12"),
					},
				},
			},
			Errors: 11,
		},
		{
			Name: "ValidatePPGAndAvailabilityZoneNotSetAtSameTime",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
					},
					Zones: pointer.To([]string{"1"}),
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateMinimumVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateVolumesSameZone",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData1)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleLog)),
					},
					Zones: pointer.To([]string{"1"}),
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateVolumesNotSameZoneErrors",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleData1)),
					},
					Zones: pointer.To([]string{"1"}),
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleLog)),
					},
					Zones: pointer.To([]string{"2"}),
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateRequiredVolumeSpecs",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // mirror
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleMirror))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleMirror)),
					},
				},
				{ // binary
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleBinary))),
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
						VolumeSpecName: pointer.To(string(VolumeSpecNameOracleBinary)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateLessThanMinimumVolumes",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateMultiProtocolFails",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 5,
		},
		{
			Name: "ValidateNoProtocolFails",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateInvalidProtocolList",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 3,
		},
		{
			Name: "ValidateInvalidProtocol",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 3,
		},
		{
			Name: "ValidateCIFSInvalidProtocolForOracle",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
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
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateVolumeSpecCantRepeat",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
					},
				},
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
			},
			Errors: 1,
		},
		{
			Name: "ValidateEndpointDstNotEnabledOnLogVolume",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
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
			Name: "ValidateCustomerManagedKeyOptionsAreSetTogether",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleLog)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateCustomerManagedKeyOptionsNotSetTogetherFailsWhenKeyVault",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:     pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						NetworkFeatures:         pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleLog)),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateCustomerManagedKeyNotSameAllVolumesFails",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
				{ // data2
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData2))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData2)),
						EncryptionKeySource:     pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointNetApp),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleLog)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE1"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidatePlatformManagedKey",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:     pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointNetApp),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:           pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup: pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:           pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:          pointer.To(string(VolumeSpecNameOracleLog)),
						EncryptionKeySource:     pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointNetApp),
					},
				},
			},
			Errors: 0,
		},
		{
			Name: "ValidateCustomerManagedKeyFailsOnBasicNetworkFeatures",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesBasic),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleLog)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
					},
				},
			},
			Errors: 2,
		},
		{
			Name: "ValidateCustomerManagedKeyPassOnStandardNetworkFeatures",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data1
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleData1))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleData1)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
				{ // log
					Name: pointer.To(fmt.Sprintf("volume-%v", string(VolumeSpecNameOracleLog))),
					Properties: volumegroups.VolumeProperties{
						ProtocolTypes:                     pointer.To([]string{"NFSv4.1"}),
						ProximityPlacementGroup:           pointer.To("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/proximityPlacementGroups/ppg1"),
						SecurityStyle:                     pointer.To(volumegroups.SecurityStyleUnix),
						VolumeSpecName:                    pointer.To(string(VolumeSpecNameOracleLog)),
						EncryptionKeySource:               pointer.To(volumegroups.EncryptionKeySourceMicrosoftPointKeyVault),
						KeyVaultPrivateEndpointResourceId: pointer.To("/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.Network/privateEndpoints/myKeyvaultPE"),
						NetworkFeatures:                   pointer.To(volumegroups.NetworkFeaturesStandard),
					},
				},
			},
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeGroupOracleVolumes(pointer.To(tc.VolumesData))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateNetAppVolumeGroupOracleVolumes to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
