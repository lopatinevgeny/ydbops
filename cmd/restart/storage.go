package restart

import (
	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-ops/internal/cobra_util"
	"github.com/ydb-platform/ydb-ops/pkg/options"
)

func NewStorageCmd() *cobra.Command {
	restartOpts := options.RestartOptionsInstance

	cmd := cobra_util.SetDefaultsOn(&cobra.Command{
		Use:   "storage",
		Short: "Restarts a specified subset of tenant nodes",
    Long:  `ydb-ops restart storage:
  Restarts a specified subset of storage nodes`,
	}, restartOpts)

	return cmd
}

func init() {
}
