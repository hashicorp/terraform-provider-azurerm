package enrollmentaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DepartmentProperties struct {
	CostCenter         *string              `json:"costCenter,omitempty"`
	DepartmentName     *string              `json:"departmentName,omitempty"`
	EnrollmentAccounts *[]EnrollmentAccount `json:"enrollmentAccounts,omitempty"`
	Status             *string              `json:"status,omitempty"`
}
