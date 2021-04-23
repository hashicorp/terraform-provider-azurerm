package pluginsdk

/*
This is tracking the changes necessary for Plugin SDKv2

1. (DONE) Base Layer
2. (DONE) Package validation
	All validation functions need to be moved into `validation` - since otherwise the test namespace needs fixing
3. (DONE) Updating the Migration functions to use the new wrapper
4. (DONE) Updating the Customize Diff functions to use the new wrapper
5. (DONE) Updating the Import functions to use the new wrappers

6. Update the `validate` package tests to include `_test` suffix for the package name
	At the same time, add the package name to the generator to reflect this:
	./azurerm/internal/tools/generator-resource-id/main.go:15

7. Gradually updating each service package to use the new wrappers
	This script _should_ be safe to re-run multiple times over the same service package
	It's only Sed, but it does most of the job.

8. Upgrade to Plugin SDKv2
	At this point there should be ~100 imports for "github.com/hashicorp/terraform-plugin-sdk", meaning
	this should be a fairly small change - at which point we can upgrade to v2.
	The following bits are notable:
		- github.com/hashicorp/terraform-plugin-sdk => github.com/hashicorp/terraform-plugin-sdk/v2
		- ./azurerm/internal/provider/provider.go:323-328 needs to become:
			```
			client, err := clients.Build(context.TODO(), clientBuilder)
			if err != nil {
				return nil, err
			}

			client.StopContext = context.TODO()
			```
			Whilst this does break cmd+c temporarily for methods using the StopContext, it means
			this otherwise functions which seems like a reasonable short-term trade-off.
		- Update the wrappers to use the real contexts
		- Update the Typed SDK
		- Remove the vendored AzureAD provider / fix testing

9. Update the CRUD|Import methods from `Create` to `CreateContext`
	This means in any resource without Locks we can remove:
		```
		ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
		defer cancel()
		```
	Since the `ctx` passed into the {name}Context method already has a timer, any resource using
	a lock will need to use the other "without timeout" method, but that's not that many resources.

10. Done? We could probably remove the wrappers, I guess.
*/
