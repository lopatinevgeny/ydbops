package tests

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ydb-platform/ydbops/pkg/options"
	"github.com/ydb-platform/ydbops/pkg/rolling/restarters"
	"github.com/ydb-platform/ydbops/tests/mock"
	"go.uber.org/zap"
)

var _ = Describe("Test storage Filter", func() {
	var (
		now                     = time.Now()
		tenMinutesAgoTimestamp  = now.Add(-10 * time.Minute)
		fiveMinutesAgoTimestamp = now.Add(-5 * time.Minute)
	)

	It("baremetal restarter filtering by --started>timestamp", func() {
		restarter := restarters.NewStorageBaremetalRestarter(zap.S())

		nodeGroups := [][]uint32{
			{1, 2, 3, 4, 5, 6, 7, 8},
		}
		nodeInfoMap := map[uint32]mock.TestNodeInfo{
			1: {
				StartTime: tenMinutesAgoTimestamp,
			},
			2: {
				StartTime: tenMinutesAgoTimestamp,
			},
			3: {
				StartTime: tenMinutesAgoTimestamp,
			},
		}

		nodes := mock.CreateNodesFromShortConfig(nodeGroups, nodeInfoMap)

		filterSpec := restarters.FilterNodeParams{
			StartedTime: &options.StartedTime{
				Direction: '<',
				Timestamp: fiveMinutesAgoTimestamp,
			},
		}

		clusterInfo := restarters.ClusterNodesInfo{
			AllNodes:        nodes,
			TenantToNodeIds: map[string][]uint32{},
		}

		filteredNodes := restarter.Filter(filterSpec, clusterInfo)

		Expect(len(filteredNodes)).To(Equal(3))

		filteredNodeIds := make(map[uint32]bool)
		for _, node := range filteredNodes {
			filteredNodeIds[node.NodeId] = true
		}

		Expect(filteredNodeIds).Should(HaveKey(uint32(1)))
		Expect(filteredNodeIds).Should(HaveKey(uint32(2)))
		Expect(filteredNodeIds).Should(HaveKey(uint32(3)))
	})

	It("baremetal restarter without arguments takes all storage nodes", func() {
		restarter := restarters.NewStorageBaremetalRestarter(zap.S())

		nodeGroups := [][]uint32{
			{1, 2, 3, 4, 5, 6, 7, 8},
			{9, 10, 11},
		}
		nodeInfoMap := map[uint32]mock.TestNodeInfo{
			9: {
				IsDynnode:  true,
				TenantName: "fakeTenant",
			},
			10: {
				IsDynnode:  true,
				TenantName: "fakeTenant",
			},
			11: {
				IsDynnode:  true,
				TenantName: "fakeTenant",
			},
		}

		nodes := mock.CreateNodesFromShortConfig(nodeGroups, nodeInfoMap)

		// empty params equivalent to no arguments
		filterSpec := restarters.FilterNodeParams{}

		clusterInfo := restarters.ClusterNodesInfo{
			AllNodes: nodes,
			TenantToNodeIds: map[string][]uint32{
				"fakeTenant": {9, 10, 11},
			},
		}

		filteredNodes := restarter.Filter(filterSpec, clusterInfo)

		Expect(len(filteredNodes)).To(Equal(8))

		filteredNodeIds := make(map[uint32]bool)
		for _, node := range filteredNodes {
			filteredNodeIds[node.NodeId] = true
		}

    for i := 1; i <= 8; i++ {
      Expect(filteredNodeIds).Should(HaveKey(uint32(i)))
    }
	})
})