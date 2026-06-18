package integrationruntime

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SelfHostedIntegrationRuntimeStatusTypeProperties struct {
	AutoUpdate                             *IntegrationRuntimeAutoUpdate                    `json:"autoUpdate,omitempty"`
	AutoUpdateETA                          *string                                          `json:"autoUpdateETA,omitempty"`
	Capabilities                           *map[string]string                               `json:"capabilities,omitempty"`
	CreateTime                             *string                                          `json:"createTime,omitempty"`
	InternalChannelEncryption              *IntegrationRuntimeInternalChannelEncryptionMode `json:"internalChannelEncryption,omitempty"`
	LatestVersion                          *string                                          `json:"latestVersion,omitempty"`
	Links                                  *[]LinkedIntegrationRuntime                      `json:"links,omitempty"`
	LocalTimeZoneOffset                    *string                                          `json:"localTimeZoneOffset,omitempty"`
	NodeCommunicationChannelEncryptionMode *string                                          `json:"nodeCommunicationChannelEncryptionMode,omitempty"`
	Nodes                                  *[]SelfHostedIntegrationRuntimeNode              `json:"nodes,omitempty"`
	PushedVersion                          *string                                          `json:"pushedVersion,omitempty"`
	ScheduledUpdateDate                    *string                                          `json:"scheduledUpdateDate,omitempty"`
	ServiceURLs                            *[]string                                        `json:"serviceUrls,omitempty"`
	TaskQueueId                            *string                                          `json:"taskQueueId,omitempty"`
	UpdateDelayOffset                      *string                                          `json:"updateDelayOffset,omitempty"`
	Version                                *string                                          `json:"version,omitempty"`
	VersionStatus                          *string                                          `json:"versionStatus,omitempty"`
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) GetAutoUpdateETAAsTime() (*time.Time, error) {
	if o.AutoUpdateETA == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AutoUpdateETA, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) SetAutoUpdateETAAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AutoUpdateETA = &formatted
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) GetScheduledUpdateDateAsTime() (*time.Time, error) {
	if o.ScheduledUpdateDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduledUpdateDate, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeStatusTypeProperties) SetScheduledUpdateDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduledUpdateDate = &formatted
}
