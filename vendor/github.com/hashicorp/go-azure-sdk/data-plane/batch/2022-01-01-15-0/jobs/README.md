
## `github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01.15.0/jobs` Documentation

The `jobs` SDK allows for interaction with <unknown source data type> `batch` (API Version `2022-01-01.15.0`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01.15.0/jobs"
```


### Client Initialization

```go
client := jobs.NewJobsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `JobsClient.JobAdd`

```go
ctx := context.TODO()

payload := jobs.JobAddParameter{
	// ...
}


read, err := client.JobAdd(ctx, payload, jobs.DefaultJobAddOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobDelete`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

read, err := client.JobDelete(ctx, id, jobs.DefaultJobDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobDisable`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

payload := jobs.JobDisableParameter{
	// ...
}


read, err := client.JobDisable(ctx, id, payload, jobs.DefaultJobDisableOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobEnable`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

read, err := client.JobEnable(ctx, id, jobs.DefaultJobEnableOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobGet`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

read, err := client.JobGet(ctx, id, jobs.DefaultJobGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobGetAllLifetimeStatistics`

```go
ctx := context.TODO()


read, err := client.JobGetAllLifetimeStatistics(ctx, jobs.DefaultJobGetAllLifetimeStatisticsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobGetTaskCounts`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

read, err := client.JobGetTaskCounts(ctx, id, jobs.DefaultJobGetTaskCountsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobList`

```go
ctx := context.TODO()


// alternatively `client.JobList(ctx, jobs.DefaultJobListOperationOptions())` can be used to do batched pagination
items, err := client.JobListComplete(ctx, jobs.DefaultJobListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.JobListFromJobSchedule`

```go
ctx := context.TODO()
id := jobs.NewJobscheduleID("jobScheduleId")

// alternatively `client.JobListFromJobSchedule(ctx, id, jobs.DefaultJobListFromJobScheduleOperationOptions())` can be used to do batched pagination
items, err := client.JobListFromJobScheduleComplete(ctx, id, jobs.DefaultJobListFromJobScheduleOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.JobListPreparationAndReleaseTaskStatus`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

// alternatively `client.JobListPreparationAndReleaseTaskStatus(ctx, id, jobs.DefaultJobListPreparationAndReleaseTaskStatusOperationOptions())` can be used to do batched pagination
items, err := client.JobListPreparationAndReleaseTaskStatusComplete(ctx, id, jobs.DefaultJobListPreparationAndReleaseTaskStatusOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `JobsClient.JobPatch`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

payload := jobs.JobPatchParameter{
	// ...
}


read, err := client.JobPatch(ctx, id, payload, jobs.DefaultJobPatchOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobTerminate`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

payload := jobs.JobTerminateParameter{
	// ...
}


read, err := client.JobTerminate(ctx, id, payload, jobs.DefaultJobTerminateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `JobsClient.JobUpdate`

```go
ctx := context.TODO()
id := jobs.NewJobID("jobId")

payload := jobs.JobUpdateParameter{
	// ...
}


read, err := client.JobUpdate(ctx, id, payload, jobs.DefaultJobUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
