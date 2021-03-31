# Change History

### New Funcs

1. NewRoleAssignmentMetricsClient(string) RoleAssignmentMetricsClient
1. NewRoleAssignmentMetricsClientWithBaseURI(string, string) RoleAssignmentMetricsClient
1. RoleAssignmentMetricsClient.GetMetricsForSubscription(context.Context) (RoleAssignmentMetricsResult, error)
1. RoleAssignmentMetricsClient.GetMetricsForSubscriptionPreparer(context.Context) (*http.Request, error)
1. RoleAssignmentMetricsClient.GetMetricsForSubscriptionResponder(*http.Response) (RoleAssignmentMetricsResult, error)
1. RoleAssignmentMetricsClient.GetMetricsForSubscriptionSender(*http.Request) (*http.Response, error)

## Struct Changes

### New Structs

1. CustomErrorResponse
1. RoleAssignmentMetricsClient
1. RoleAssignmentMetricsResult

### New Struct Fields

1. RoleAssignmentProperties.DelegatedManagedIdentityResourceID
