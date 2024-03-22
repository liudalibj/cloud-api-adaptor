//go:build azure

// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package azure

import (
	pv "github.com/liudalibj/cloud-api-adaptor/src/cloud-api-adaptor/test/provisioner"
)

func init() {
	pv.NewProvisionerFunctions["azure"] = NewAzureCloudProvisioner
	pv.NewInstallOverlayFunctions["azure"] = NewAzureInstallOverlay
}
