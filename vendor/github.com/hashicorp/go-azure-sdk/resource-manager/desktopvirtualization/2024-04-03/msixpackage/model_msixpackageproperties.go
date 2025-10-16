package msixpackage

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSIXPackageProperties struct {
	DisplayName           *string                    `json:"displayName,omitempty"`
	ImagePath             *string                    `json:"imagePath,omitempty"`
	IsActive              *bool                      `json:"isActive,omitempty"`
	IsRegularRegistration *bool                      `json:"isRegularRegistration,omitempty"`
	LastUpdated           *string                    `json:"lastUpdated,omitempty"`
	PackageApplications   *[]MsixPackageApplications `json:"packageApplications,omitempty"`
	PackageDependencies   *[]MsixPackageDependencies `json:"packageDependencies,omitempty"`
	PackageFamilyName     *string                    `json:"packageFamilyName,omitempty"`
	PackageName           *string                    `json:"packageName,omitempty"`
	PackageRelativePath   *string                    `json:"packageRelativePath,omitempty"`
	Version               *string                    `json:"version,omitempty"`
}

func (o *MSIXPackageProperties) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *MSIXPackageProperties) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
