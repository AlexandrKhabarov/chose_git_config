package cli


import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type Block struct {
	width int
	height int
	startXPos int
	startYPos int
	rows [][]byte
	currentRowIndex int
	selectedRowIndex int
	blockColor termbox.Attribute
	selectedRowColor termbox.Attribute
	switchRowsColor termbox.Attribute

}


func (self *Block) renderDefaultBlock() {
	cellBuffer := termbox.CellBuffer()

	boxStartXPos := self.startXPos
	boxWidth := self.width

	boxStartYPos := self.startYPos
	boxHeight := self.height
	
	w, _ := termbox.Size()

	for h := 0; h <  boxHeight; h++ {
		start := (boxStartYPos + h) * w + boxStartXPos
		for w := 0; w < boxWidth; w++ {
			cell := cellBuffer[start + w]
			cell.Bg = self.blockColor	
			cellBuffer[start + w] = cell
		}
	}
	self.colorToSwitchColor()
	termbox.Flush()
}

func (self *Block)changeBackGraoundColor(color termbox.Attribute) {
	cellBuffer := termbox.CellBuffer()

	boxStartXPos := self.startXPos
	boxWidth := self.width

	boxStartYPos := self.startYPos
	boxHeight := self.height
	
	w, _ := termbox.Size()

	for h := 0; h <  boxHeight; h++ {
		start := (boxStartYPos + h) * w + boxStartXPos
		for w := 0; w < boxWidth; w++ {
			cell := cellBuffer[start + w]
			cell.Bg = color	
			cellBuffer[start + w] = cell
		}
	}
	self.colorToSwitchColor()
	termbox.Flush()
}

func (self *Block) handleArrowUp() {
	if self.currentRowIndex - 1 >= 0 {
		fmt.Println("Arrow up")
		self.colorToDefaultColor()
		self.currentRowIndex -= 1
		self.colorToSwitchColor()
	}
}

func (self *Block) handleArrowDown() {
	if self.currentRowIndex + 1 < len(self.rows) {
		self.colorToDefaultColor()
		self.currentRowIndex += 1
		self.colorToSwitchColor()
	}
}
func (self *Block) colorToDefaultColor() {
	self.colorRow(self.blockColor, self.currentRowIndex)
}

func (self *Block) colorToSwitchColor() {
	self.colorRow(self.switchRowsColor, self.currentRowIndex)
}

func (self *Block) colorRow(color termbox.Attribute, colorNum int) {
	w, _ := termbox.Size()
	cellBuffer := termbox.CellBuffer()
	start := (self.startYPos + colorNum) * w + self.startXPos
	for i := 0; i < self.width; i++ {
		cell := cellBuffer[start + i]
		cell.Bg = color
		cellBuffer[start + i] = cell
	}
	termbox.Flush()
}

func (self *Block) handleBackSpace() {
	
}

func (self *Block) getSelectedRow() []byte {
	return self.rows[self.selectedRowIndex]
}

func (self *Block) addRow(row []byte) {
	if len(row) > 0 {
		self.rows = append(self.rows, row)
	}
	if len(self.rows) < self.height {
		cellBuffer := termbox.CellBuffer()
		w, _ := termbox.Size()
		start := (self.startYPos + len(self.rows) - 1) * w + self.startXPos
		for i, ch := range string(row) {
			cell := cellBuffer[start + i]
			cell.Ch = ch
			cell.Fg = termbox.ColorBlack
			cellBuffer[start + i] = cell
		}
	}
	termbox.Flush()
} 