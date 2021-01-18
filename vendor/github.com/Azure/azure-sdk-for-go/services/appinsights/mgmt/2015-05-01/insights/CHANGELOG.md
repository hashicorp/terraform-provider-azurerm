Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewApplicationInsightsComponentListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationInsightsComponentListResult) (ApplicationInsightsComponentListResult, error))` to `(ApplicationInsightsComponentListResult, func(context.Context, ApplicationInsightsComponentListResult) (ApplicationInsightsComponentListResult, error))`
- Function `NewWebTestListResultPage` parameter(s) have been changed from `(func(context.Context, WebTestListResult) (WebTestListResult, error))` to `(WebTestListResult, func(context.Context, WebTestListResult) (WebTestListResult, error))`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`

## New Content

- New function `MyWorkbooksClient.DeleteSender(*http.Request) (*http.Response, error)`
- New function `MyWorkbooksClient.ListByResourceGroupPreparer(context.Context, string, CategoryType, []string, *bool) (*http.Request, error)`
- New function `MyWorkbooksClient.ListBySubscription(context.Context, CategoryType, []string, *bool) (MyWorkbooksListResult, error)`
- New function `MyWorkbooksClient.Get(context.Context, string, string) (MyWorkbook, error)`
- New function `MyWorkbooksClient.GetSender(*http.Request) (*http.Response, error)`
- New function `NewMyWorkbooksClientWithBaseURI(string, string) MyWorkbooksClient`
- New function `MyWorkbooksClient.ListByResourceGroupSender(*http.Request) (*http.Response, error)`
- New function `MyWorkbooksClient.CreateOrUpdateResponder(*http.Response) (MyWorkbook, error)`
- New function `NewMyWorkbooksClient(string) MyWorkbooksClient`
- New function `*MyWorkbook.UnmarshalJSON([]byte) error`
- New function `MyWorkbooksClient.GetPreparer(context.Context, string, string) (*http.Request, error)`
- New function `MyWorkbooksClient.ListBySubscriptionPreparer(context.Context, CategoryType, []string, *bool) (*http.Request, error)`
- New function `MyWorkbooksClient.Delete(context.Context, string, string) (autorest.Response, error)`
- New function `MyWorkbooksClient.ListBySubscriptionResponder(*http.Response) (MyWorkbooksListResult, error)`
- New function `MyWorkbooksClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)`
- New function `MyWorkbooksClient.CreateOrUpdatePreparer(context.Context, string, string, MyWorkbook) (*http.Request, error)`
- New function `MyWorkbookProperties.MarshalJSON() ([]byte, error)`
- New function `MyWorkbooksClient.Update(context.Context, string, string, MyWorkbook) (MyWorkbook, error)`
- New function `MyWorkbooksClient.UpdatePreparer(context.Context, string, string, MyWorkbook) (*http.Request, error)`
- New function `MyWorkbooksClient.ListByResourceGroupResponder(*http.Response) (MyWorkbooksListResult, error)`
- New function `MyWorkbookResource.MarshalJSON() ([]byte, error)`
- New function `MyWorkbooksClient.ListByResourceGroup(context.Context, string, CategoryType, []string, *bool) (MyWorkbooksListResult, error)`
- New function `MyWorkbooksClient.DeleteResponder(*http.Response) (autorest.Response, error)`
- New function `MyWorkbooksClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)`
- New function `MyWorkbook.MarshalJSON() ([]byte, error)`
- New function `MyWorkbooksClient.CreateOrUpdate(context.Context, string, string, MyWorkbook) (MyWorkbook, error)`
- New function `MyWorkbooksClient.DeletePreparer(context.Context, string, string) (*http.Request, error)`
- New function `MyWorkbooksClient.UpdateResponder(*http.Response) (MyWorkbook, error)`
- New function `MyWorkbooksClient.UpdateSender(*http.Request) (*http.Response, error)`
- New function `MyWorkbooksClient.GetResponder(*http.Response) (MyWorkbook, error)`
- New struct `MyWorkbook`
- New struct `MyWorkbookError`
- New struct `MyWorkbookProperties`
- New struct `MyWorkbookResource`
- New struct `MyWorkbooksClient`
- New struct `MyWorkbooksListResult`
