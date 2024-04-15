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
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
)

// drawView draw all containers view.
func (w *Widget) drawView() (*container.Container, error) {
	return container.New(
		w.terminal,
		container.Border(linestyle.Light),
		container.BorderTitle("GEX: PRESS Q or ESC TO QUIT"),
		container.BorderColor(cell.ColorNumber(2)),
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(
						container.SplitHorizontal(
							container.Top(
								container.SplitVertical(
									container.Left(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Network"),
												container.PlaceWidget(w.currentNetwork),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("Moniker"),
												container.PlaceWidget(w.moniker),
											),
										),
									),
									container.Right(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Health"),
												container.PlaceWidget(w.health),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("System Time"),
												container.PlaceWidget(w.time),
											),
										),
									),
								),
							),
							container.Bottom(
								// Insert new bottom rows
								container.SplitVertical(
									container.Left(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Block Time"),
												container.PlaceWidget(w.secondsPerBlock),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("Max Block Size"),
												container.PlaceWidget(w.maxBlockSize),
											),
										),
									),
									container.Right(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Connected Peers"),
												container.PlaceWidget(w.peers),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("Validators"),
												container.PlaceWidget(w.validators),
											),
										),
									),
								),
							),
						),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Current Block Round"),
						container.PlaceWidget(w.blockProgress),
					),
				),
			),
			container.Bottom(
				container.SplitVertical(
					container.Left(
						container.SplitHorizontal(
							container.Top(
								container.SplitVertical(
									container.Left(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Gas Max"),
												container.PlaceWidget(w.gasMax),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("Gas Ø Block"),
												container.PlaceWidget(w.gasAvgBlock),
											),
										),
									),
									container.Right(
										container.SplitVertical(
											container.Left(
												container.Border(linestyle.Light),
												container.BorderTitle("Gas Ø Tx"),
												container.PlaceWidget(w.gasAvgTransaction),
											),
											container.Right(
												container.Border(linestyle.Light),
												container.BorderTitle("Gas Latest Tx"),
												container.PlaceWidget(w.latestGas),
											),
										),
									),
								),
							),
							container.Bottom(
								container.Border(linestyle.Light),
								container.BorderTitle("Latest Blocks"),
								container.PlaceWidget(w.blocks),
							),
						),
					), container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Latest Confirmed Transactions"),
						container.PlaceWidget(w.transactions),
					),
				),
			),
		),
	)
}
