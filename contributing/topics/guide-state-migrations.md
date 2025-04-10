# Guide: State Migrations

State migrations come into play if a resource's implementation needs to change, this can happen for a number a reasons, such as the implementation being incorrect or the API that the resource interacts with changes.

Common scenarios where a state migration would be required in Azure are:
* To correct the format of a Resource ID, the most common example is updating the casing of a segment e.g. `/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/resGroup1` -> `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1`
* Updating the default value of a property in the schema
* Recasting property values in the schema, unlike the scenario's above this also requires changes to the user's config, thus should only be in a major version release

> **Note:** State migrations are one-way by design meaning they're not backward compatible. Once they've been run you can no longer downgrade to an older version of the provider. Care should be taken when adding state migrations and thorough manual testing should be done. See the section on Testing below.

## Conventions within the AzureRM Provider

State migrations are service specific and are thus kept under a `migration` folder of a service e.g.

```
├── compute
│   ├── client
│   ├── migration
│   │   ├── managed_disk_v0_to_v1.go
│   ├── managed_disk_resource.go
...
```


The migration file follows the naming convention of `[resourceName]_[initialVersion]_to_[finalVersion].go` e.g. `managed_disk_v0_to_v1.go`

## Walkthrough for adding a state migration

We will step through an example on how to add a state migration for a made up resource, `capybara_resource.go` in the `animals` service, where one of the Resource ID segments has been cased incorrectly. The state migration will make the following modification:
`/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/Capybaras/capybara1` -> `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/capybaras/capybara1`

1. Create an empty file under the service's migration folder called `capybara_v0_to_v1.go` e.g. (e.g. `./internal/services/animals/migration/capybara_v0_to_v1.go`)

2. The bare minimum required within the file is shown below. Regardless of what the state migration is modifying, `Schema()` and `UpgradeFunc()` must be specified since these are referenced by the resource.
```go
package migration

import (
	"context"
	
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CapybaraV0ToV1 struct{}

func (CapybaraV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// TODO implement me!
	}
}

func (CapybaraV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// TODO implement me!
		return nil, nil
	}
}
```

3. Copy over the schema for `capybara_resource.go`. If nothing in the schema is changing then this can be copied over 1:1, however you will want to go through and remove some property attributes that are not required.
   The `Schema()` is a point-in-time reference to the Terraform Schema for this Resource at this point - and is used by Terraform to deserialize/serialize the object from the Terraform State. For this reason only a subset of attributes should be defined here (including `Type`, `Required`, `Optional`, `Computed` and `Elem` [for maps/lists/sets, including any custom hash functions]) - and the following attributes can be removed from the Schema:
   
   * Default
   * ValidateFunc
   * ForceNew
   * MaxItems
   * MinItems
   * AtLeastOneOf
   * ConflictsWith
   * ExactlyOneOf
   * RequiredWith
   
   Other caveats to look out for when copying the schema over are:
   * in-lining any schema elements which are returned by functions
   * removing any if/else logic within the Schema, in most cases this will be feature flags e.g. `features.FivePointOh()`
   
4. Fill out the UpgradeFunc to update the Terraform State for this resource. Typically this involves parsing the old Resource ID case-insensitively and then setting the correct casing for the `id` field (which is what this example assumes) - however note that State Migrations aren't limited to the `id` field. The file should now look like this:
```go
package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/animals/2023-11-01/capybaras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CapybaraV0ToV1 struct{}

func (s CapybaraV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cuteness": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"pet_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (s CapybaraV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsed, err := capybaras.ParseCapybaraIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsed.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
```

5. Finally, we hook the state migration up to the resource. For typed resources this looks like the following
```go
package animal

import (
	"context"
	"fmt"
	"time"
	
	"github.com/hashicorp/go-azure-sdk/resource-manager/animals/2023-11-01/capybaras"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/animals/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CapybaraResource struct{}

var (
	_ sdk.ResourceWithStateMigration = CapybaraResource{}
)

type CapybaraResourceModel struct {
	Name       string   `tfschema:"name"`
	Cuteness   string   `tfschema:"cuteness"`
	PetNames   []string `tfschema:"pet_names"`

}

func (r CapybaraResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1, // This field references the version which the state migration updates the schema to i.e. v0 -> v1
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.CapybaraV0ToV1{},
		},
	}
}

// The rest of the resource e.g. Create/Update/Read/Delete methods have been omitted for brevity

```

## Testing

Currently, no automated testing for state migrations exist since the testing framework is unable to run different versions of the provider simultaneously. As a result testing for state migrations must be done manually and usually involves the following high level steps:

1. Create the resource using an older version of the provider
2. Locally build a version of the provider containing the state migration
3. Enable development overrides for Terraform
4. Run `terraform plan` and/or `terraform apply` using the locally built version of the provider
5. Verify that there are no plan differences