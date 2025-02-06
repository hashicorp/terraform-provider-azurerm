// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualDesktopScalingPlanHostPoolAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopScalingPlanHostPoolAssociationCreate,
		Read:   resourceVirtualDesktopScalingPlanHostPoolAssociationRead,
		Update: resourceVirtualDesktopScalingPlanHostPoolAssociationUpdate,
		Delete: resourceVirtualDesktopScalingPlanHostPoolAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ScalingPlanHostPoolAssociationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"scaling_plan_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: scalingplan.ValidateScalingPlanID,
			},

			"host_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: scalingplan.ValidateHostPoolID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceVirtualDesktopScalingPlanHostPoolAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan <-> Host Pool Association creation.")
	scalingPlanId, err := scalingplan.ParseScalingPlanID(d.Get("scaling_plan_id").(string))
	if err != nil {
		return err
	}
	hostPoolId, err := scalingplan.ParseHostPoolID(d.Get("host_pool_id").(string))
	if err != nil {
		return err
	}
	associationId := parse.NewScalingPlanHostPoolAssociationId(*scalingPlanId, *hostPoolId).ID()

	locks.ByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)
	defer locks.UnlockByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)

	locks.ByName(hostPoolId.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(hostPoolId.HostPoolName, hostPoolResourceType)

	existing, err := client.Get(ctx, *scalingPlanId)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *scalingPlanId)
		}

		return fmt.Errorf("retrieving %s: %+v", *scalingPlanId, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *scalingPlanId)
	}
	model := *existing.Model

	hostPoolAssociations := []scalingplan.ScalingHostPoolReference{}
	if props := model.Properties; props != nil && props.HostPoolReferences != nil {
		hostPoolAssociations = *props.HostPoolReferences
	}

	hostPoolStr := hostPoolId.ID()
	if scalingPlanHostPoolAssociationExists(model.Properties, hostPoolStr) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan_host_pool_association", associationId)
	}
	hostPoolAssociations = append(hostPoolAssociations, scalingplan.ScalingHostPoolReference{
		HostPoolArmPath:    &hostPoolStr,
		ScalingPlanEnabled: utils.Bool(d.Get("enabled").(bool)),
	})

	payload := scalingplan.ScalingPlanPatch{
		Properties: &scalingplan.ScalingPlanPatchProperties{
			HostPoolReferences: &hostPoolAssociations,
			Schedules:          model.Properties.Schedules,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, *scalingPlanId, payload); err != nil {
		return fmt.Errorf("creating association between %s and %s: %+v", *scalingPlanId, *hostPoolId, err)
	}

	d.SetId(associationId)
	return resourceVirtualDesktopScalingPlanHostPoolAssociationRead(d, meta)
}

func resourceVirtualDesktopScalingPlanHostPoolAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanHostPoolAssociationID(d.Id())
	if err != nil {
		return err
	}

	scalingPlan, err := client.Get(ctx, id.ScalingPlan)
	if err != nil {
		if response.WasNotFound(scalingPlan.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.ScalingPlan)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
	}
	if model := scalingPlan.Model; model != nil {
		hostPoolId := id.HostPool.ID()
		exists := scalingPlanHostPoolAssociationExists(model.Properties, hostPoolId)
		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.ScalingPlan, id.HostPool)
			d.SetId("")
			return nil
		}
		if props := model.Properties; props != nil && props.HostPoolReferences != nil {
			for _, referenceId := range *props.HostPoolReferences {
				if referenceId.HostPoolArmPath != nil {
					if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
						d.Set("enabled", referenceId.ScalingPlanEnabled)
					}
				}
			}
		}

		d.Set("scaling_plan_id", id.ScalingPlan.ID())
		d.Set("host_pool_id", hostPoolId)
	}

	return nil
}

func resourceVirtualDesktopScalingPlanHostPoolAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanHostPoolAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)
	defer locks.UnlockByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)

	locks.ByName(id.HostPool.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPool.HostPoolName, hostPoolResourceType)

	existing, err := client.Get(ctx, id.ScalingPlan)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", id.ScalingPlan)
		}

		return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id.ScalingPlan)
	}
	model := *existing.Model
	if !scalingPlanHostPoolAssociationExists(model.Properties, id.HostPool.ID()) {
		log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.ScalingPlan, id.HostPool)
		d.SetId("")
		return nil
	}

	hostPoolReferences := []scalingplan.ScalingHostPoolReference{}
	hostPoolId := id.HostPool.ID()
	if props := model.Properties; props != nil && props.HostPoolReferences != nil {
		for _, referenceId := range *props.HostPoolReferences {
			if referenceId.HostPoolArmPath != nil {
				if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
					referenceId.ScalingPlanEnabled = utils.Bool(d.Get("enabled").(bool))
				}
			}
			hostPoolReferences = append(hostPoolReferences, referenceId)
		}
	}

	payload := scalingplan.ScalingPlanPatch{
		Properties: &scalingplan.ScalingPlanPatchProperties{
			HostPoolReferences: &hostPoolReferences,
			Schedules:          model.Properties.Schedules,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, id.ScalingPlan, payload); err != nil {
		return fmt.Errorf("updating association between %s and %s: %+v", id.ScalingPlan, id.HostPool, err)
	}

	return nil
}

func resourceVirtualDesktopScalingPlanHostPoolAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanHostPoolAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)
	defer locks.UnlockByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)

	locks.ByName(id.HostPool.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPool.HostPoolName, hostPoolResourceType)

	existing, err := client.Get(ctx, id.ScalingPlan)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", id.ScalingPlan)
		}

		return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id.ScalingPlan)
	}
	model := *existing.Model

	hostPoolReferences := []scalingplan.ScalingHostPoolReference{}
	hostPoolId := id.HostPool.ID()
	if props := model.Properties; props != nil && props.HostPoolReferences != nil {
		for _, referenceId := range *props.HostPoolReferences {
			if referenceId.HostPoolArmPath != nil {
				if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
					continue
				}
			}

			hostPoolReferences = append(hostPoolReferences, referenceId)
		}
	}

	payload := scalingplan.ScalingPlanPatch{
		Properties: &scalingplan.ScalingPlanPatchProperties{
			HostPoolReferences: &hostPoolReferences,
			Schedules:          model.Properties.Schedules,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, id.ScalingPlan, payload); err != nil {
		return fmt.Errorf("removing association between %s and %s: %+v", id.ScalingPlan, id.HostPool, err)
	}

	return nil
}

func scalingPlanHostPoolAssociationExists(props *scalingplan.ScalingPlanProperties, applicationGroupId string) bool {
	if props == nil || props.HostPoolReferences == nil {
		return false
	}

	for _, id := range *props.HostPoolReferences {
		if id.HostPoolArmPath != nil {
			if strings.EqualFold(*id.HostPoolArmPath, applicationGroupId) {
				return true
			}
		}
	}

	return false
}
