package cli

import (
	"bytes"

	"github.com/nsf/termbox-go"
)

type Block struct {
	width     int
	height    int
	startXPos int
	startYPos int
	// TODO: imanagine way for craeting set for rows (delete duplicates)
	rows             *[][]byte
	currentRowIndex  int
	selectedRowIndex int
	isRowSelected    bool
	blockColor       termbox.Attribute
	selectedRowColor termbox.Attribute
	switchRowColor   termbox.Attribute
}

func (self *Block) drawBlock() {
	self.colorBlock(self.blockColor)
}

func (self *Block) colorBlock(color termbox.Attribute) {
	boxHeight := self.height
	for h := 0; h < boxHeight; h++ {
		self.colorRow(color, h)
	}
	if self.isRowSelected {
		self.colorCurrentRow(self.selectedRowColor)
	} else {
		self.colorCurrentRow(self.switchRowColor)
	}
	termbox.Flush()
}

func (self *Block) handleArrowUp() {
	if self.currentRowIndex-1 >= 0 && !self.isRowSelected {
		self.colorCurrentRow(self.blockColor)
		self.currentRowIndex -= 1
		self.colorCurrentRow(self.switchRowColor)
	}
}

func (self *Block) handleArrowDown() {
	rows := self.rows
	if self.currentRowIndex+1 < len(*rows) && !self.isRowSelected {
		self.colorCurrentRow(self.blockColor)
		self.currentRowIndex += 1
		self.colorCurrentRow(self.switchRowColor)
	}
}
func (self *Block) colorCurrentRow(color termbox.Attribute) {
	self.colorRow(color, self.currentRowIndex)
}

func (self *Block) colorRow(color termbox.Attribute, colorNum int) {
	w, _ := termbox.Size()
	cellBuffer := termbox.CellBuffer()
	start := (self.startYPos+colorNum)*w + self.startXPos
	for i := 0; i < self.width; i++ {
		cell := cellBuffer[start+i]
		cell.Bg = color
		cellBuffer[start+i] = cell
	}
	termbox.Flush()
}

func (self *Block) handleBackSpace() {
	if !self.isRowSelected {
		self.isRowSelected = true
		self.colorCurrentRow(self.selectedRowColor)
		self.selectedRowIndex = self.currentRowIndex
	} else {
		self.isRowSelected = false
		self.colorCurrentRow(self.switchRowColor)
		self.selectedRowIndex = 0
	}
}

func (self *Block) getSelectedRow() []byte {
	rows := self.rows
	return (*rows)[self.selectedRowIndex]
}

// TODO: implement in more ellegant way filtering of incoming rows
// TODO: reimplement method for coloring row
// TODO: add scroling for overloading of rows in block
func (self *Block) addRow(row []byte) {
	rows := self.rows
	rowsLen := len(*rows)
	alreadyExist := false
	if len(row) > 0 {
		for _, r := range *rows {
			if bytes.Compare(r, row) == 0 {
				alreadyExist = true
			}
		}
		if !alreadyExist {
			(*rows) = append((*rows), row)
		}
	}
	if rowsLen < self.height && !alreadyExist && len(row) > 0 {
		cellBuffer := termbox.CellBuffer()
		w, _ := termbox.Size()
		start := (self.startYPos+rowsLen)*w + self.startXPos
		for i, ch := range string(row) {
			cell := cellBuffer[start+i]
			cell.Ch = ch
			cell.Fg = termbox.ColorBlack
			cellBuffer[start+i] = cell
		}
		termbox.Flush()
	}
}
