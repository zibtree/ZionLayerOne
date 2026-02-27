package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zionlayer/zionlayer/consensus"
	"github.com/zionlayer/zionlayer/core/mempool"
	"github.com/zionlayer/zionlayer/core/state"
	"github.com/zionlayer/zionlayer/core/transaction"
	"github.com/zionlayer/zionlayer/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "ziond",
	Short: "ZionLayer Node",
	Long:  "ZionLayer — The General-Purpose AI-Native Layer 1 Blockchain",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the ZionLayer node",
	RunE:  runNode,
}

var (
	flagRPCPort       int
	flagValidatorAddr string
	flagDataDir       string
)

func init() {
	startCmd.Flags().IntVar(&flagRPCPort, "rpc-port", 8545, "JSON-RPC port")
	startCmd.Flags().StringVar(&flagValidatorAddr, "validator", "", "Validator address")
	startCmd.Flags().StringVar(&flagDataDir, "data-dir", "./data", "Data directory")
	rootCmd.AddCommand(startCmd)
}

func runNode(cmd *cobra.Command, args []string) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("⛓️  ZionLayer starting",
		zap.String("version", "0.1.0"),
		zap.Int("rpc-port", flagRPCPort),
	)

	// Initialize components
	stateDB := state.NewStateDB()
	pool := mempool.NewPool()
	engine := consensus.NewZionBFT(stateDB, logger)

	// Start consensus (tx feed channel)
	txFeed := make(chan []*transaction.Tx, 10)
	validatorAddr := flagValidatorAddr
	if validatorAddr == "" {
		validatorAddr = "0xDevnetValidator0000000000000000000000001"
	}
	engine.Start(validatorAddr, txFeed)

	// Feed mempool batches to consensus
	go func() {
		for {
			batch := pool.Pop(100)
			if len(batch) > 0 {
				txFeed <- batch
			}
		}
	}()

	// Log finalized blocks
	go func() {
		for b := range engine.Blocks() {
			logger.Info("✅ block finalized",
				zap.Uint64("height", b.Header.Height),
				zap.Int("txs", len(b.Txs)),
			)
		}
	}()

	// Start RPC server in background
	rpcServer := rpc.NewServer(stateDB, pool, logger, flagRPCPort)
	go func() {
		if err := rpcServer.Start(); err != nil {
			logger.Fatal("RPC server error", zap.Error(err))
		}
	}()

	logger.Info("🚀 node ready",
		zap.String("rpc", fmt.Sprintf("http://localhost:%d", flagRPCPort)),
	)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	engine.Stop()
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
