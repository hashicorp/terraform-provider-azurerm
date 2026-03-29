// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestMultiCloudConnectorARMResourceID(t *testing.T) {
	t.Parallel()

	valid := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.HybridConnectivity/publicCloudConnectors/c1"
	_, errs := MultiCloudConnectorARMResourceID(valid, "multi_cloud_connector_id")
	if len(errs) > 0 {
		t.Fatalf("expected valid id, got %v", errs)
	}

	_, errs = MultiCloudConnectorARMResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sa", "multi_cloud_connector_id")
	if len(errs) == 0 {
		t.Fatal("expected error for wrong resource type")
	}
}

func TestAwsS3BucketARMResourceID(t *testing.T) {
	t.Parallel()

	valid := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.AwsConnector/s3Buckets/b1"
	_, errs := AwsS3BucketARMResourceID(valid, "aws_s3_bucket_id")
	if len(errs) > 0 {
		t.Fatalf("expected valid id, got %v", errs)
	}

	_, errs = AwsS3BucketARMResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Storage/storageAccounts/sa", "aws_s3_bucket_id")
	if len(errs) == 0 {
		t.Fatal("expected error for wrong resource type")
	}
}
