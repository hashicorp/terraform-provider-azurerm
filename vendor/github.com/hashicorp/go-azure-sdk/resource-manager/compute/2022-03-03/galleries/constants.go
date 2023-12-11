package galleries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryExpandParams string

const (
	GalleryExpandParamsSharingProfileGroups GalleryExpandParams = "SharingProfile/Groups"
)

func PossibleValuesForGalleryExpandParams() []string {
	return []string{
		string(GalleryExpandParamsSharingProfileGroups),
	}
}

func (s *GalleryExpandParams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGalleryExpandParams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGalleryExpandParams(input string) (*GalleryExpandParams, error) {
	vals := map[string]GalleryExpandParams{
		"sharingprofile/groups": GalleryExpandParamsSharingProfileGroups,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryExpandParams(input)
	return &out, nil
}

type GalleryProvisioningState string

const (
	GalleryProvisioningStateCreating  GalleryProvisioningState = "Creating"
	GalleryProvisioningStateDeleting  GalleryProvisioningState = "Deleting"
	GalleryProvisioningStateFailed    GalleryProvisioningState = "Failed"
	GalleryProvisioningStateMigrating GalleryProvisioningState = "Migrating"
	GalleryProvisioningStateSucceeded GalleryProvisioningState = "Succeeded"
	GalleryProvisioningStateUpdating  GalleryProvisioningState = "Updating"
)

func PossibleValuesForGalleryProvisioningState() []string {
	return []string{
		string(GalleryProvisioningStateCreating),
		string(GalleryProvisioningStateDeleting),
		string(GalleryProvisioningStateFailed),
		string(GalleryProvisioningStateMigrating),
		string(GalleryProvisioningStateSucceeded),
		string(GalleryProvisioningStateUpdating),
	}
}

func (s *GalleryProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGalleryProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGalleryProvisioningState(input string) (*GalleryProvisioningState, error) {
	vals := map[string]GalleryProvisioningState{
		"creating":  GalleryProvisioningStateCreating,
		"deleting":  GalleryProvisioningStateDeleting,
		"failed":    GalleryProvisioningStateFailed,
		"migrating": GalleryProvisioningStateMigrating,
		"succeeded": GalleryProvisioningStateSucceeded,
		"updating":  GalleryProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryProvisioningState(input)
	return &out, nil
}

type GallerySharingPermissionTypes string

const (
	GallerySharingPermissionTypesCommunity GallerySharingPermissionTypes = "Community"
	GallerySharingPermissionTypesGroups    GallerySharingPermissionTypes = "Groups"
	GallerySharingPermissionTypesPrivate   GallerySharingPermissionTypes = "Private"
)

func PossibleValuesForGallerySharingPermissionTypes() []string {
	return []string{
		string(GallerySharingPermissionTypesCommunity),
		string(GallerySharingPermissionTypesGroups),
		string(GallerySharingPermissionTypesPrivate),
	}
}

func (s *GallerySharingPermissionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGallerySharingPermissionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGallerySharingPermissionTypes(input string) (*GallerySharingPermissionTypes, error) {
	vals := map[string]GallerySharingPermissionTypes{
		"community": GallerySharingPermissionTypesCommunity,
		"groups":    GallerySharingPermissionTypesGroups,
		"private":   GallerySharingPermissionTypesPrivate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GallerySharingPermissionTypes(input)
	return &out, nil
}

type SelectPermissions string

const (
	SelectPermissionsPermissions SelectPermissions = "Permissions"
)

func PossibleValuesForSelectPermissions() []string {
	return []string{
		string(SelectPermissionsPermissions),
	}
}

func (s *SelectPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSelectPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSelectPermissions(input string) (*SelectPermissions, error) {
	vals := map[string]SelectPermissions{
		"permissions": SelectPermissionsPermissions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SelectPermissions(input)
	return &out, nil
}

type SharingProfileGroupTypes string

const (
	SharingProfileGroupTypesAADTenants    SharingProfileGroupTypes = "AADTenants"
	SharingProfileGroupTypesSubscriptions SharingProfileGroupTypes = "Subscriptions"
)

func PossibleValuesForSharingProfileGroupTypes() []string {
	return []string{
		string(SharingProfileGroupTypesAADTenants),
		string(SharingProfileGroupTypesSubscriptions),
	}
}

func (s *SharingProfileGroupTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharingProfileGroupTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharingProfileGroupTypes(input string) (*SharingProfileGroupTypes, error) {
	vals := map[string]SharingProfileGroupTypes{
		"aadtenants":    SharingProfileGroupTypesAADTenants,
		"subscriptions": SharingProfileGroupTypesSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingProfileGroupTypes(input)
	return &out, nil
}

type SharingState string

const (
	SharingStateFailed     SharingState = "Failed"
	SharingStateInProgress SharingState = "InProgress"
	SharingStateSucceeded  SharingState = "Succeeded"
	SharingStateUnknown    SharingState = "Unknown"
)

func PossibleValuesForSharingState() []string {
	return []string{
		string(SharingStateFailed),
		string(SharingStateInProgress),
		string(SharingStateSucceeded),
		string(SharingStateUnknown),
	}
}

func (s *SharingState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharingState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharingState(input string) (*SharingState, error) {
	vals := map[string]SharingState{
		"failed":     SharingStateFailed,
		"inprogress": SharingStateInProgress,
		"succeeded":  SharingStateSucceeded,
		"unknown":    SharingStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingState(input)
	return &out, nil
}
