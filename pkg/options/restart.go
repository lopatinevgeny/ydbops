package options

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/ydb-platform/ydb-go-genproto/draft/protos/Ydb_Maintenance"
	"github.com/ydb-platform/ydb-ops/internal/util"
	"google.golang.org/protobuf/types/known/durationpb"

	cmsOptions "github.com/ydb-platform/ydb-rolling-restart/pkg/options"
)

var AvailabilityModes = []string{"strong", "weak", "force"}

type RestartOptions struct {
	CMS     *cmsOptions.CMS
	GRPC    *cmsOptions.GRPC

	AvailabilityMode   string
	Tenants            []string
	Hosts              []string
	ExcludeHosts       []string
	RestartDuration    int
	RestartRetryNumber int

	Continue bool
}

var RestartOptionsInstance = &RestartOptions{
	CMS: &cmsOptions.CMS{},
	GRPC: &cmsOptions.GRPC{},
}

func (o *RestartOptions) Validate() error {
	if !util.Contains(AvailabilityModes, o.AvailabilityMode) {
		return fmt.Errorf("specified not supported availability mode: %s", o.AvailabilityMode)
	}

	if o.RestartDuration < 0 {
		return fmt.Errorf("specified invalid restart duration seconds: %d. Must be positive", o.RestartDuration)
	}

	if o.RestartRetryNumber < 0 {
		return fmt.Errorf("specified invalid restart retry number: %d. Must be positive", o.RestartRetryNumber)
	}

	if _, err := o.GetNodeIds(); err != nil {
		return err
	}

	if err := o.CMS.Validate(); err != nil {
		return err
	}
	if err := o.GRPC.Validate(); err != nil {
		return err
	}

	return nil
}

func (o *RestartOptions) DefineFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Continue, "continue", /* TODO: false??? */ false, "TODO Continue previous rolling restart")

	fs.StringSliceVar(&o.ExcludeHosts, "exclude-hosts", []string{}, "TODO Never restart these hosts")

	fs.StringVarP(&o.AvailabilityMode, "availability-mode", "", AvailabilityModes[0],
		fmt.Sprintf("Availability mode. Available choices: %s", strings.Join(AvailabilityModes, ", ")))

	fs.IntVarP(&o.RestartDuration, "restart-duration", "", 60,
		"Restart duration in seconds")

	fs.IntVarP(&o.RestartRetryNumber, "restart-retry-number", "", 3,
		"Retry number of restart")

	fs.StringArrayVarP(&o.Tenants, "tenants", "", o.Tenants,
		"Restart only specified tenants")

	fs.StringArrayVarP(&o.Hosts, "nodes", "", o.Hosts,
		"Restart only specified nodes")

	o.CMS.DefineFlags(fs)
	o.GRPC.DefineFlags(fs)
}

func (o *RestartOptions) GetAvailabilityMode() Ydb_Maintenance.AvailabilityMode {
	title := strings.ToUpper(fmt.Sprintf("availability_mode_%s", o.AvailabilityMode))
	value := Ydb_Maintenance.AvailabilityMode_value[title]

	return Ydb_Maintenance.AvailabilityMode(value)
}

func (o *RestartOptions) GetRestartDuration() *durationpb.Duration {
	return durationpb.New(time.Second * time.Duration(o.RestartDuration) * time.Duration(o.RestartRetryNumber))
}

func (o *RestartOptions) GetNodeIds() ([]uint32, error) {
	ids := make([]uint32, 0, len(o.Hosts))

	for _, nodeId := range o.Hosts {
		id, err := strconv.Atoi(nodeId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse node id: %+v", err)
		}
		if id < 0 {
			return nil, fmt.Errorf("invalid node id specified: %d, must be positive", id)
		}
		ids = append(ids, uint32(id))
	}

	return ids, nil
}