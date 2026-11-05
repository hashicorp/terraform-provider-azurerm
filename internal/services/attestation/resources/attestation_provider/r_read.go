// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func resourceAttestationProviderRead(d *pluginsdk.ResourceData, meta any) error {
	attestationClients := meta.(*clients.Client).Attestation
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	resp, err := attestationClients.ProviderClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceAttestationProviderFlatten(ctx, attestationClients, d, id, resp.Model, true)
}

func base64DataFromAttestationJWT(input *string) (*string, error) {
	if input == nil {
		return nil, nil
	}

	split := strings.Split(*input, ".")
	if len(split) != 3 {
		return nil, fmt.Errorf("expected the first token to have 3 segments but got %d", len(split))
	}
	// decode the JWT into a PolicyResult object
	decodedJwtSegment, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		return nil, fmt.Errorf("base64-decoding the first JWT Segment %q: %+v", split[1], err)
	}
	var firstResult attestation.PolicyResult
	if err := json.Unmarshal(decodedJwtSegment, &firstResult); err != nil {
		return nil, fmt.Errorf("unmarshaling into PolicyResult: %+v", err)
	}
	if firstResult.Policy == nil {
		return nil, nil
	}

	return firstResult.Policy, nil
}
