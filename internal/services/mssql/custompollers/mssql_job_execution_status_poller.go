package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobexecutions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobstepexecutions"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var (
	_ pollers.PollerType = &MssqlJobExecutionStatusPoller{}

	pollingSuccess = pollers.PollResult{
		Status: pollers.PollingStatusSucceeded,
	}

	pollingFailed = pollers.PollResult{
		Status: pollers.PollingStatusFailed,
	}

	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}

	pollingUnknown = pollers.PollResult{
		Status:       pollers.PollingStatusUnknown,
		PollInterval: 10 * time.Second,
	}

	executionLifecycleToResult = map[jobexecutions.JobExecutionLifecycle]pollers.PollResult{
		jobexecutions.JobExecutionLifecycleSkipped:              pollingSuccess,
		jobexecutions.JobExecutionLifecycleSucceeded:            pollingSuccess,
		jobexecutions.JobExecutionLifecycleSucceededWithSkipped: pollingSuccess,

		jobexecutions.JobExecutionLifecycleCreated:                      pollingInProgress,
		jobexecutions.JobExecutionLifecycleInProgress:                   pollingInProgress,
		jobexecutions.JobExecutionLifecycleWaitingForChildJobExecutions: pollingInProgress,
		jobexecutions.JobExecutionLifecycleWaitingForRetry:              pollingInProgress,

		jobexecutions.JobExecutionLifecycleCanceled: pollingFailed,
		jobexecutions.JobExecutionLifecycleFailed:   pollingFailed,
		jobexecutions.JobExecutionLifecycleTimedOut: pollingFailed,
	}
)

type MssqlJobExecutionStatusPoller struct {
	client                  *jobexecutions.JobExecutionsClient
	jobStepExecutionsClient *jobstepexecutions.JobStepExecutionsClient
	id                      jobexecutions.ExecutionId
}

func (p MssqlJobExecutionStatusPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return &pollingFailed, fmt.Errorf("polling %s: %w", p.id, err)
	}

	if resp.Model != nil && resp.Model.Properties != nil {
		if pollingState, ok := executionLifecycleToResult[pointer.From(resp.Model.Properties.Lifecycle)]; ok {
			if pollingState.Status == pollers.PollingStatusFailed {
				lastMessage := "no message returned"
				if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.LastMessage != nil {
					lastMessage = *resp.Model.Properties.LastMessage
				}

				id := jobstepexecutions.NewExecutionID(p.id.SubscriptionId, p.id.ResourceGroupName, p.id.ServerName, p.id.JobAgentName, p.id.JobName, p.id.JobExecutionId)
				jobStepsResp, err := p.jobStepExecutionsClient.ListByJobExecutionComplete(ctx, id, jobstepexecutions.DefaultListByJobExecutionOperationOptions())
				if err != nil {
					return &pollingFailed, pollers.PollingFailedError{
						Message: fmt.Sprintf("execution failed, last job message: `%s`", lastMessage),
					}
				}

				for _, step := range jobStepsResp.Items {
					if props := step.Properties; props != nil && pointer.From(props.ProvisioningState) == jobstepexecutions.ProvisioningStateFailed {
						if props.LastMessage != nil {
							lastMessage = *props.LastMessage
						}
					}
				}

				return &pollingFailed, pollers.PollingFailedError{
					Message: fmt.Sprintf("execution failed, last job step message: `%s`", lastMessage),
				}
			}
			return &pollingState, nil
		}
	}

	return &pollingUnknown, nil
}

func NewMssqlJobExecutionStatusPoller(client *jobexecutions.JobExecutionsClient, jobStepExecutionsClient *jobstepexecutions.JobStepExecutionsClient, id jobexecutions.ExecutionId) *MssqlJobExecutionStatusPoller {
	return &MssqlJobExecutionStatusPoller{
		client:                  client,
		jobStepExecutionsClient: jobStepExecutionsClient,
		id:                      id,
	}
}
