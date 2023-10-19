package flags

import (
	"fmt"
	"time"

	"github.com/urfave/cli"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	oppprof "github.com/ethereum-optimism/optimism/op-service/pprof"
	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
)

const EnvVarPrefix = "OP_PROPOSER"

var (
	// Required Flags
	L1EthRpcFlag = cli.StringFlag{
		Name:   "l1-eth-rpc",
		Usage:  "HTTP provider URL for L1. Multiple alternative addresses are supported, separated by commas, and the first address is used by default",
		EnvVar: opservice.PrefixEnvVar(EnvVarPrefix, "L1_ETH_RPC"),
	}
	RollupRpcFlag = cli.StringFlag{
		Name:   "rollup-rpc",
		Usage:  "HTTP provider URL for the rollup node",
		EnvVar: opservice.PrefixEnvVar(EnvVarPrefix, "ROLLUP_RPC"),
	}
	L2OOAddressFlag = cli.StringFlag{
		Name:   "l2oo-address",
		Usage:  "Address of the L2OutputOracle contract",
		EnvVar: opservice.PrefixEnvVar(EnvVarPrefix, "L2OO_ADDRESS"),
	}

	// Optional flags
	PollIntervalFlag = cli.DurationFlag{
		Name:   "poll-interval",
		Usage:  "How frequently to poll L2 for new blocks",
		Value:  6 * time.Second,
		EnvVar: opservice.PrefixEnvVar(EnvVarPrefix, "POLL_INTERVAL"),
	}
	AllowNonFinalizedFlag = cli.BoolFlag{
		Name:   "allow-non-finalized",
		Usage:  "Allow the proposer to submit proposals for L2 blocks derived from non-finalized L1 blocks.",
		EnvVar: opservice.PrefixEnvVar(EnvVarPrefix, "ALLOW_NON_FINALIZED"),
	}
	// Legacy Flags
	L2OutputHDPathFlag = txmgr.L2OutputHDPathFlag
)

var requiredFlags = []cli.Flag{
	L1EthRpcFlag,
	RollupRpcFlag,
	L2OOAddressFlag,
}

var optionalFlags = []cli.Flag{
	PollIntervalFlag,
	AllowNonFinalizedFlag,
	L2OutputHDPathFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oprpc.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, opmetrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oppprof.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, txmgr.CLIFlags(EnvVarPrefix)...)

	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.GlobalIsSet(f.GetName()) {
			return fmt.Errorf("flag %s is required", f.GetName())
		}
	}
	return nil
}
