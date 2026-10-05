package packageresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PackageProvisioningState string

const (
	PackageProvisioningStateActivitiesStored            PackageProvisioningState = "ActivitiesStored"
	PackageProvisioningStateCanceled                    PackageProvisioningState = "Canceled"
	PackageProvisioningStateConnectionTypeImported      PackageProvisioningState = "ConnectionTypeImported"
	PackageProvisioningStateContentDownloaded           PackageProvisioningState = "ContentDownloaded"
	PackageProvisioningStateContentRetrieved            PackageProvisioningState = "ContentRetrieved"
	PackageProvisioningStateContentStored               PackageProvisioningState = "ContentStored"
	PackageProvisioningStateContentValidated            PackageProvisioningState = "ContentValidated"
	PackageProvisioningStateCreated                     PackageProvisioningState = "Created"
	PackageProvisioningStateCreating                    PackageProvisioningState = "Creating"
	PackageProvisioningStateFailed                      PackageProvisioningState = "Failed"
	PackageProvisioningStateModuleDataStored            PackageProvisioningState = "ModuleDataStored"
	PackageProvisioningStateModuleImportRunbookComplete PackageProvisioningState = "ModuleImportRunbookComplete"
	PackageProvisioningStateRunningImportModuleRunbook  PackageProvisioningState = "RunningImportModuleRunbook"
	PackageProvisioningStateStartingImportModuleRunbook PackageProvisioningState = "StartingImportModuleRunbook"
	PackageProvisioningStateSucceeded                   PackageProvisioningState = "Succeeded"
	PackageProvisioningStateUpdating                    PackageProvisioningState = "Updating"
)

func PossibleValuesForPackageProvisioningState() []string {
	return []string{
		string(PackageProvisioningStateActivitiesStored),
		string(PackageProvisioningStateCanceled),
		string(PackageProvisioningStateConnectionTypeImported),
		string(PackageProvisioningStateContentDownloaded),
		string(PackageProvisioningStateContentRetrieved),
		string(PackageProvisioningStateContentStored),
		string(PackageProvisioningStateContentValidated),
		string(PackageProvisioningStateCreated),
		string(PackageProvisioningStateCreating),
		string(PackageProvisioningStateFailed),
		string(PackageProvisioningStateModuleDataStored),
		string(PackageProvisioningStateModuleImportRunbookComplete),
		string(PackageProvisioningStateRunningImportModuleRunbook),
		string(PackageProvisioningStateStartingImportModuleRunbook),
		string(PackageProvisioningStateSucceeded),
		string(PackageProvisioningStateUpdating),
	}
}

func (s *PackageProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePackageProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePackageProvisioningState(input string) (*PackageProvisioningState, error) {
	vals := map[string]PackageProvisioningState{
		"activitiesstored":            PackageProvisioningStateActivitiesStored,
		"canceled":                    PackageProvisioningStateCanceled,
		"connectiontypeimported":      PackageProvisioningStateConnectionTypeImported,
		"contentdownloaded":           PackageProvisioningStateContentDownloaded,
		"contentretrieved":            PackageProvisioningStateContentRetrieved,
		"contentstored":               PackageProvisioningStateContentStored,
		"contentvalidated":            PackageProvisioningStateContentValidated,
		"created":                     PackageProvisioningStateCreated,
		"creating":                    PackageProvisioningStateCreating,
		"failed":                      PackageProvisioningStateFailed,
		"moduledatastored":            PackageProvisioningStateModuleDataStored,
		"moduleimportrunbookcomplete": PackageProvisioningStateModuleImportRunbookComplete,
		"runningimportmodulerunbook":  PackageProvisioningStateRunningImportModuleRunbook,
		"startingimportmodulerunbook": PackageProvisioningStateStartingImportModuleRunbook,
		"succeeded":                   PackageProvisioningStateSucceeded,
		"updating":                    PackageProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PackageProvisioningState(input)
	return &out, nil
}
