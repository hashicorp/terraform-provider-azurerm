package dataset

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSetKind string

const (
	DataSetKindAdlsGenOneFile       DataSetKind = "AdlsGen1File"
	DataSetKindAdlsGenOneFolder     DataSetKind = "AdlsGen1Folder"
	DataSetKindAdlsGenTwoFile       DataSetKind = "AdlsGen2File"
	DataSetKindAdlsGenTwoFileSystem DataSetKind = "AdlsGen2FileSystem"
	DataSetKindAdlsGenTwoFolder     DataSetKind = "AdlsGen2Folder"
	DataSetKindBlob                 DataSetKind = "Blob"
	DataSetKindBlobFolder           DataSetKind = "BlobFolder"
	DataSetKindContainer            DataSetKind = "Container"
	DataSetKindKustoCluster         DataSetKind = "KustoCluster"
	DataSetKindKustoDatabase        DataSetKind = "KustoDatabase"
	DataSetKindSqlDBTable           DataSetKind = "SqlDBTable"
	DataSetKindSqlDWTable           DataSetKind = "SqlDWTable"
)

func PossibleValuesForDataSetKind() []string {
	return []string{
		string(DataSetKindAdlsGenOneFile),
		string(DataSetKindAdlsGenOneFolder),
		string(DataSetKindAdlsGenTwoFile),
		string(DataSetKindAdlsGenTwoFileSystem),
		string(DataSetKindAdlsGenTwoFolder),
		string(DataSetKindBlob),
		string(DataSetKindBlobFolder),
		string(DataSetKindContainer),
		string(DataSetKindKustoCluster),
		string(DataSetKindKustoDatabase),
		string(DataSetKindSqlDBTable),
		string(DataSetKindSqlDWTable),
	}
}

func (s *DataSetKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataSetKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataSetKind(input string) (*DataSetKind, error) {
	vals := map[string]DataSetKind{
		"adlsgen1file":       DataSetKindAdlsGenOneFile,
		"adlsgen1folder":     DataSetKindAdlsGenOneFolder,
		"adlsgen2file":       DataSetKindAdlsGenTwoFile,
		"adlsgen2filesystem": DataSetKindAdlsGenTwoFileSystem,
		"adlsgen2folder":     DataSetKindAdlsGenTwoFolder,
		"blob":               DataSetKindBlob,
		"blobfolder":         DataSetKindBlobFolder,
		"container":          DataSetKindContainer,
		"kustocluster":       DataSetKindKustoCluster,
		"kustodatabase":      DataSetKindKustoDatabase,
		"sqldbtable":         DataSetKindSqlDBTable,
		"sqldwtable":         DataSetKindSqlDWTable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataSetKind(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
