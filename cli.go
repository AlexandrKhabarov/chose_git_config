package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

type ConsoleUI struct {
	finishChan chan struct{}
	arrowDownChan chan struct{}
	arrowUpChan chan struct{}
	selectChan chan struct{}
	changeWindow chan struct{}
	completeChanges chan struct{}
}

func NewConsoleUI() ConsoleUI {
	return ConsoleUI {
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
	}
}

func(ui *ConsoleUI) RunUI() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	ui.runKeyBoardEventHandler()
	ui.renderUI()
}

func(ui *ConsoleUI) runKeyBoardEventHandler() {
	for {
		e := termbox.PollEvent()
		switch e.Key {
		case termbox.KeyEsc:
			ui.finishChan <- struct{}{}
		case termbox.KeyArrowDown:
			ui.arrowDownChan <- struct{}{}
		case termbox.KeyArrowUp:
			ui.arrowUpChan <- struct{}{}
		case termbox.KeyBackspace:
			ui.selectChan <- struct{}{}
		case termbox.KeyTab:
			ui.changeWindow <- struct{}{}
		case termbox.KeyEnter:
			ui.completeChanges <- struct{}{}
		}
	}
}

func(ui *ConsoleUI) renderUI() {
	
}

type CustomCell struct {
	cell termbox.Cell
	x    int
	y    int
}

func InitCli() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	finishChan := make(chan struct{})
	arrowDownChan := make(chan struct{})
	arrowUpChan := make(chan struct{})
	backsapceChan := make(chan struct{})

	blocks := make([][][]CustomCell, 0, 8)

	step := 0
	for i := 1; i < 8; i++ {
		rows := make([][]CustomCell, 0, 30)
		for j := 0; j < 30; j++ {
			row := make([]CustomCell, 0, 5)
			for k := 0; k < 3; k++ {
				cell := termbox.Cell{Ch: rune(' '), Fg: termbox.Attribute(i), Bg: termbox.Attribute(i)}
				customCell := CustomCell{cell: cell, x: j, y: k + step}
				termbox.SetCell(j, k+step, cell.Ch, cell.Fg, cell.Bg)
				row = append(row, customCell)
			}
			rows = append(rows, row)
		}
		blocks = append(blocks, rows)
		step += 5
	}

	termbox.Flush()

	go func(finishChan chan<- struct{}) {
		for {
			e := termbox.PollEvent()
			switch e.Key {
			case termbox.KeyEsc:
				finishChan <- struct{}{}
			case termbox.KeyArrowDown:
				arrowDownChan <- struct{}{}
			case termbox.KeyArrowUp:
				arrowUpChan <- struct{}{}
			case termbox.KeyBackspace:
				backspaceChan <- struct{}{}
			}
		}
	}(finishChan)

loop:
	for i := 0; ; {
		select {
		case <-finishChan:
			break loop
		case <-arrowUpChan:
			if i-1 < 0 {
				continue
			} else {
				previousBlock := blocks[i]
				i -= 1
				nextBlock := blocks[i]
				for _, row := range nextBlock {
					for _, cell := range row {
						termbox.SetCell(cell.x, cell.y, cell.cell.Ch, termbox.ColorDefault, termbox.ColorDefault)
					}
				}

				for _, row := range previousBlock {
					for _, cell := range row {
						termbox.SetCell(cell.x, cell.y, cell.cell.Ch, cell.cell.Fg, cell.cell.Bg)
					}
				}

			}
			termbox.Flush()
		case <-arrowDownChan:
			if i+1 > len(blocks)-1 {
				continue
			} else {
				previousBlock := blocks[i]
				i += 1
				nextBlock := blocks[i]
				for _, row := range nextBlock {
					for _, cell := range row {
						termbox.SetCell(cell.x, cell.y, cell.cell.Ch, termbox.ColorDefault, termbox.ColorDefault)
					}
				}
				for _, row := range previousBlock {
					for _, cell := range row {
						termbox.SetCell(cell.x, cell.y, cell.cell.Ch, cell.cell.Fg, cell.cell.Bg)
					}
				}
			}
			termbox.Flush()
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	close(finishChan)
	close(arrowDownChan)
	close(arrowUpChan)
}
