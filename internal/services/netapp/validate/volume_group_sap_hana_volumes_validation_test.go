package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumegroups"
)

func TestValidateNetAppVolumeGroupSAPHanaVolumes(t *testing.T) {
	cases := []struct {
		Name        string
		VolumesData []volumegroups.VolumeGroupVolumeProperties
		Errors      int
	}{
		{
			Name:        "ValidateCorrectSettings",
			VolumesData: []volumegroups.VolumeGroupVolumeProperties{
				{ // data
					Properties: volumeGroups
					ExportPolicy: &volumegroups.ExportPolicy{
						Rules: []volumegroups.ExportPolicyRule{
							{
								RuleIndex: pointer.ToInt32(1),
								AllowedClients: pointer.ToString("0.0.0.0/0"),
								Cifs: pointer.ToBool(false),
								Nfsv3: pointer.ToBool(false),
								Nfsv41: pointer.ToBool(true),
								UnixReadOnly: pointer.ToBool(false),
								UnixReadWrite: pointer.ToBool(true),
								HasRootAccess: pointer.ToBoll(true),
							},
						},
					},
				},
			},

			Errors:      0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeGroupSAPHanaVolumes(pointer.To(tc.VolumesData))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateNetAppVolumeGroupSAPHanaVolumes to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
