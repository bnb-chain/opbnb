package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"os"

	"github.com/ethereum-optimism/optimism/op-chain-ops/safe"
	"github.com/ethereum-optimism/optimism/op-service/jsonutil"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

func main() {
	color := isatty.IsTerminal(os.Stderr.Fd())
	oplog.SetGlobalLogHandler(log.NewTerminalHandler(os.Stderr, color))

	app := &cli.App{
		Name:  "opbnb-upgrade",
		Usage: "Build transactions useful for upgrading the opBNB",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "l1-rpc-url",
				Value:   "http://127.0.0.1:8545",
				Usage:   "L1 RPC URL, the chain ID will be used to determine the opBNB",
				EnvVars: []string{"L1_RPC_URL"},
			},
			&cli.PathFlag{
				Name:    "outfile",
				Usage:   "The file to write the output to. If not specified, output is written to stdout",
				EnvVars: []string{"OUTFILE"},
			},
		},
		Action: entrypoint,
	}

	if err := app.Run(os.Args); err != nil {
		log.Crit("error opBNB-upgrade", "err", err)
	}
}

// entrypoint contains the main logic of the script
func entrypoint(ctx *cli.Context) error {
	client, err := ethclient.Dial(ctx.String("l1-rpc-url"))
	if err != nil {
		return err
	}

	// Fetch the L1 chain ID to determine the superchain name
	l1ChainID, err := client.ChainID(ctx.Context)
	if err != nil {
		return err
	}

	proxyAddresses := BscQAnetProxyContracts
	implAddresses := BscQAnetImplContracts
	if l1ChainID.Uint64() == bscTestnet {
		proxyAddresses = BscTestnetProxyContracts
		implAddresses = BscTestnetImplContracts
		fmt.Println("upgrade bscTestnet")
	} else if l1ChainID.Uint64() == bscMainnet {
		proxyAddresses = BscMainnetProxyContracts
		implAddresses = BscMainnetImplContracts
		fmt.Println("upgrade bscMainnet")
	} else {
		fmt.Println("upgrade BscQAnet")
	}

	versions, err := GetContractVersions(ctx.Context, proxyAddresses, client)
	log.Info("current contract versions")
	log.Info("L1CrossDomainMessenger", "version", versions.L1CrossDomainMessenger, "address", proxyAddresses["L1CrossDomainMessengerProxy"])
	log.Info("L1ERC721Bridge", "version", versions.L1ERC721Bridge, "address", proxyAddresses["L1ERC721BridgeProxy"])
	log.Info("L1StandardBridge", "version", versions.L1StandardBridge, "address", proxyAddresses["L1StandardBridgeProxy"])
	log.Info("L2OutputOracle", "version", versions.L2OutputOracle, "address", proxyAddresses["L2OutputOracleProxy"])
	log.Info("OptimismMintableERC20Factory", "version", versions.OptimismMintableERC20Factory, "address", proxyAddresses["OptimismMintableERC20FactoryProxy"])
	log.Info("OptimismPortal", "version", versions.OptimismPortal, "address", proxyAddresses["OptimismPortalProxy"])
	log.Info("SystemConfig", "version", versions.SystemConfig, "address", proxyAddresses["SystemConfigProxy"])

	versions, err = GetContractVersions(ctx.Context, implAddresses, client)
	log.Info("Upgrading to the following versions")
	log.Info("L1CrossDomainMessenger", "version", versions.L1CrossDomainMessenger, "address", proxyAddresses["L1CrossDomainMessengerProxy"])
	log.Info("L1ERC721Bridge", "version", versions.L1ERC721Bridge, "address", proxyAddresses["L1ERC721BridgeProxy"])
	log.Info("L1StandardBridge", "version", versions.L1StandardBridge, "address", proxyAddresses["L1StandardBridgeProxy"])
	log.Info("L2OutputOracle", "version", versions.L2OutputOracle, "address", proxyAddresses["L2OutputOracleProxy"])
	log.Info("OptimismMintableERC20Factory", "version", versions.OptimismMintableERC20Factory, "address", proxyAddresses["OptimismMintableERC20FactoryProxy"])
	log.Info("OptimismPortal", "version", versions.OptimismPortal, "address", proxyAddresses["OptimismPortalProxy"])
	log.Info("SystemConfig", "version", versions.SystemConfig, "address", proxyAddresses["SystemConfigProxy"])

	// Create a batch of transactions
	batch := safe.Batch{}

	// Build the batch
	if err := L1(&batch, proxyAddresses, implAddresses, client, l1ChainID); err != nil {
		return err
	}

	// Write the batch to disk or stdout
	if outfile := ctx.Path("outfile"); outfile != "" {
		if err := jsonutil.WriteJSON(outfile, batch, 0o666); err != nil {
			return err
		}
	} else {
		data, err := json.MarshalIndent(batch, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	}
	return nil

}
