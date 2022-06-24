package azurebackupjob

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupJob struct {
	ActivityID                 string             `json:"activityID"`
	BackupInstanceFriendlyName string             `json:"backupInstanceFriendlyName"`
	BackupInstanceId           *string            `json:"backupInstanceId,omitempty"`
	DataSourceId               string             `json:"dataSourceId"`
	DataSourceLocation         string             `json:"dataSourceLocation"`
	DataSourceName             string             `json:"dataSourceName"`
	DataSourceSetName          *string            `json:"dataSourceSetName,omitempty"`
	DataSourceType             string             `json:"dataSourceType"`
	DestinationDataStoreName   *string            `json:"destinationDataStoreName,omitempty"`
	Duration                   *string            `json:"duration,omitempty"`
	EndTime                    *string            `json:"endTime,omitempty"`
	ErrorDetails               *[]UserFacingError `json:"errorDetails,omitempty"`
	Etag                       *string            `json:"etag,omitempty"`
	ExtendedInfo               *JobExtendedInfo   `json:"extendedInfo,omitempty"`
	IsUserTriggered            bool               `json:"isUserTriggered"`
	Operation                  string             `json:"operation"`
	OperationCategory          string             `json:"operationCategory"`
	PolicyId                   *string            `json:"policyId,omitempty"`
	PolicyName                 *string            `json:"policyName,omitempty"`
	ProgressEnabled            bool               `json:"progressEnabled"`
	ProgressUrl                *string            `json:"progressUrl,omitempty"`
	RestoreType                *string            `json:"restoreType,omitempty"`
	SourceDataStoreName        *string            `json:"sourceDataStoreName,omitempty"`
	SourceResourceGroup        string             `json:"sourceResourceGroup"`
	SourceSubscriptionID       string             `json:"sourceSubscriptionID"`
	StartTime                  string             `json:"startTime"`
	Status                     string             `json:"status"`
	SubscriptionId             string             `json:"subscriptionId"`
	SupportedActions           []string           `json:"supportedActions"`
	VaultName                  string             `json:"vaultName"`
}

func (o *AzureBackupJob) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureBackupJob) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *AzureBackupJob) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AzureBackupJob) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
