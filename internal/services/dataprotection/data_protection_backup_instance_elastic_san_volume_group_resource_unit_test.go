// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
)

func TestElasticSanVolumeGroupBackupDatasourceParameters(t *testing.T) {
	expected := "volume1"
	expanded := expandElasticSanVolumeGroupBackupDatasourceParameters(expected)
	if expanded == nil || len(*expanded) != 1 {
		t.Fatalf("expected one backup datasource parameter, got %#v", expanded)
	}

	encoded, err := json.Marshal((*expanded)[0])
	if err != nil {
		t.Fatalf("marshaling backup datasource parameters: %+v", err)
	}
	decoded, err := backupinstanceresources.UnmarshalBackupDatasourceParametersImplementation(encoded)
	if err != nil {
		t.Fatalf("unmarshaling backup datasource parameters: %+v", err)
	}

	actual, err := flattenElasticSanVolumeGroupBackupDatasourceParameters(&[]backupinstanceresources.BackupDatasourceParameters{decoded})
	if err != nil {
		t.Fatalf("flattening backup datasource parameters: %+v", err)
	}
	if actual != expected {
		t.Fatalf("expected volume selector %q, got %q", expected, actual)
	}
}
