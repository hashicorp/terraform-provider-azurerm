package staticsites

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteBuildARMResourceProperties struct {
	BuildId                  *string                              `json:"buildId,omitempty"`
	CreatedTimeUtc           *string                              `json:"createdTimeUtc,omitempty"`
	DatabaseConnections      *[]DatabaseConnectionOverview        `json:"databaseConnections,omitempty"`
	Hostname                 *string                              `json:"hostname,omitempty"`
	LastUpdatedOn            *string                              `json:"lastUpdatedOn,omitempty"`
	LinkedBackends           *[]StaticSiteLinkedBackend           `json:"linkedBackends,omitempty"`
	PullRequestTitle         *string                              `json:"pullRequestTitle,omitempty"`
	SourceBranch             *string                              `json:"sourceBranch,omitempty"`
	Status                   *BuildStatus                         `json:"status,omitempty"`
	UserProvidedFunctionApps *[]StaticSiteUserProvidedFunctionApp `json:"userProvidedFunctionApps,omitempty"`
}

func (o *StaticSiteBuildARMResourceProperties) GetCreatedTimeUtcAsTime() (*time.Time, error) {
	if o.CreatedTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *StaticSiteBuildARMResourceProperties) SetCreatedTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTimeUtc = &formatted
}

func (o *StaticSiteBuildARMResourceProperties) GetLastUpdatedOnAsTime() (*time.Time, error) {
	if o.LastUpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *StaticSiteBuildARMResourceProperties) SetLastUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedOn = &formatted
}
