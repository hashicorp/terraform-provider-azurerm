package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationDetails struct {
	AdminDetails *[]AdministratorDetails `json:"admin_details,omitempty"`
	Id           *string                 `json:"id,omitempty"`
}
