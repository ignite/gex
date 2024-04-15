package explorer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cometbft/cometbft/libs/json"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/gex/pkg/client"
	"github.com/ignite/gex/pkg/number"
	"github.com/ignite/gex/pkg/widget"
)

const (
	statusConnected    = "✔️ good"
	statusNotConnected = "✖️ not connected"
)

var (
	RoundStepPropose   = strings.ToUpper("RoundStepPropose")
	RoundStepPreVote   = strings.ToUpper("RoundStepPrevote")
	RoundStepPreCommit = strings.ToUpper("RoundStepPrecommit")
	RoundStepCommit    = strings.ToUpper("RoundStepCommit")
	RoundStepNewHeight = strings.ToUpper("RoundStepNewHeight")
)

// info holds all cross infos.
type info struct {
	sync.RWMutex
	blocks          int64
	maxGasWanted    int64
	transactions    int64
	totalGasWanted  int64
	lastTxGasWanted int64
}

// Run runs the explorer view listening to the provided host.
func Run(ctx context.Context, host string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c, err := client.New(ctx, host)
	if err != nil {
		return err
	}

	w, err := widget.New()
	if err != nil {
		return err
	}

	var (
		info  = &info{}
		start = time.Now()
	)

	errGroup, _ := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		status, err := c.Status(ctx)
		if err != nil {
			return err
		}

		if err := w.SetCurrentNetwork(status.NodeInfo.Network); err != nil {
			return err
		}
		return w.SetMoniker(status.NodeInfo.Moniker)
	})

	client.Callback(ctx, 1*time.Second, func() error {
		now := time.Now()
		if err := w.SetTime(now.Format("2006-01-02\n15:04:05")); err != nil {
			return err
		}
		secondsPassed := now.Sub(start).Seconds()

		info.RLock()
		var (
			lastTxGasWanted  = info.lastTxGasWanted
			maxGasWanted     = info.maxGasWanted
			blocksPerSecond  = 0.0
			totalGasPerBlock = int64(0)
			averageGasPerTx  = int64(0)
		)
		if info.blocks > 0 {
			totalGasPerBlock = info.totalGasWanted / info.blocks
			blocksPerSecond = secondsPassed / float64(info.blocks)
		}
		if info.transactions > 0 {
			averageGasPerTx = info.totalGasWanted / info.transactions
		}
		info.RUnlock()

		if err := w.SetSecondsPerBlock(fmt.Sprintf("%.2f seconds", blocksPerSecond)); err != nil {
			return err
		}

		if err := w.SetGasMax(number.WithComma(maxGasWanted)); err != nil {
			return err
		}

		if err := w.SetGasAvgBlock(number.WithComma(totalGasPerBlock)); err != nil {
			return err
		}

		if err := w.SetLatestGas(number.WithComma(lastTxGasWanted)); err != nil {
			return err
		}

		return w.SetGasAvgTransaction(number.WithComma(averageGasPerTx))
	})

	c.ConsensusParams(ctx, func(params coretypes.ResultConsensusParams) error {
		info.Lock()
		info.maxGasWanted = params.ConsensusParams.Block.MaxGas
		info.Unlock()
		return w.SetMaxBlockSize(number.ByteCountDecimal(params.ConsensusParams.Block.MaxBytes))
	})

	c.NetInfo(ctx, func(info coretypes.ResultNetInfo) error {
		return w.SetPeers(info.NPeers)
	})

	c.Health(ctx, func(health *coretypes.ResultHealth, err error) error {
		if health != nil && err == nil {
			return w.SetHealth(statusConnected)
		}
		return w.SetHealth(statusNotConnected)
	})

	c.Validators(ctx, func(validators coretypes.ResultValidators) error {
		return w.SetValidators(validators.Total)
	})

	err = c.NewRoundStep(ctx, func(state types.EventDataRoundState) error {
		progress := 0
		switch strings.ToUpper(state.Step) {
		case RoundStepPropose:
			progress = 20
		case RoundStepPreVote:
			progress = 40
		case RoundStepPreCommit:
			progress = 60
		case RoundStepCommit:
			progress = 80
		case RoundStepNewHeight:
			progress = 100
		}
		return w.SetBlockProgress(progress)
	})
	if err != nil {
		return err
	}

	err = c.NewBlock(ctx, func(block types.EventDataNewBlock) error {
		info.Lock()
		info.blocks++
		info.Unlock()

		return w.AddBlock(
			fmt.Sprintf(
				"%d %s txs:%d",
				block.Block.Height,
				block.Block.Header.Hash(),
				block.Block.Txs.Len(),
			),
		)
	})
	if err != nil {
		return err
	}

	err = c.Tx(ctx, func(tx types.EventDataTx) error {
		info.Lock()
		info.transactions++
		info.lastTxGasWanted = tx.Result.GasWanted
		info.totalGasWanted += info.lastTxGasWanted
		info.Unlock()

		result, err := json.Marshal(tx.Result)
		if err != nil {
			return err
		}
		return w.AddTransaction(string(result))
	})
	if err != nil {
		return err
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return w.Run(ctx)
}
