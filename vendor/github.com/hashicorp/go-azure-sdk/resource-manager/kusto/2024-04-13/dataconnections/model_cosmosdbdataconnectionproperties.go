package dataconnections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDbDataConnectionProperties struct {
	CosmosDbAccountResourceId string             `json:"cosmosDbAccountResourceId"`
	CosmosDbContainer         string             `json:"cosmosDbContainer"`
	CosmosDbDatabase          string             `json:"cosmosDbDatabase"`
	ManagedIdentityObjectId   *string            `json:"managedIdentityObjectId,omitempty"`
	ManagedIdentityResourceId string             `json:"managedIdentityResourceId"`
	MappingRuleName           *string            `json:"mappingRuleName,omitempty"`
	ProvisioningState         *ProvisioningState `json:"provisioningState,omitempty"`
	RetrievalStartDate        *string            `json:"retrievalStartDate,omitempty"`
	TableName                 string             `json:"tableName"`
}

func (o *CosmosDbDataConnectionProperties) GetRetrievalStartDateAsTime() (*time.Time, error) {
	if o.RetrievalStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RetrievalStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *CosmosDbDataConnectionProperties) SetRetrievalStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RetrievalStartDate = &formatted
}
