package machinelearning

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	synapseValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseSpark() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSparkCreate,
		Read:   resourceSynapseSparkRead,
		Delete: resourceSynapseSparkDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ComputeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{2,16}$`),
					"It can include letters, digits and dashes. It must start with a letter, end with a letter or digit, and be between 2 and 16 characters in length."),
			},

			"machine_learning_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"location": commonschema.Location(),

			"synapse_spark_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: synapseValidate.SparkPoolID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceSynapseSparkCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, _ := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ComputeGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Machine Learning Compute (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_synapse_spark", id.ID())
		}
	}

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var compute machinelearningcomputes.Compute = machinelearningcomputes.SynapseSpark{
		Properties:       nil,
		ComputeLocation:  pointer.To(d.Get("location").(string)),
		Description:      pointer.To(d.Get("description").(string)),
		ResourceId:       pointer.To(d.Get("synapse_spark_pool_id").(string)),
		DisableLocalAuth: pointer.To(!d.Get("local_auth_enabled").(bool)),
	}
	parameters := machinelearningcomputes.ComputeResource{
		Properties: pointer.To(compute),
		Identity:   identity,
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.ComputeCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseSparkRead(d, meta)
}

func resourceSynapseSparkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Machine Learning Compute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Machine Learning Compute (%q): %+v", id, err)
	}

	d.Set("name", id.ComputeName)
	workspaceId := parse.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if spark, ok := (*props).(machinelearningcomputes.SynapseSpark); ok {
				localAuthEnabled := true
				if spark.DisableLocalAuth != nil {
					localAuthEnabled = !*spark.DisableLocalAuth
				}
				d.Set("local_auth_enabled", localAuthEnabled)
				d.Set("description", spark.Description)
				d.Set("synapse_spark_pool_id", spark.ResourceId)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceSynapseSparkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	opts := machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: utils.ToPtr(machinelearningcomputes.UnderlyingResourceActionDetach),
	}

	if err := client.ComputeDeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting Spark %s: %+v", *id, err)
	}

	return nil
}
