// Copyright 2018 Google Inc.
// Copyright 2020 Tobias Schwarz
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package widget

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/text"
)

// SetCurrentNetwork reset the widget and set current network text.
func (w *Widget) SetCurrentNetwork(txt string, opts ...text.WriteOption) error {
	w.currentNetwork.Reset()
	return w.currentNetwork.Write(txt, opts...)
}

// SetHealth resets the widget and sets health text.
func (w *Widget) SetHealth(txt string, opts ...text.WriteOption) error {
	w.health.Reset()
	return w.health.Write(txt, opts...)
}

// SetTime resets the widget and sets time text.
func (w *Widget) SetTime(txt string, opts ...text.WriteOption) error {
	w.time.Reset()
	return w.time.Write(txt, opts...)
}

// SetPeers resets the widget and sets peers text.
func (w *Widget) SetPeers(peers int, opts ...text.WriteOption) error {
	w.peers.Reset()
	return w.peers.Write(strconv.Itoa(peers), opts...)
}

// SetSecondsPerBlock resets the widget and sets seconds per block text.
func (w *Widget) SetSecondsPerBlock(txt string, opts ...text.WriteOption) error {
	w.secondsPerBlock.Reset()
	return w.secondsPerBlock.Write(txt, opts...)
}

// SetMaxBlockSize resets the widget and sets max block size text.
func (w *Widget) SetMaxBlockSize(txt string, opts ...text.WriteOption) error {
	w.maxBlockSize.Reset()
	return w.maxBlockSize.Write(txt, opts...)
}

// SetValidators resets the widget and sets validators text.
func (w *Widget) SetValidators(validators int, opts ...text.WriteOption) error {
	w.validators.Reset()
	return w.validators.Write(strconv.Itoa(validators), opts...)
}

// SetGasMax resets the widget and sets gas max text.
func (w *Widget) SetGasMax(txt string, opts ...text.WriteOption) error {
	w.gasMax.Reset()
	return w.gasMax.Write(txt, opts...)
}

// SetGasAvgBlock resets the widget and sets gas average per block text.
func (w *Widget) SetGasAvgBlock(txt string, opts ...text.WriteOption) error {
	w.gasAvgBlock.Reset()
	return w.gasAvgBlock.Write(txt, opts...)
}

// SetGasAvgTransaction resets the widget and sets gas average per transaction text.
func (w *Widget) SetGasAvgTransaction(txt string, opts ...text.WriteOption) error {
	w.gasAvgTransaction.Reset()
	return w.gasAvgTransaction.Write(txt, opts...)
}

// SetLatestGas resets the widget and sets latest gas text.
func (w *Widget) SetLatestGas(txt string, opts ...text.WriteOption) error {
	w.latestGas.Reset()
	return w.latestGas.Write(txt, opts...)
}

// SetMoniker resets the widget and sets moniker text.
func (w *Widget) SetMoniker(text string, opts ...text.WriteOption) error {
	w.moniker.Reset()
	return w.moniker.Write(text, opts...)
}

// AddTransaction adds a new transaction to the widget.
func (w *Widget) AddTransaction(txt string, opts ...text.WriteOption) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	if err := w.transactions.Write(
		fmt.Sprintf("\n\nNew Transaction (%s)\n", now),
		text.WriteCellOpts(cell.Bold(), cell.Inverse()),
	); err != nil {
		return err
	}
	return w.transactions.Write(txt, opts...)
}

// AddBlock adds a new block to the widget.
func (w *Widget) AddBlock(txt string, opts ...text.WriteOption) error {
	return w.blocks.Write(txt+"\n", opts...)
}

// SetBlockProgress sets the progress of the block in the widget.
func (w *Widget) SetBlockProgress(percent int, opts ...donut.Option) error {
	return w.blockProgress.Percent(percent, opts...)
}
