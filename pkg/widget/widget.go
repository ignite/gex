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
	"context"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/text"
)

const loading = "⌛ loading..."

// Widget widget container holder.
type Widget struct {
	terminal          *termbox.Terminal
	container         *container.Container
	currentNetwork    *text.Text
	health            *text.Text
	time              *text.Text
	peers             *text.Text
	secondsPerBlock   *text.Text
	maxBlockSize      *text.Text
	validators        *text.Text
	gasMax            *text.Text
	gasAvgBlock       *text.Text
	gasAvgTransaction *text.Text
	latestGas         *text.Text
	transactions      *text.Text
	blocks            *text.Text
	moniker           *text.Text
	blockProgress     *donut.Donut
}

// New initialize widgets.
func New() (*Widget, error) {
	var (
		widget = new(Widget)
		err    error
	)

	// Creates Network Widget.
	if widget.currentNetwork, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.currentNetwork.Write(loading); err != nil {
		return widget, err
	}

	// Creates Health Widget.
	if widget.health, err = text.New(); err != nil {
		return widget, err
	}
	if err := widget.health.Write(loading); err != nil {
		return widget, err
	}

	// Creates System Time Widget.
	if widget.time, err = text.New(); err != nil {
		return widget, err
	}
	currentTime := time.Now()
	if err := widget.time.Write(currentTime.Format("2006-01-02\n15:04:05")); err != nil {
		return widget, err
	}

	// Creates Connected Peers Widget.
	if widget.peers, err = text.New(); err != nil {
		return widget, err
	}
	if err := widget.peers.Write("0"); err != nil {
		return widget, err
	}

	// Creates Seconds Between Blocks Widget.
	if widget.secondsPerBlock, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.secondsPerBlock.Write("0"); err != nil {
		return widget, err
	}

	// Creates Max Block Size Widget.
	if widget.maxBlockSize, err = text.New(); err != nil {
		return widget, err
	}
	if err := widget.maxBlockSize.Write("0"); err != nil {
		return widget, err
	}

	// Creates Validators widget.
	if widget.validators, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.validators.Write(loading); err != nil {
		return widget, err
	}

	// Creates Validators widget.
	if widget.gasMax, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.gasMax.Write(loading); err != nil {
		return widget, err
	}

	// Creates Gas per Average Block Widget.
	if widget.gasAvgBlock, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.gasAvgBlock.Write(loading); err != nil {
		return widget, err
	}

	// Creates Gas per Average Transaction Widget.
	if widget.gasAvgTransaction, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.gasAvgTransaction.Write(loading); err != nil {
		return widget, err
	}

	// Creates Gas per Latest Transaction Widget.
	if widget.latestGas, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.latestGas.Write(loading); err != nil {
		return widget, err
	}

	// Add big widgets.

	// Block Status Donut widget.
	if widget.blockProgress, err = donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorGreen)),
		donut.Label("New Block Status", cell.FgColor(cell.ColorGreen)),
	); err != nil {
		return widget, err
	}

	// Transaction parsing widget.
	if widget.transactions, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.transactions.Write("Transactions will appear as soon as they are confirmed in a moniker."); err != nil {
		return widget, err
	}

	// Blocks parsing widget.
	if widget.blocks, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}

	// Create Blocks parsing widget.
	if widget.moniker, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return widget, err
	}
	if err := widget.moniker.Write(loading); err != nil {
		return widget, err
	}

	if widget.terminal, err = termbox.New(); err != nil {
		return widget, err
	}

	if widget.container, err = widget.drawView(); err != nil {
		return widget, err
	}

	return widget, nil
}

// Cleanup widgets.
func (w *Widget) Cleanup() {
	w.terminal.Close()
}

// Run the widget view.
func (w *Widget) Run(ctx context.Context) error {
	defer w.Cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' || k.Key == keyboard.KeyEsc {
			cancel()
		}
	}
	return termdash.Run(ctx, w.terminal, w.container, termdash.KeyboardSubscriber(quitter))
}
