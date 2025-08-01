// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
)

func DbSource(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v != string(autonomousdatabases.SourceBackupFromTimestamp) && v != string(autonomousdatabases.SourceCrossRegionDisasterRecovery) {
		return []string{}, append(errors, fmt.Errorf("%v must be %v or %v", k, string(autonomousdatabases.SourceBackupFromTimestamp), string(autonomousdatabases.SourceCrossRegionDisasterRecovery)))
	}

	return []string{}, []error{}
}

func DisasterRecoveryType(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v != string(autonomousdatabases.DisasterRecoveryTypeAdg) && v != string(autonomousdatabases.DisasterRecoveryTypeBackupBased) {
		return []string{}, append(errors, fmt.Errorf("%v must be %v or %v", k, string(autonomousdatabases.DisasterRecoveryTypeAdg), string(autonomousdatabases.DisasterRecoveryTypeBackupBased)))
	}

	return []string{}, []error{}
}

func DatabaseType(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v != string(autonomousdatabases.DataBaseTypeClone) && v != string(autonomousdatabases.DataBaseTypeCloneFromBackupTimestamp) && v != string(autonomousdatabases.DataBaseTypeCrossRegionDisasterRecovery) && v != string(autonomousdatabases.DataBaseTypeRegular) {
		return []string{}, append(errors, fmt.Errorf("%v must be %v or %v", k, string(autonomousdatabases.DataBaseTypeClone), string(autonomousdatabases.DataBaseTypeCloneFromBackupTimestamp), string(autonomousdatabases.DataBaseTypeCrossRegionDisasterRecovery), string(autonomousdatabases.DataBaseTypeRegular)))
	}

	return []string{}, []error{}
}
