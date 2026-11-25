// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket_with_server

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandNetAppBucketServer(input []models.NetAppVolumeBucketServer) *buckets.BucketServerProperties {
	if len(input) == 0 {
		return nil
	}

	srv := input[0]
	out := &buckets.BucketServerProperties{}

	if srv.Fqdn != "" {
		out.Fqdn = pointer.To(srv.Fqdn)
	}
	if srv.CertificatePem != "" {
		out.CertificateObject = pointer.To(srv.CertificatePem)
	}
	if srv.OnCertificateConflictAction != "" {
		out.OnCertificateConflictAction = pointer.ToEnum[buckets.OnCertificateConflictAction](srv.OnCertificateConflictAction)
	}

	return out
}

func expandNetAppBucketServerPatch(input []models.NetAppVolumeBucketServer, rd *pluginsdk.ResourceData) *buckets.BucketServerPatchProperties {
	if len(input) == 0 {
		return nil
	}

	srv := input[0]
	out := &buckets.BucketServerPatchProperties{}

	if rd.HasChange("server.0.fqdn") {
		out.Fqdn = pointer.To(srv.Fqdn)
	}
	if rd.HasChange("server.0.certificate_pem") && srv.CertificatePem != "" {
		out.CertificateObject = pointer.To(srv.CertificatePem)
	}
	if rd.HasChange("server.0.on_certificate_conflict_action") && srv.OnCertificateConflictAction != "" {
		out.OnCertificateConflictAction = pointer.ToEnum[buckets.OnCertificateConflictAction](srv.OnCertificateConflictAction)
	}

	return out
}

func flattenNetAppBucketServer(input *buckets.BucketServerProperties) []models.NetAppVolumeBucketServer {
	if input == nil {
		return []models.NetAppVolumeBucketServer{}
	}

	out := models.NetAppVolumeBucketServer{
		Fqdn: pointer.From(input.Fqdn),
	}
	if input.OnCertificateConflictAction != nil {
		out.OnCertificateConflictAction = pointer.FromEnum(input.OnCertificateConflictAction)
	}

	if out.Fqdn == "" && out.OnCertificateConflictAction == "" {
		return []models.NetAppVolumeBucketServer{}
	}

	return []models.NetAppVolumeBucketServer{out}
}
