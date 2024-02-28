module github.com/ydb-platform/ydb-ops

go 1.21.6

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/ydb-platform/ydb-go-genproto v0.0.0-20240219184408-1de5f3f077de
	github.com/ydb-platform/ydb-rolling-restart v0.0.1
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/grpc v1.62.0 // indirect
)

replace github.com/ydb-platform/ydb-rolling-restart v0.0.1 => /home/jorres/work/nebius/nb-repos/ydb-rolling-restart