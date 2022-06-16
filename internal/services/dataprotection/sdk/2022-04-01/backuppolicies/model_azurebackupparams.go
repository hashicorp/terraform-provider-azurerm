package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupParameters = AzureBackupParams{}

type AzureBackupParams struct {
	BackupType string `json:"backupType"`

	// Fields inherited from BackupParameters
}

var _ json.Marshaler = AzureBackupParams{}

func (s AzureBackupParams) MarshalJSON() ([]byte, error) {
	type wrapper AzureBackupParams
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureBackupParams: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupParams: %+v", err)
	}
	decoded["objectType"] = "AzureBackupParams"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureBackupParams: %+v", err)
	}

	return encoded, nil
}
