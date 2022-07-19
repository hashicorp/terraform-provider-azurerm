package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CompanyInfo struct {
	Business        *string `json:"business,omitempty"`
	Country         *string `json:"country,omitempty"`
	Domain          *string `json:"domain,omitempty"`
	EmployeesNumber *string `json:"employeesNumber,omitempty"`
	State           *string `json:"state,omitempty"`
}
