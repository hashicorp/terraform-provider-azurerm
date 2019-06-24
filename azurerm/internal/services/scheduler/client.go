package scheduler

import "github.com/Azure/azure-sdk-for-go/services/scheduler/mgmt/2016-03-01/scheduler"

// TODO: remove in 2.0
type Client struct {
	JobCollectionsClient scheduler.JobCollectionsClient //nolint: megacheck
	JobsClient           scheduler.JobsClient           //nolint: megacheck
}
