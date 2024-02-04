package enrollmentaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnrollmentAccountListResult struct {
	NextLink *string              `json:"nextLink,omitempty"`
	Value    *[]EnrollmentAccount `json:"value,omitempty"`
}
