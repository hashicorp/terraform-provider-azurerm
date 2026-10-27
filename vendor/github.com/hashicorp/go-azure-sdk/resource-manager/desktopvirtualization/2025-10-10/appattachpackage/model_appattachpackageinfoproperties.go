package appattachpackage

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppAttachPackageInfoProperties struct {
	CertificateExpiry     *string                    `json:"certificateExpiry,omitempty"`
	CertificateName       *string                    `json:"certificateName,omitempty"`
	DisplayName           *string                    `json:"displayName,omitempty"`
	ImagePath             *string                    `json:"imagePath,omitempty"`
	IsActive              *bool                      `json:"isActive,omitempty"`
	IsPackageTimestamped  *PackageTimestamped        `json:"isPackageTimestamped,omitempty"`
	IsRegularRegistration *bool                      `json:"isRegularRegistration,omitempty"`
	LastUpdated           *string                    `json:"lastUpdated,omitempty"`
	PackageAlias          *string                    `json:"packageAlias,omitempty"`
	PackageApplications   *[]MsixPackageApplications `json:"packageApplications,omitempty"`
	PackageDependencies   *[]MsixPackageDependencies `json:"packageDependencies,omitempty"`
	PackageFamilyName     *string                    `json:"packageFamilyName,omitempty"`
	PackageFullName       *string                    `json:"packageFullName,omitempty"`
	PackageName           *string                    `json:"packageName,omitempty"`
	PackageRelativePath   *string                    `json:"packageRelativePath,omitempty"`
	Version               *string                    `json:"version,omitempty"`
}

func (o *AppAttachPackageInfoProperties) GetCertificateExpiryAsTime() (*time.Time, error) {
	if o.CertificateExpiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CertificateExpiry, "2006-01-02T15:04:05Z07:00")
}

func (o *AppAttachPackageInfoProperties) SetCertificateExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CertificateExpiry = &formatted
}

func (o *AppAttachPackageInfoProperties) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *AppAttachPackageInfoProperties) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
