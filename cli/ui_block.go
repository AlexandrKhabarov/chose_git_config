package cli


import (
	"github.com/nsf/termbox-go"
)

type Block struct {
	width int
	height int
	startXPos int
	startYPos int
	rows *[][]byte
	currentRowIndex int
	selectedRowIndex int
	blockColor termbox.Attribute
	selectedRowColor termbox.Attribute
	switchRowColor termbox.Attribute

}


func (self *Block) drawBlock() {
	self.colorBlock(self.blockColor)
}

func (self *Block)colorBlock(color termbox.Attribute) {
	boxHeight := self.height
	for h := 0; h <  boxHeight; h++ {
		self.colorRow(color, h)
	}
	self.colorCurrentRow(self.switchRowColor)
	termbox.Flush()
}

func (self *Block) handleArrowUp() {
	if self.currentRowIndex - 1 >= 0 {
		self.colorCurrentRow(self.blockColor)
		self.currentRowIndex -= 1
		self.colorCurrentRow(self.switchRowColor)
	}
}

func (self *Block) handleArrowDown() {
	rows := self.rows
	if self.currentRowIndex + 1 < len(*rows) {
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
	start := (self.startYPos + colorNum) * w + self.startXPos
	for i := 0; i < self.width; i++ {
		cell := cellBuffer[start + i]
		cell.Bg = color
		cellBuffer[start + i] = cell
	}
	termbox.Flush()
}

func (self *Block) handleBackSpace() {
	// todo: realize handle backspace button
	
}

func (self *Block) getSelectedRow() []byte {
	rows := self.rows
	return (*rows)[self.selectedRowIndex]
}

func (self *Block) addRow(row []byte) {
	rows := self.rows
	rowsLen := len(*rows)
	if len(row) > 0 {
		(*rows) = append((*rows), row)
	}
	if rowsLen < self.height {
		cellBuffer := termbox.CellBuffer()
		w, _ := termbox.Size()
		start := (self.startYPos + rowsLen - 1) * w + self.startXPos
		for i, ch := range string(row) {
			cell := cellBuffer[start + i]
			cell.Ch = ch
			cell.Fg = termbox.ColorBlack
			cellBuffer[start + i] = cell
		}
	}
	termbox.Flush()
} 