// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenMaintenanceConfigurationWindow(input *maintenanceconfigurations.MaintenanceWindow) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		if startDateTime := v.StartDateTime; startDateTime != nil {
			output["start_date_time"] = *startDateTime
		}

		if expirationDateTime := v.ExpirationDateTime; expirationDateTime != nil {
			output["expiration_date_time"] = *expirationDateTime
		}

		if duration := v.Duration; duration != nil {
			output["duration"] = *duration
		}

		if timeZone := v.TimeZone; timeZone != nil {
			output["time_zone"] = *timeZone
		}

		if recurEvery := v.RecurEvery; recurEvery != nil {
			output["recur_every"] = *recurEvery
		}

		results = append(results, output)
	}

	return results
}

func flattenMaintenanceConfigurationInstallPatches(input *maintenanceconfigurations.InputPatchConfiguration) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		if rebootSetting := v.RebootSetting; rebootSetting != nil {
			// https://github.com/Azure/azure-rest-api-specs/issues/27222
			if strings.EqualFold(string(*rebootSetting), "AlwaysReboot") {
				output["reboot"] = "Always"
			} else {
				output["reboot"] = *rebootSetting
			}
		}

		if windowsParameters := v.WindowsParameters; windowsParameters != nil {
			output["windows"] = flattenMaintenanceConfigurationInstallPatchesWindows(windowsParameters)
		}

		if linuxParameters := v.LinuxParameters; linuxParameters != nil {
			output["linux"] = flattenMaintenanceConfigurationInstallPatchesLinux(linuxParameters)
		}

		results = append(results, output)
	}

	return results
}

func flattenMaintenanceConfigurationInstallPatchesWindows(input *maintenanceconfigurations.InputWindowsParameters) []any {
	results := make([]any, 0)

	if v := input; v != nil {
		output := make(map[string]any)

		if classificationsToInclude := v.ClassificationsToInclude; classificationsToInclude != nil {
			output["classifications_to_include"] = pluginsdk.FlattenSlice(classificationsToInclude)
		}

		if kbNumbersToExclude := v.KbNumbersToExclude; kbNumbersToExclude != nil {
			output["kb_numbers_to_exclude"] = pluginsdk.FlattenSlice(kbNumbersToExclude)
		}

		if kbNumbersToInclude := v.KbNumbersToInclude; kbNumbersToInclude != nil {
			output["kb_numbers_to_include"] = pluginsdk.FlattenSlice(kbNumbersToInclude)
		}

		results = append(results, output)
	}

	return results
}

func flattenMaintenanceConfigurationInstallPatchesLinux(input *maintenanceconfigurations.InputLinuxParameters) []any {
	results := make([]any, 0)

	if input != nil {
		classificationsToInclude := make([]any, 0)
		if input.ClassificationsToInclude != nil {
			classificationsToInclude = pluginsdk.FlattenSlice(input.ClassificationsToInclude)
		}
		packageNamesMaskToExclude := make([]any, 0)
		if input.PackageNameMasksToExclude != nil {
			packageNamesMaskToExclude = pluginsdk.FlattenSlice(input.PackageNameMasksToExclude)
		}
		packageNamesMaskToInclude := make([]any, 0)
		if input.PackageNameMasksToInclude != nil {
			packageNamesMaskToInclude = pluginsdk.FlattenSlice(input.PackageNameMasksToInclude)
		}

		results = append(results, map[string]any{
			"classifications_to_include":    classificationsToInclude,
			"package_names_mask_to_exclude": packageNamesMaskToExclude,
			"package_names_mask_to_include": packageNamesMaskToInclude,
		})
	}

	return results
}

func flattenExtensionProperties(input *map[string]string) map[string]any {
	output := make(map[string]any)
	if input != nil {
		for k, v := range *input {
			if k == "inGuestPatchMode" {
				output["InGuestPatchMode"] = v
				continue
			}
			output[k] = v
		}
	}
	return output
}
