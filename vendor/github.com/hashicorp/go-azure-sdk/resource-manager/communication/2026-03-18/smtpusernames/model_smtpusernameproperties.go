package smtpusernames

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SmtpUsernameProperties struct {
	EntraApplicationId string `json:"entraApplicationId"`
	TenantId           string `json:"tenantId"`
	Username           string `json:"username"`
}
