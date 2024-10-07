package sourcecontrolsyncjobstreams

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlSyncJobStreamOperationPredicate struct {
	Id *string
}

func (p SourceControlSyncJobStreamOperationPredicate) Matches(input SourceControlSyncJobStream) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	return true
}
