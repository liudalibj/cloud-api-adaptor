// (C) Copyright IBM Corp. 2022.
// SPDX-License-Identifier: Apache-2.0

package routing

import (
	"testing"

	testutils "github.com/liudalibj/cloud-api-adaptor/src/cloud-api-adaptor/pkg/internal/testing"
	"github.com/liudalibj/cloud-api-adaptor/src/cloud-api-adaptor/pkg/podnetwork/tuntest"
)

func TestRouting(t *testing.T) {
	// TODO: enable this test once https://github.com/liudalibj/cloud-api-adaptor/issues/52 is fixed
	testutils.SkipTestIfRunningInCI(t)

	tuntest.RunTunnelTest(t, "routing", NewWorkerNodeTunneler, NewPodNodeTunneler, true)

}
