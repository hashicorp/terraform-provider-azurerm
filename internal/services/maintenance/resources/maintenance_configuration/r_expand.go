// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandMaintenanceConfigurationWindow(input []any) *maintenanceconfigurations.MaintenanceWindow {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	startDateTime := v["start_date_time"].(string)
	expirationDateTime := v["expiration_date_time"].(string)
	duration := v["duration"].(string)
	timeZone := v["time_zone"].(string)
	recurEvery := v["recur_every"].(string)
	window := maintenanceconfigurations.MaintenanceWindow{
		StartDateTime:      pointer.To(startDateTime),
		ExpirationDateTime: pointer.To(expirationDateTime),
		Duration:           pointer.To(duration),
		TimeZone:           pointer.To(timeZone),
		RecurEvery:         pointer.To(recurEvery),
	}
	return &window
}

func expandMaintenanceConfigurationInstallPatches(input []any) *maintenanceconfigurations.InputPatchConfiguration {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]any)
	if !ok {
		return nil
	}
	installPatches := maintenanceconfigurations.InputPatchConfiguration{
		WindowsParameters: expandMaintenanceConfigurationInstallPatchesWindows(v["windows"].([]any)),
		LinuxParameters:   expandMaintenanceConfigurationInstallPatchesLinux(v["linux"].([]any)),
		RebootSetting:     pointer.ToEnum[maintenanceconfigurations.RebootOptions](v["reboot"].(string)),
	}
	return &installPatches
}

func expandMaintenanceConfigurationInstallPatchesWindows(input []any) *maintenanceconfigurations.InputWindowsParameters {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]any)
	if !ok {
		return nil
	}
	windowsInput := maintenanceconfigurations.InputWindowsParameters{}
	if v, ok := v["classifications_to_include"]; ok {
		windowsInput.ClassificationsToInclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	if v, ok := v["kb_numbers_to_exclude"]; ok {
		windowsInput.KbNumbersToExclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	if v, ok := v["kb_numbers_to_include"]; ok {
		windowsInput.KbNumbersToInclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	return &windowsInput
}

func expandMaintenanceConfigurationInstallPatchesLinux(input []any) *maintenanceconfigurations.InputLinuxParameters {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]any)
	if !ok {
		return nil
	}
	linuxParameters := maintenanceconfigurations.InputLinuxParameters{}
	if v, ok := v["classifications_to_include"]; ok {
		linuxParameters.ClassificationsToInclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	if v, ok := v["package_names_mask_to_exclude"]; ok {
		linuxParameters.PackageNameMasksToExclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	if v, ok := v["package_names_mask_to_include"]; ok {
		linuxParameters.PackageNameMasksToInclude = pluginsdk.ExpandStringSlice(v.([]any))
	}
	return &linuxParameters
}

func expandExtensionProperties(input map[string]any) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}
