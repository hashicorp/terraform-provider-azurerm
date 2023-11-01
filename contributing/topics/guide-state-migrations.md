# Guide: State Migrations

State migrations come into play if a resource's implementation needs to change, this can happen for a number a reasons, such as the implementation being incorrect or the API that the resource interacts with changes.

Common scenarios where a state migration would be required in Azure are:
* To correct the format of a Resource ID, the most common example is updating the casing of a segment e.g. `/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/resGroup1` -> `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1`
* Updating the default value of a property in the schema
* Recasting property values in the schema

**Note:** In a lot of cases the changes made by a state migration are not backward compatible, care should be taken when adding state migrations and thorough manual testing should be done. See the section on Testing below.

## Conventions within the AzureRM Provider

State migrations are service specific and are thus kept under a `migration` folder of a service e.g.

├── compute
│   ├── client
│   ├── migration
│   │   ├── managed_disk_v0_to_v1.go
│   ├── managed_disk_resource.go
...

The migration file follows the naming convention of `[resourceName]_[initialVersion]_to_[finalVersion].go` e.g. `managed_disk_v0_to_v1.go`

## Walkthrough for adding a state migration

We will step through an example on how to add a state migration for a made up resource, `capybara_resource.go` in the `animals` service, where one of the Resource ID segments has been cased incorrectly. The state migration will make the following modification:
`/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/Capybaras/capybara1` -> `/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/capybaras/capybara1`

1. Create an empty file under the service's migration folder called `capybara_v0_to_v1.go`

2. The bare minimum required within the file is shown below. Regardless of what the state migration is modifying, `Schema()` and `UpgradeFunc()` must be specified since these are referenced by the resource.
```go
package migration

import (
	"context"
	
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CapybaraV0ToV1 struct{}

func (s CapybaraV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// TODO implement me!
	}
}

func (s CapybaraV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// TODO implement me!
		return nil, nil
	}
}
```

3. Copy over the schema for `capybara_resource.go`. If nothing in the schema is changing then this can be copied over 1:1, however you will want to go through and remove some property attributes that are not required. 
The information in `Schema()` is used to help core serialize/deserialize the values when sending the data over RPC. For this reason the following property attributes should be removed from the schema since they unnecessarily bloat the code base:
   * Default
   * ValidateFunc
   * ForceNew
   * MaxItems
   * MinItems
   * AtLeastOneOf
   * ConflictsWith
   * ExactlyOneOf
   * RequiredWith
4. Fill out the UpgradeFunc to make the modification to the Resource ID. In most cases this involves parsing the old ID insensitively and then overwriting the value for `id` in the state. The file should now look like this:
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
		newId, err := capybaras.ParseCapybaraIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
```

5. Finally we hook the state migration up to the resource. For typed resources this looks like the following
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
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.CapybaraV0ToV1{},
		},
	}
}

// The rest of the resource e.g. Create/Update/Read/Delete methods have been omitted for brevity

```

## Testing

Currently no automated testing for state migrations exist since the testing framework is unable to run different versions of the provider simultaneously. As a result testing for state migrations must be done manually and usually involves the following high level steps:

1. Create the resource using an older version of the provider
2. Locally build a version of the provider containing the state migration
3. Enable development overrides for Terraform
4. Run `terraform plan` using the locally built version of the provider
5. Verify that there are no plan differences