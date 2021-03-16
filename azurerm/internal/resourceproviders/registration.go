package resourceproviders

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
)

func EnsureRegistered(ctx context.Context, client resources.ProvidersClient, availableRPs []resources.Provider, requiredRPs map[string]struct{}) error {
	log.Printf("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister := resourceproviders.DetermineResourceProvidersRequiringRegistration(availableRPs, requiredRPs)

	if len(providersToRegister) > 0 {
		log.Printf("[DEBUG] Registering %d Resource Providers", len(providersToRegister))
		if err := resourceproviders.RegisterForSubscription(ctx, client, providersToRegister); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] All required Resource Providers are registered")
	}

	return nil
}
