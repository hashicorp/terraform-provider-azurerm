package sourcecontrolsyncjobstreams

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlSyncJobStreamById struct {
	Id         *string                                   `json:"id,omitempty"`
	Properties *SourceControlSyncJobStreamByIdProperties `json:"properties,omitempty"`
}
