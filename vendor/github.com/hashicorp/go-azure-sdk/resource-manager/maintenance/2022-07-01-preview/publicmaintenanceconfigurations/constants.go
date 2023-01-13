package publicmaintenanceconfigurations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceScope string

const (
	MaintenanceScopeExtension          MaintenanceScope = "Extension"
	MaintenanceScopeHost               MaintenanceScope = "Host"
	MaintenanceScopeInGuestPatch       MaintenanceScope = "InGuestPatch"
	MaintenanceScopeOSImage            MaintenanceScope = "OSImage"
	MaintenanceScopeResource           MaintenanceScope = "Resource"
	MaintenanceScopeSQLDB              MaintenanceScope = "SQLDB"
	MaintenanceScopeSQLManagedInstance MaintenanceScope = "SQLManagedInstance"
)

func PossibleValuesForMaintenanceScope() []string {
	return []string{
		string(MaintenanceScopeExtension),
		string(MaintenanceScopeHost),
		string(MaintenanceScopeInGuestPatch),
		string(MaintenanceScopeOSImage),
		string(MaintenanceScopeResource),
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
		"resource":           MaintenanceScopeResource,
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

type RebootOptions string

const (
	RebootOptionsAlways     RebootOptions = "Always"
	RebootOptionsIfRequired RebootOptions = "IfRequired"
	RebootOptionsNever      RebootOptions = "Never"
)

func PossibleValuesForRebootOptions() []string {
	return []string{
		string(RebootOptionsAlways),
		string(RebootOptionsIfRequired),
		string(RebootOptionsNever),
	}
}

func parseRebootOptions(input string) (*RebootOptions, error) {
	vals := map[string]RebootOptions{
		"always":     RebootOptionsAlways,
		"ifrequired": RebootOptionsIfRequired,
		"never":      RebootOptionsNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RebootOptions(input)
	return &out, nil
}

type TaskScope string

const (
	TaskScopeGlobal   TaskScope = "Global"
	TaskScopeResource TaskScope = "Resource"
)

func PossibleValuesForTaskScope() []string {
	return []string{
		string(TaskScopeGlobal),
		string(TaskScopeResource),
	}
}

func parseTaskScope(input string) (*TaskScope, error) {
	vals := map[string]TaskScope{
		"global":   TaskScopeGlobal,
		"resource": TaskScopeResource,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TaskScope(input)
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
