# Change History

## Additive Changes

### New Funcs

1. *Variable.UnmarshalJSON([]byte) error
1. *VariableListResultIterator.Next() error
1. *VariableListResultIterator.NextWithContext(context.Context) error
1. *VariableListResultPage.Next() error
1. *VariableListResultPage.NextWithContext(context.Context) error
1. *VariableValue.UnmarshalJSON([]byte) error
1. *VariableValueListResultIterator.Next() error
1. *VariableValueListResultIterator.NextWithContext(context.Context) error
1. *VariableValueListResultPage.Next() error
1. *VariableValueListResultPage.NextWithContext(context.Context) error
1. NewVariableListResultIterator(VariableListResultPage) VariableListResultIterator
1. NewVariableListResultPage(VariableListResult, func(context.Context, VariableListResult) (VariableListResult, error)) VariableListResultPage
1. NewVariableValueListResultIterator(VariableValueListResultPage) VariableValueListResultIterator
1. NewVariableValueListResultPage(VariableValueListResult, func(context.Context, VariableValueListResult) (VariableValueListResult, error)) VariableValueListResultPage
1. NewVariableValuesClient(string) VariableValuesClient
1. NewVariableValuesClientWithBaseURI(string, string) VariableValuesClient
1. NewVariablesClient(string) VariablesClient
1. NewVariablesClientWithBaseURI(string, string) VariablesClient
1. Variable.MarshalJSON() ([]byte, error)
1. VariableListResult.IsEmpty() bool
1. VariableListResult.MarshalJSON() ([]byte, error)
1. VariableListResultIterator.NotDone() bool
1. VariableListResultIterator.Response() VariableListResult
1. VariableListResultIterator.Value() Variable
1. VariableListResultPage.NotDone() bool
1. VariableListResultPage.Response() VariableListResult
1. VariableListResultPage.Values() []Variable
1. VariableValue.MarshalJSON() ([]byte, error)
1. VariableValueListResult.IsEmpty() bool
1. VariableValueListResult.MarshalJSON() ([]byte, error)
1. VariableValueListResultIterator.NotDone() bool
1. VariableValueListResultIterator.Response() VariableValueListResult
1. VariableValueListResultIterator.Value() VariableValue
1. VariableValueListResultPage.NotDone() bool
1. VariableValueListResultPage.Response() VariableValueListResult
1. VariableValueListResultPage.Values() []VariableValue
1. VariableValuesClient.CreateOrUpdate(context.Context, string, string, VariableValue) (VariableValue, error)
1. VariableValuesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, string, VariableValue) (VariableValue, error)
1. VariableValuesClient.CreateOrUpdateAtManagementGroupPreparer(context.Context, string, string, string, VariableValue) (*http.Request, error)
1. VariableValuesClient.CreateOrUpdateAtManagementGroupResponder(*http.Response) (VariableValue, error)
1. VariableValuesClient.CreateOrUpdateAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.CreateOrUpdatePreparer(context.Context, string, string, VariableValue) (*http.Request, error)
1. VariableValuesClient.CreateOrUpdateResponder(*http.Response) (VariableValue, error)
1. VariableValuesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.Delete(context.Context, string, string) (autorest.Response, error)
1. VariableValuesClient.DeleteAtManagementGroup(context.Context, string, string, string) (autorest.Response, error)
1. VariableValuesClient.DeleteAtManagementGroupPreparer(context.Context, string, string, string) (*http.Request, error)
1. VariableValuesClient.DeleteAtManagementGroupResponder(*http.Response) (autorest.Response, error)
1. VariableValuesClient.DeleteAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.DeletePreparer(context.Context, string, string) (*http.Request, error)
1. VariableValuesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. VariableValuesClient.DeleteSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.Get(context.Context, string, string) (VariableValue, error)
1. VariableValuesClient.GetAtManagementGroup(context.Context, string, string, string) (VariableValue, error)
1. VariableValuesClient.GetAtManagementGroupPreparer(context.Context, string, string, string) (*http.Request, error)
1. VariableValuesClient.GetAtManagementGroupResponder(*http.Response) (VariableValue, error)
1. VariableValuesClient.GetAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. VariableValuesClient.GetResponder(*http.Response) (VariableValue, error)
1. VariableValuesClient.GetSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.List(context.Context, string) (VariableValueListResultPage, error)
1. VariableValuesClient.ListComplete(context.Context, string) (VariableValueListResultIterator, error)
1. VariableValuesClient.ListForManagementGroup(context.Context, string, string) (VariableValueListResultPage, error)
1. VariableValuesClient.ListForManagementGroupComplete(context.Context, string, string) (VariableValueListResultIterator, error)
1. VariableValuesClient.ListForManagementGroupPreparer(context.Context, string, string) (*http.Request, error)
1. VariableValuesClient.ListForManagementGroupResponder(*http.Response) (VariableValueListResult, error)
1. VariableValuesClient.ListForManagementGroupSender(*http.Request) (*http.Response, error)
1. VariableValuesClient.ListPreparer(context.Context, string) (*http.Request, error)
1. VariableValuesClient.ListResponder(*http.Response) (VariableValueListResult, error)
1. VariableValuesClient.ListSender(*http.Request) (*http.Response, error)
1. VariablesClient.CreateOrUpdate(context.Context, string, Variable) (Variable, error)
1. VariablesClient.CreateOrUpdateAtManagementGroup(context.Context, string, string, Variable) (Variable, error)
1. VariablesClient.CreateOrUpdateAtManagementGroupPreparer(context.Context, string, string, Variable) (*http.Request, error)
1. VariablesClient.CreateOrUpdateAtManagementGroupResponder(*http.Response) (Variable, error)
1. VariablesClient.CreateOrUpdateAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariablesClient.CreateOrUpdatePreparer(context.Context, string, Variable) (*http.Request, error)
1. VariablesClient.CreateOrUpdateResponder(*http.Response) (Variable, error)
1. VariablesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. VariablesClient.Delete(context.Context, string) (autorest.Response, error)
1. VariablesClient.DeleteAtManagementGroup(context.Context, string, string) (autorest.Response, error)
1. VariablesClient.DeleteAtManagementGroupPreparer(context.Context, string, string) (*http.Request, error)
1. VariablesClient.DeleteAtManagementGroupResponder(*http.Response) (autorest.Response, error)
1. VariablesClient.DeleteAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariablesClient.DeletePreparer(context.Context, string) (*http.Request, error)
1. VariablesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. VariablesClient.DeleteSender(*http.Request) (*http.Response, error)
1. VariablesClient.Get(context.Context, string) (Variable, error)
1. VariablesClient.GetAtManagementGroup(context.Context, string, string) (Variable, error)
1. VariablesClient.GetAtManagementGroupPreparer(context.Context, string, string) (*http.Request, error)
1. VariablesClient.GetAtManagementGroupResponder(*http.Response) (Variable, error)
1. VariablesClient.GetAtManagementGroupSender(*http.Request) (*http.Response, error)
1. VariablesClient.GetPreparer(context.Context, string) (*http.Request, error)
1. VariablesClient.GetResponder(*http.Response) (Variable, error)
1. VariablesClient.GetSender(*http.Request) (*http.Response, error)
1. VariablesClient.List(context.Context) (VariableListResultPage, error)
1. VariablesClient.ListComplete(context.Context) (VariableListResultIterator, error)
1. VariablesClient.ListForManagementGroup(context.Context, string) (VariableListResultPage, error)
1. VariablesClient.ListForManagementGroupComplete(context.Context, string) (VariableListResultIterator, error)
1. VariablesClient.ListForManagementGroupPreparer(context.Context, string) (*http.Request, error)
1. VariablesClient.ListForManagementGroupResponder(*http.Response) (VariableListResult, error)
1. VariablesClient.ListForManagementGroupSender(*http.Request) (*http.Response, error)
1. VariablesClient.ListPreparer(context.Context) (*http.Request, error)
1. VariablesClient.ListResponder(*http.Response) (VariableListResult, error)
1. VariablesClient.ListSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. Variable
1. VariableColumn
1. VariableListResult
1. VariableListResultIterator
1. VariableListResultPage
1. VariableProperties
1. VariableValue
1. VariableValueColumnValue
1. VariableValueListResult
1. VariableValueListResultIterator
1. VariableValueListResultPage
1. VariableValueProperties
1. VariableValuesClient
1. VariablesClient
