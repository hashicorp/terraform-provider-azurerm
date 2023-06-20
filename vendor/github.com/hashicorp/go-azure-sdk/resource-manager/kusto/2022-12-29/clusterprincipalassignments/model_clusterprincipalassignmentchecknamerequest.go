package clusterprincipalassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPrincipalAssignmentCheckNameRequest struct {
	Name string                  `json:"name"`
	Type PrincipalAssignmentType `json:"type"`
}
