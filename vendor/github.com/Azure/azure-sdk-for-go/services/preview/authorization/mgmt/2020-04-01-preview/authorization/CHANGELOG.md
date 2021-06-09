# Change History

## Breaking Changes

### Struct Changes

#### Removed Structs

1. CustomErrorResponse

### Signature Changes

#### Funcs

1. RoleAssignmentsClient.Delete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.DeleteByID
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.DeleteByIDPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.DeletePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.Get
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.GetByID
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.GetByIDPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.GetPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.List
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.ListComplete
	- Params
		- From: context.Context, string
		- To: context.Context, string, string
1. RoleAssignmentsClient.ListForResource
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, string
1. RoleAssignmentsClient.ListForResourceComplete
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, string
1. RoleAssignmentsClient.ListForResourceGroup
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListForResourceGroupComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListForResourceGroupPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListForResourcePreparer
	- Params
		- From: context.Context, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, string
1. RoleAssignmentsClient.ListForScope
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListForScopeComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListForScopePreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, string, string
1. RoleAssignmentsClient.ListPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, string, string

## Additive Changes

### New Funcs

1. ErrorAdditionalInfo.MarshalJSON() ([]byte, error)
1. ErrorDetail.MarshalJSON() ([]byte, error)
1. Principal.MarshalJSON() ([]byte, error)
1. RoleAssignmentMetricsResult.MarshalJSON() ([]byte, error)
