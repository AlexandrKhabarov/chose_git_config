package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

type ConsoleUI struct {
	blocks             []*Block
	selectedBlockIndex int
	namesBlockIndex    int
	emailsBlockIndex   int
	selectionColor     termbox.Attribute
	selectedBlockColor termbox.Attribute
}

func NewConsoleUI(names, email chan []byte) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	w, h := termbox.Size()

	w, h = w/3, h/5

	namesRows := make([][]byte, 0, 0)
	namesBlock := Block{
		height:           h,
		width:            w,
		startXPos:        w,
		startYPos:        h,
		rows:             &namesRows,
		currentRowIndex:  0,
		isRowSelected:    false,
		blockColor:       termbox.ColorGreen,
		switchRowColor:   termbox.ColorBlue,
		selectedRowColor: termbox.ColorCyan,
	}

	emailRows := make([][]byte, 0, 0)
	emailsBlock := Block{
		height:           h,
		width:            w,
		startXPos:        w,
		startYPos:        h * 3,
		rows:             &emailRows,
		currentRowIndex:  0,
		isRowSelected:    false,
		blockColor:       termbox.ColorGreen,
		switchRowColor:   termbox.ColorBlue,
		selectedRowColor: termbox.ColorCyan,
	}
	ui := ConsoleUI{
		blocks: []*Block{
			&namesBlock,
			&emailsBlock,
		},
		selectedBlockIndex: 0,
		namesBlockIndex:    0,
		emailsBlockIndex:   1,
		selectionColor:     termbox.ColorRed,
	}
	ui.RunUI(names, email)
}

func (self *ConsoleUI) RunUI(names, email chan []byte) {
	self.renderBlocks()
	self.selectBlock()
	go self.fillNameBlock(names)
	go self.fillEmailBlock(email)
	self.runEventLoop()
}

func (self *ConsoleUI) getNamesBlock() *Block {
	return self.blocks[self.namesBlockIndex]
}

func (self *ConsoleUI) getEmailsBlock() *Block {
	return self.blocks[self.emailsBlockIndex]
}

func (self *ConsoleUI) renderBlocks() {
	for _, block := range self.blocks {
		block.drawBlock()
	}
}

func (self *ConsoleUI) getSelectedBlock() *Block {
	return self.blocks[self.selectedBlockIndex]
}

func (self *ConsoleUI) changeBlock() {
	selectedBlockIndex := self.selectedBlockIndex
	self.unSelectBlock()
	nextBlockIndex := (selectedBlockIndex + 1) % len(self.blocks)
	self.selectedBlockIndex = nextBlockIndex
	self.selectBlock()
}

func (self *ConsoleUI) unSelectBlock() {
	selectedBlock := self.getSelectedBlock()
	selectedBlock.blockColor = self.selectedBlockColor
	selectedBlock.drawBlock()
}

func (self *ConsoleUI) selectBlock() {
	selectedBlock := self.getSelectedBlock()
	self.selectedBlockColor = selectedBlock.blockColor
	selectedBlock.blockColor = self.selectionColor
	selectedBlock.drawBlock()
}

func (self *ConsoleUI) fillNameBlock(names chan []byte) {
	namesBlock := self.getNamesBlock()
	for name := range names {
		namesBlock.addRow(name)
	}
}

func (self *ConsoleUI) fillEmailBlock(emails chan []byte) {
	emailsBlock := self.getEmailsBlock()
	for email := range emails {
		emailsBlock.addRow(email)
	}
}

func (self *ConsoleUI) runEventLoop() {
loop:
	for {
		switch e := termbox.PollEvent(); e.Key {
		case termbox.KeyEsc:
			break loop
		case termbox.KeyArrowDown:
			block := self.getSelectedBlock()
			block.handleArrowDown()
		case termbox.KeyArrowUp:
			block := self.getSelectedBlock()
			block.handleArrowUp()
		case termbox.KeySpace:
			block := self.getSelectedBlock()
			block.handleBackSpace()
		case termbox.KeyTab:
			self.changeBlock()
		case termbox.KeyEnter:
			namesBlock := self.getNamesBlock()
			emailsBlock := self.getEmailsBlock()
			name := namesBlock.getSelectedRow()
			email := emailsBlock.getSelectedRow()
			user := User{
				UserName:  name,
				UserEmail: email,
			}
			//TODO: handle errors and imagine another way for passing path for .git/config
			UpdateUserInfo(".git/config", []byte(user.UserRepresentation()))
			break loop
		default:
			time.Sleep(time.Microsecond * 10)
		}
	}
}
