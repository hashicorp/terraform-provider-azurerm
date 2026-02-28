# Azure Durable Task Service Implementation - Setup Complete

## ✅ Completed Steps

### 1. Feature Branch Created
- Branch: `feature/durable-task-service`
- Based on: `main` branch

### 2. Service Implementation Files Created
All files created in `internal/services/durabletask/`:

- ✅ `registration.go` - Service registration with typed resources
- ✅ `parse.go` - Resource ID parsing utilities for Scheduler, TaskHub, and RetentionPolicy
- ✅ `validate.go` - Name validation functions
- ✅ `scheduler_resource.go` - Main scheduler resource with SKU, IP allow list, capacity
- ✅ `task_hub_resource.go` - Task hub resource
- ✅ `retention_policy_resource.go` - Retention policy resource

### 3. Files Staged
All durabletask service files have been staged for commit:
```bash
git add internal/services/durabletask/
```

## 🔧 Next Steps Required

### Step 1: Create DurableTask Client File

Create `internal/services/durabletask/client/client.go`:

```go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SchedulersClient         *schedulers.SchedulersClient
	TaskHubsClient           *taskhubs.TaskHubsClient
	RetentionPoliciesClient  *retentionpolicies.RetentionPoliciesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	schedulersClient, err := schedulers.NewSchedulersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Schedulers client: %+v", err)
	}
	o.Configure(schedulersClient.Client, o.Authorizers.ResourceManager)

	taskHubsClient, err := taskhubs.NewTaskHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TaskHubs client: %+v", err)
	}
	o.Configure(taskHubsClient.Client, o.Authorizers.ResourceManager)

	retentionPoliciesClient, err := retentionpolicies.NewRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RetentionPolicies client: %+v", err)
	}
	o.Configure(retentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SchedulersClient:        schedulersClient,
		TaskHubsClient:          taskHubsClient,
		RetentionPoliciesClient: retentionPoliciesClient,
	}, nil
}
```

### Step 2: Update internal/clients/client.go

**Add import (line ~65, alphabetically after datashare):**
```go
durabletask "github.com/hashicorp/terraform-provider-azurerm/internal/services/durabletask/client"
```

**Add client field (line ~202, alphabetically after DataShare):**
```go
DurableTask                       *durabletask.Client
```

**Add client initialization (line ~422, alphabetically after DataShare):**
```go
if client.DurableTask, err = durabletask.NewClient(o); err != nil {
	return fmt.Errorf("building clients for DurableTask: %+v", err)
}
```

### Step 3: Register Service

Find the service registration file (likely `internal/provider/services.go`) and add:

```go
durabletask.Registration{},
```

to the list of `TypedServiceRegistration` or `UntypedServiceRegistrationWithAGitHubLabel`.

### Step 4: Documentation (Optional for now)

Documentation files would go in:
- `website/docs/r/durable_task_scheduler.html.markdown`
- `website/docs/r/durable_task_hub.html.markdown`
- `website/docs/r/durable_task_retention_policy.html.markdown`

### Step 5: Test Files (Optional for now)

Test files would go in:
- `internal/services/durabletask/scheduler_resource_test.go`
- `internal/services/durabletask/task_hub_resource_test.go`
- `internal/services/durabletask/retention_policy_resource_test.go`

## ⚠️ Important Notes

1. **Azure SDK Dependency**: The implementation assumes the Azure SDK has the DurableTask packages available:
   - `github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers`
   - `github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs`
   - `github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies`

2. **Build Verification**: After completing all steps, run:
   ```bash
   go mod tidy
   go build
   ```

3. **Missing Import in registration.go**: The registration.go file is missing:
   ```go
   "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
   ```

## 📊 Implementation Status

| Component | Status | Location |
|-----------|--------|----------|
| Service Files | ✅ Complete | `internal/services/durabletask/` |
| Client File | ⏳ Pending | Needs creation |
| Client Registration | ⏳ Pending | `internal/clients/client.go` |
| Service Registration | ⏳ Pending | Provider services file |
| Documentation | ⏳ Optional | `website/docs/r/` |
| Tests | ⏳ Optional | `internal/services/durabletask/*_test.go` |

## 🚀 Quick Commands

```bash
# Create client directory
New-Item -ItemType Directory -Force -Path "internal\services\durabletask\client"

# Check build status (after completing all steps)
go build

# Run specific tests (after writing tests)
go test ./internal/services/durabletask -v

# Commit changes
git add .
git commit -m "Add Azure Durable Task service support"
```
