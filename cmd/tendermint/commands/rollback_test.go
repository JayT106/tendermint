package commands_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/rpc/client/local"
	rpctest "github.com/tendermint/tendermint/rpc/test"
	e2e "github.com/tendermint/tendermint/test/e2e/app"
)

func TestRollbackIntegration(t *testing.T) {
	var height int64
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := rpctest.CreateConfig(t, t.Name())
	require.NoError(t, err)
	cfg.BaseConfig.DBBackend = "goleveldb"

	app, err := e2e.NewApplication(e2e.DefaultConfig(dir))
	require.NoError(t, err)

	t.Run("First run", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		require.NoError(t, err)
		node, _, err := rpctest.StartTendermint(ctx, cfg, app, rpctest.SuppressStdout)
		require.NoError(t, err)
		require.True(t, node.IsRunning())

		time.Sleep(3 * time.Second)
		cancel()
		node.Wait()

		require.False(t, node.IsRunning())
	})
	t.Run("Rollback", func(t *testing.T) {
		time.Sleep(time.Second)
		require.NoError(t, app.Rollback())
		height, _, err = commands.RollbackState(cfg)
		require.NoError(t, err, "%d", height)
		require.Equal(t, height, app.Info(types.RequestInfo{}).LastBlockHeight)
	})
	t.Run("Rollback for appHash mismatch case", func(t *testing.T) {
		require.NoError(t, app.Rollback())
		height2, _, err := commands.RollbackState(cfg)
		require.NoError(t, err, "%d", height2)
		require.Equal(t, height-1, height2)
		require.Equal(t, height2, app.Info(types.RequestInfo{}).LastBlockHeight)

		// reset the pval state for the pval can vote the new block height
		pval, err := privval.LoadOrGenFilePV(cfg.PrivValidator.KeyFile(), cfg.PrivValidator.StateFile())
		require.NoError(t, err)
		require.NoError(t, pval.Reset())
	})
	t.Run("Restart", func(t *testing.T) {
		require.True(t, height > 0, "%d", height)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		node2, _, err2 := rpctest.StartTendermint(ctx, cfg, app, rpctest.SuppressStdout)
		require.NoError(t, err2)
		t.Cleanup(node2.Wait)

		logger := log.NewNopLogger()

		client, err := local.New(logger, node2.(local.NodeService))
		require.NoError(t, err)

		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				t.Fatalf("failed to make progress after 10 seconds. Min height: %d", height)
			case <-ticker.C:
				status, err := client.Status(ctx)
				require.NoError(t, err)

				if status.SyncInfo.LatestBlockHeight > height {
					return
				}
			}
		}
	})

}
