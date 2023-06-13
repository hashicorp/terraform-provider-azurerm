package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogzOrganizationProperties struct {
	CompanyName     *string `json:"companyName,omitempty"`
	EnterpriseAppId *string `json:"enterpriseAppId,omitempty"`
	Id              *string `json:"id,omitempty"`
	SingleSignOnUrl *string `json:"singleSignOnUrl,omitempty"`
}
