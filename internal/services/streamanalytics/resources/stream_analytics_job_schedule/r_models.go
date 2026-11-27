// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

type JobScheduleResourceModel struct {
	StreamAnalyticsJob string `tfschema:"stream_analytics_job_id"`
	StartMode          string `tfschema:"start_mode"`
	StartTime          string `tfschema:"start_time"`
	LastOutputTime     string `tfschema:"last_output_time"`
}

func (r Resource) ModelObject() any {
	return &JobScheduleResourceModel{}
}
