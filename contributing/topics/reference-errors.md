# Working with Errors

Following typical Go conventions, error variables within the AzureRM Provider codebase should be named `err`, e.g.

```go
err := someMethodWhichReturnsAnError(...)
```

.. or in the case of a method which returns multiple return types:

```go
model, err := someMethodWhichReturnsAnObjectAndAnError(...)
```

These errors should also be wrapped with more context:

```go
func doSomething() error {
    err := doSomethingWhichCanError()
	if err != nil {
		return fmt.Errorf("performing somethingWhichCanError: %+v", err)
    }
	return nil
}
```

Since this method only returns an error, we can instead reduce this to:

```go
if err := doSomethingWhichCanError(); err != nil {
    return fmt.Errorf("performing somethingWhichCanError: %+v", err)
}
return nil
```

Note that when calling code from within a Terraform Data Source/Resource, the Resource ID type (**note: not** the raw Resource ID) can be used as a formatting argument, for example:

```go
id := someResource.NewResourceGroupID("subscription-id", "my-resource-group")
return fmt.Errorf("deleting %s: %+v", id, err)
```

which will output:

```
deleting Resource Group "my-resource-group" (Subscription ID "subscription-id"): some error"
```

When parsing existing Resource IDs it is sufficient to return the error as is since all the parsing functions return standardised and descriptive error messages:

```go
id, err := someResource.ParseResourceID(state.ID)
if err != nil {
    return err
}
```

# Internal Errors

Internal errors, which are entirely outside the users control (such as failed expectations) that occur within the provider should be prefixed with `internal-error`, for example:

```go
deadline, ok := ctx.Deadline()
if !ok {
    return fmt.Errorf("internal-error: context had no deadline")
}
```

### Notes

Error messages should be both short and clear, using the context as relevant - for example use:

* `return fmt.Errorf("updating %s: %+v", id, err)`
* `return fmt.Errorf("waiting for %s to finish provisioning: %+v", id, err)`
* `return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)`

instead of:

* `return err`
* `return fmt.Errorf("failed updating thing: %+v", err)`
* `return fmt.Errorf("something went wrong: %+v", err)`


This type of error wrapping should be applied to **all** error handling including any nested function that contains two or more error checks (e.g., a function that calls an update API and waits for the update to finish or builds an SDK struct) so practitioners and code maintainers have a clear idea which generated the error.

**NOTE:** Wrapped error messages should generally not start with `failed`, `error`, or an uppercase letter as there will a function higher up the stack that will prefix this.

When returning errors in those situations, it is important to consider the calling context and to exclude any information the calling function is likely to include, while including any additional context then calling function may not have.
