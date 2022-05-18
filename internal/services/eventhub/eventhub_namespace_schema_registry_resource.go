package eventhub

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/sdk/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/sdk/2021-11-01/schemaregistry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"log"
	"time"
)

func resourceEventHubNamespaceSchemaRegistry() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceSchemaRegistryCreateUpdate,
		Read:   resourceEventHubNamespaceSchemaRegistryRead,
		Update: resourceEventHubNamespaceSchemaRegistryCreateUpdate,
		Delete: resourceEventHubNamespaceSchemaRegistryDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schemaregistry.ParseSchemagroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateSchemaGroupName(),
			},

			"namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: namespaces.ValidateNamespaceID,
			},

			"schema_compatibility_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(schemaregistry.SchemaCompatibilityNone),
					string(schemaregistry.SchemaCompatibilityBackward),
					string(schemaregistry.SchemaCompatibilityForward),
				}, false),
			},

			"schema_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(schemaregistry.SchemaTypeUnknown),
					string(schemaregistry.SchemaTypeAvro),
				}, false),
			},

			"schema_group_property": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceEventHubNamespaceSchemaRegistryCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespaceSchemaGroupClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace Schema Registry creation.")

	namespaceId, err := namespaces.ParseNamespaceID(d.Get("namespace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing eventhub namespace %s error: %+v", namespaceId.ID(), err)
	}

	id := schemaregistry.NewSchemagroupID(subscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_schema_group", id.ID())
		}
	}

	schemaCompatibilityType := schemaregistry.SchemaCompatibility(d.Get("schema_compatibility_type").(string))
	schemaType := schemaregistry.SchemaType(d.Get("schema_type").(string))

	parameters := schemaregistry.SchemaGroup{
		Properties: &schemaregistry.SchemaGroupProperties{
			SchemaCompatibility: &schemaCompatibilityType,
			SchemaType:          &schemaType,
		},
	}

	if value, ok := d.GetOk("schema_group_property"); ok {
		result := make(map[string]string)
		for k, v := range value.(map[string]interface{}) {
			result[k] = v.(string)
			parameters.Properties.GroupProperties = &result
		}
	}
	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceEventHubNamespaceSchemaRegistryRead(d, meta)
}

func resourceEventHubNamespaceSchemaRegistryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespaceSchemaGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schemaregistry.ParseSchemagroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.SchemaGroupName)

	namespaceId := namespaces.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	d.Set("namespace_id", namespaceId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.SchemaCompatibility != nil {
				d.Set("schema_compatibility_type", string(*props.SchemaCompatibility))
			}
			if props.SchemaType != nil {
				d.Set("schema_type", string(*props.SchemaType))
			}
			if props.GroupProperties != nil {
				properties := make(map[string]string)
				for k, v := range *props.GroupProperties {
					properties[k] = v
				}
				d.Set("schema_group_property", properties)
			}
		}
	}

	return nil
}

func resourceEventHubNamespaceSchemaRegistryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespaceSchemaGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schemaregistry.ParseSchemagroupIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
