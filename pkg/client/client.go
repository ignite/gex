package client

import (
	"context"
	"time"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const tickerTime = 1 * time.Second

// Client gex client.
type Client struct {
	cosmosclient.Client
}

// New creates a new Client.
func New(ctx context.Context, host string) (Client, error) {
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(host))
	if err != nil {
		return Client{}, err
	}

	if err := client.RPC.Start(); err != nil {
		return Client{}, err
	}

	return Client{Client: client}, nil
}

// NewBlock listen the new block event from the websocket subscriber.
func (c Client) NewBlock(ctx context.Context, fn func(types.EventDataNewBlock) error) error {
	return c.Subscribe(
		ctx,
		types.EventNewBlock,
		types.EventQueryNewBlock.String(),
		func(event coretypes.ResultEvent) error {
			blockEvent, ok := event.Data.(types.EventDataNewBlock)
			if !ok {
				return errors.Errorf("invalid event new block type: %v", event.Data)
			}
			return fn(blockEvent)
		},
	)
}

// NewRoundStep listen the new round step event from the websocket subscriber.
func (c Client) NewRoundStep(ctx context.Context, fn func(types.EventDataRoundState) error) error {
	return c.Subscribe(
		ctx,
		types.EventNewRoundStep,
		types.EventQueryNewRoundStep.String(),
		func(event coretypes.ResultEvent) error {
			newRoundEvent, ok := event.Data.(types.EventDataRoundState)
			if !ok {
				return errors.Errorf("invalid event new round step type: %v", event.Data)
			}
			return fn(newRoundEvent)
		},
	)
}

// Tx listen the new transaction event from the websocket subscriber.
func (c Client) Tx(ctx context.Context, fn func(types.EventDataTx) error) error {
	return c.Subscribe(
		ctx,
		types.EventTx,
		types.EventQueryTx.String(),
		func(event coretypes.ResultEvent) error {
			txEvent, ok := event.Data.(types.EventDataTx)
			if !ok {
				return errors.Errorf("invalid event tx type: %v", event.Data)
			}
			return fn(txEvent)
		},
	)
}

// Subscribe listen websocket events based in the query.
func (c Client) Subscribe(ctx context.Context, subscriber, query string, fn func(coretypes.ResultEvent) error) error {
	out, err := c.RPC.Subscribe(ctx, subscriber, query)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resultEvent := <-out:
				if err := fn(resultEvent); err != nil {
					// TODO find a better way to send logs
					// fmt.Println(err)
					return
				}
			}
		}
	}()

	return nil
}

// NetInfo fetch the network information for each new block.
func (c Client) NetInfo(ctx context.Context, fn func(coretypes.ResultNetInfo) error) {
	c.BlockCallback(ctx, func(int64) error {
		netInfo, err := c.RPC.NetInfo(ctx)
		if err != nil {
			return err
		}
		return fn(*netInfo)
	})
}

// Health fetch the health information for each new block.
func (c Client) Health(ctx context.Context, fn func(*coretypes.ResultHealth, error) error) {
	c.BlockCallback(ctx, func(int64) error {
		return fn(c.RPC.Health(ctx))
	})
}

// Validators fetch the validators information for each new block.
func (c Client) Validators(ctx context.Context, fn func(coretypes.ResultValidators) error) {
	c.BlockCallback(ctx, func(height int64) error {
		page := 1
		count := 1_000
		validators, err := c.RPC.Validators(ctx, &height, &page, &count)
		if err != nil {
			return err
		}
		return fn(*validators)
	})
}

// ConsensusParams fetch the consensus parameters for each new block.
func (c Client) ConsensusParams(ctx context.Context, fn func(coretypes.ResultConsensusParams) error) {
	c.BlockCallback(ctx, func(height int64) error {
		params, err := c.RPC.ConsensusParams(ctx, &height)
		if err != nil {
			return err
		}
		return fn(*params)
	})
}

// BlockCallback execute the callback for each new block.
func (c Client) BlockCallback(ctx context.Context, fn func(height int64) error) {
	Callback(ctx, tickerTime, func() error {
		if err := c.WaitForNextBlock(ctx); err != nil {
			return err
		}
		height, err := c.LatestBlockHeight(ctx)
		if err != nil {
			return err
		}
		return fn(height)
	})
}

// Callback execute the callback for each time duration.
func Callback(ctx context.Context, d time.Duration, fn func() error) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := fn(); err != nil {
					// TODO find a better way to send logs
					// fmt.Println(err)
					return
				}
			}
		}
	}()
}
