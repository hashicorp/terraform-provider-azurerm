package maintenanceconfigurations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceScope string

const (
	MaintenanceScopeExtension          MaintenanceScope = "Extension"
	MaintenanceScopeHost               MaintenanceScope = "Host"
	MaintenanceScopeInGuestPatch       MaintenanceScope = "InGuestPatch"
	MaintenanceScopeOSImage            MaintenanceScope = "OSImage"
	MaintenanceScopeSQLDB              MaintenanceScope = "SQLDB"
	MaintenanceScopeSQLManagedInstance MaintenanceScope = "SQLManagedInstance"
)

func PossibleValuesForMaintenanceScope() []string {
	return []string{
		string(MaintenanceScopeExtension),
		string(MaintenanceScopeHost),
		string(MaintenanceScopeInGuestPatch),
		string(MaintenanceScopeOSImage),
		string(MaintenanceScopeSQLDB),
		string(MaintenanceScopeSQLManagedInstance),
	}
}

func parseMaintenanceScope(input string) (*MaintenanceScope, error) {
	vals := map[string]MaintenanceScope{
		"extension":          MaintenanceScopeExtension,
		"host":               MaintenanceScopeHost,
		"inguestpatch":       MaintenanceScopeInGuestPatch,
		"osimage":            MaintenanceScopeOSImage,
		"sqldb":              MaintenanceScopeSQLDB,
		"sqlmanagedinstance": MaintenanceScopeSQLManagedInstance,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MaintenanceScope(input)
	return &out, nil
}

type Visibility string

const (
	VisibilityCustom Visibility = "Custom"
	VisibilityPublic Visibility = "Public"
)

func PossibleValuesForVisibility() []string {
	return []string{
		string(VisibilityCustom),
		string(VisibilityPublic),
	}
}

func parseVisibility(input string) (*Visibility, error) {
	vals := map[string]Visibility{
		"custom": VisibilityCustom,
		"public": VisibilityPublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Visibility(input)
	return &out, nil
}
