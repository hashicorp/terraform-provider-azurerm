package dataset

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSet interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataSetImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalDataSetImplementation(input []byte) (DataSet, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataSet into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AdlsGen1File") {
		var out ADLSGen1FileDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ADLSGen1FileDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AdlsGen1Folder") {
		var out ADLSGen1FolderDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ADLSGen1FolderDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AdlsGen2File") {
		var out ADLSGen2FileDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ADLSGen2FileDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AdlsGen2FileSystem") {
		var out ADLSGen2FileSystemDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ADLSGen2FileSystemDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AdlsGen2Folder") {
		var out ADLSGen2FolderDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ADLSGen2FolderDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Container") {
		var out BlobContainerDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobContainerDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Blob") {
		var out BlobDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "BlobFolder") {
		var out BlobFolderDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobFolderDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KustoCluster") {
		var out KustoClusterDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KustoClusterDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KustoDatabase") {
		var out KustoDatabaseDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KustoDatabaseDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDBTable") {
		var out SqlDBTableDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDBTableDataSet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDWTable") {
		var out SqlDWTableDataSet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDWTableDataSet: %+v", err)
		}
		return out, nil
	}

	out := RawDataSetImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
