// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArcKubernetesCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArcKubernetesClusterCreate,
		Read:   resourceArcKubernetesClusterRead,
		Update: resourceArcKubernetesClusterUpdate,
		Delete: resourceArcKubernetesClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := arckubernetes.ParseConnectedClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
					"The name of Arc Kubernetes Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"agent_public_key_certificate": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.Base64EncodedString,
				// agentPublicKeyCertificate input must be empty for Connected Cluster of Kind: Provisioned Cluster
				ConflictsWith: []string{"kind"},
			},

			"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

			"location": commonschema.Location(),

			"aad_profile": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				RequiredWith: []string{"kind"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"admin_group_object_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},

						"azure_rbac_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"arc_agent_desired_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"arc_agent_auto_upgrade_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"azure_hybrid_benefit": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(arckubernetes.AzureHybridBenefitNotApplicable),
				ValidateFunc: validation.StringInSlice([]string{
					string(arckubernetes.AzureHybridBenefitTrue),
					string(arckubernetes.AzureHybridBenefitFalse),
					string(arckubernetes.AzureHybridBenefitNotApplicable),
				}, false),
			},

			"kind": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"agent_public_key_certificate"},
				ValidateFunc: validation.StringInSlice([]string{
					string(arckubernetes.ConnectedClusterKindProvisionedCluster),
				}, false),
			},

			"agent_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"distribution": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"infrastructure": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kubernetes_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"offering": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"total_core_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"total_node_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArcKubernetesClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := arckubernetes.NewConnectedClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.ConnectedClusterGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_arc_kubernetes_cluster", id.ID())
	}

	identityValue, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	autoUpgradeOption := arckubernetes.AutoUpgradeOptionsEnabled
	if !d.Get("arc_agent_auto_upgrade_enabled").(bool) {
		autoUpgradeOption = arckubernetes.AutoUpgradeOptionsDisabled
	}

	location := location.Normalize(d.Get("location").(string))
	props := arckubernetes.ConnectedCluster{
		Identity: *identityValue,
		Location: location,
		Properties: arckubernetes.ConnectedClusterProperties{
			AgentPublicKeyCertificate: d.Get("agent_public_key_certificate").(string),
			ArcAgentProfile: &arckubernetes.ArcAgentProfile{
				AgentAutoUpgrade: pointer.To(autoUpgradeOption),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if kindVal := d.Get("kind").(string); kindVal != "" {
		props.Kind = pointer.To(arckubernetes.ConnectedClusterKind(kindVal))
	}

	if hybridBenefitVal := d.Get("azure_hybrid_benefit").(string); hybridBenefitVal != "" {
		props.Properties.AzureHybridBenefit = pointer.To(arckubernetes.AzureHybridBenefit(hybridBenefitVal))
	}

	if aadProfileVal := d.Get("aad_profile").([]interface{}); len(aadProfileVal) != 0 {
		props.Properties.AadProfile = expandArcKubernetesClusterAadProfile(aadProfileVal)
	}

	if desiredVersion := d.Get("arc_agent_desired_version").(string); desiredVersion != "" {
		props.Properties.ArcAgentProfile.DesiredAgentVersion = pointer.To(desiredVersion)
	}

	if err := client.ConnectedClusterCreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArcKubernetesClusterRead(d, meta)
}

func resourceArcKubernetesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectedClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if err := d.Set("identity", identity.FlattenSystemAssigned(&model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("kind", string(pointer.From(model.Kind)))
		d.Set("location", location.Normalize(model.Location))
		props := model.Properties
		d.Set("aad_profile", flattenArcKubernetesClusterAadProfile(props.AadProfile))
		d.Set("azure_hybrid_benefit", string(pointer.From(props.AzureHybridBenefit)))
		d.Set("agent_public_key_certificate", props.AgentPublicKeyCertificate)
		d.Set("agent_version", props.AgentVersion)
		d.Set("distribution", props.Distribution)
		d.Set("infrastructure", props.Infrastructure)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("offering", props.Offering)
		d.Set("total_core_count", props.TotalCoreCount)
		d.Set("total_node_count", props.TotalNodeCount)

		arcAgentAutoUpgradeEnabled := true
		arcAgentdesiredVersion := ""
		if arcAgentProfile := props.ArcAgentProfile; arcAgentProfile != nil {
			arcAgentdesiredVersion = pointer.From(arcAgentProfile.DesiredAgentVersion)
			if arcAgentProfile.AgentAutoUpgrade != nil && *arcAgentProfile.AgentAutoUpgrade == arckubernetes.AutoUpgradeOptionsDisabled {
				arcAgentAutoUpgradeEnabled = false
			}
		}
		d.Set("arc_agent_auto_upgrade_enabled", arcAgentAutoUpgradeEnabled)
		d.Set("arc_agent_desired_version", arcAgentdesiredVersion)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceArcKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := resp.Model

	if d.HasChange("aad_profile") {
		payload.Properties.AadProfile = expandArcKubernetesClusterAadProfile(d.Get("aad_profile").([]interface{}))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("azure_hybrid_benefit") {
		payload.Properties.AzureHybridBenefit = pointer.To(arckubernetes.AzureHybridBenefit(d.Get("azure_hybrid_benefit").(string)))
	}

	if d.HasChange("arc_agent_desired_version") {
		if desiredVersion := d.Get("arc_agent_desired_version").(string); desiredVersion != "" {
			payload.Properties.ArcAgentProfile.DesiredAgentVersion = pointer.To(desiredVersion)
		} else {
			payload.Properties.ArcAgentProfile.DesiredAgentVersion = nil
		}
	}

	if d.HasChange("arc_agent_auto_upgrade_enabled") {

		autoUpgradeOption := arckubernetes.AutoUpgradeOptionsEnabled
		if !d.Get("arc_agent_auto_upgrade_enabled").(bool) {
			autoUpgradeOption = arckubernetes.AutoUpgradeOptionsDisabled
		}

		payload.Properties.ArcAgentProfile.AgentAutoUpgrade = pointer.To(autoUpgradeOption)
	}

	if err := client.ConnectedClusterCreateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return resourceArcKubernetesClusterRead(d, meta)
}

func resourceArcKubernetesClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandArcKubernetesClusterAadProfile(input []interface{}) *arckubernetes.AadProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	output := arckubernetes.AadProfile{
		EnableAzureRBAC: pointer.To(v["azure_rbac_enabled"].(bool)),
	}

	if tenantIdVal := v["tenant_id"].(string); tenantIdVal != "" {
		output.TenantID = pointer.To(tenantIdVal)
	}

	if groupVal := v["admin_group_object_ids"].([]interface{}); len(groupVal) != 0 {
		output.AdminGroupObjectIDs = utils.ExpandStringSlice(groupVal)
	}

	return &output
}

func flattenArcKubernetesClusterAadProfile(input *arckubernetes.AadProfile) []interface{} {
	if input == nil || (input.EnableAzureRBAC == nil && input.AdminGroupObjectIDs == nil && input.TenantID == nil) {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"azure_rbac_enabled":     pointer.From(input.EnableAzureRBAC),
			"admin_group_object_ids": utils.FlattenStringSlice(input.AdminGroupObjectIDs),
			"tenant_id":              pointer.From(input.TenantID),
		},
	}
}
