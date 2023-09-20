package main

import (
	"os"

	log "github.com/ethereum/go-ethereum/log"
	cli "github.com/urfave/cli"

	watch "github.com/ethereum-optimism/optimism/op-challenger/cmd/watch"
	config "github.com/ethereum-optimism/optimism/op-challenger/config"
	flags "github.com/ethereum-optimism/optimism/op-challenger/flags"
	version "github.com/ethereum-optimism/optimism/op-challenger/version"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

var (
	GitCommit = ""
	GitDate   = ""
)

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := version.Version
	if GitCommit != "" {
		v += "-" + GitCommit[:8]
	}
	if GitDate != "" {
		v += "-" + GitDate
	}
	if version.Meta != "" {
		v += "-" + version.Meta
	}
	return v
}()

func main() {
	args := os.Args
	if err := run(args, Main); err != nil {
		log.Crit("Application failed", "err", err)
	}
}

type ConfigAction func(log log.Logger, version string, config *config.Config) error

// run parses the supplied args to create a config.Config instance, sets up logging
// then calls the supplied ConfigAction.
// This allows testing the translation from CLI arguments to Config
func run(args []string, action ConfigAction) error {
	// Set up logger with a default INFO level in case we fail to parse flags,
	// otherwise the final critical log won't show what the parsing error was.
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Version = VersionWithMeta
	app.Flags = flags.Flags
	app.Name = "op-challenger"
	app.Usage = "Challenge Invalid L2OutputOracle Outputs"
	app.Description = "A modular op-stack challenge agent for dispute games written in golang."
	app.Action = func(ctx *cli.Context) error {
		logger, err := config.LoggerFromCLI(ctx)
		if err != nil {
			return err
		}
		logger.Info("Starting challenger", "version", VersionWithMeta)

		cfg, err := config.NewConfigFromCLI(ctx)
		if err != nil {
			return err
		}
		return action(logger, VersionWithMeta, cfg)
	}
	app.Commands = []cli.Command{
		{
			Name:        "watch",
			Subcommands: watch.Subcommands,
		},
	}

	return app.Run(args)
}
