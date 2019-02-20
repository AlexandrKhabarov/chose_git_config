package cli

import (
	"fmt"
	"time"
	"github.com/nsf/termbox-go"
)


type ConsoleUI struct {
	finishChan chan struct{}
	arrowDownChan chan struct{}
	arrowUpChan chan struct{}
	selectChan chan struct{}
	changeWindow chan struct{}
	completeChanges chan struct{}
	blocks []Block
	selectedBlockIndex int
	namesBlockIndex int
	emailsBlockIndex int
	selectionColor termbox.Attribute
}

func NewConsoleUI(names, email chan []byte) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	w, h := termbox.Size()

	w, h = w/3, h/5
	
	namesBlock := Block{
		height: h,
		width: w,
		startXPos: w,
		startYPos: h,
		rows: make([][]byte, 0, 0),
		currentRowIndex: 0,
		blockColor: termbox.ColorGreen,
		switchRowsColor: termbox.ColorBlue,
	}
	
	emailsBlock := Block{
		height: h,
		width: w,
		startXPos: w,
		startYPos: h*3,
		rows: make([][]byte, 0, 0),
		currentRowIndex: 0,
		blockColor: termbox.ColorGreen,
		switchRowsColor: termbox.ColorBlue,
	}
	ui := ConsoleUI {
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		make(chan struct{}, 0),
		[]Block {
			namesBlock,
			emailsBlock,
		},
		0,
		0,
		1,
		termbox.ColorRed,
	}
	ui.RunUI(names, email)
}


func(self *ConsoleUI) RunUI(names, email chan []byte) {
	go self.runKeyBoardEventHandler()
	self.renderBlocks()
	self.selectBlock()
	self.fillNameBlock(names)
	fmt.Println(len(self.getNamesBlock().rows))
	self.fillEmailBlock(email)
	fmt.Println(len(self.getEmailsBlock().rows))
	self.runEventLoop()
}

func (self *ConsoleUI) getNamesBlock() Block {
	return self.blocks[self.namesBlockIndex]
}

func (self *ConsoleUI) setNamesBlock(block Block) {
	self.blocks[self.namesBlockIndex] = block
}

func (self *ConsoleUI) getEmailsBlock() Block {
	return self.blocks[self.emailsBlockIndex]
}

func (self *ConsoleUI) setEmailsBlock(block Block) {
	self.blocks[self.emailsBlockIndex] = block
}

func (self *ConsoleUI) renderBlocks() {
	for _, block := range self.blocks {
		block.renderDefaultBlock()
	}
}

func (self *ConsoleUI) getSelectedBlock() Block {
	return self.blocks[self.selectedBlockIndex]
}

func (self *ConsoleUI) setSelectedBlock(block Block){
	self.blocks[self.selectedBlockIndex] = block
}

func (self *ConsoleUI) changeBlock() {
	selectedBlockIndex := self.selectedBlockIndex
	selectedBlock := self.getSelectedBlock()
	selectedBlock.renderDefaultBlock()
	nextBlockIndex := (selectedBlockIndex + 1) % len(self.blocks)
	self.selectedBlockIndex = nextBlockIndex
	self.selectBlock()
}

func (self *ConsoleUI) selectBlock() {
	selectedBlock := self.getSelectedBlock()
	selectedBlock.changeBackGraoundColor(self.selectionColor)
}

func(self *ConsoleUI) runKeyBoardEventHandler() {
	for {
		e := termbox.PollEvent()
		switch e.Key {
		case termbox.KeyEsc:
			self.finishChan <- struct{}{}
		case termbox.KeyArrowDown:
			self.arrowDownChan <- struct{}{}
		case termbox.KeyArrowUp:
			self.arrowUpChan <- struct{}{}
		case termbox.KeyBackspace:
			self.selectChan <- struct{}{}
		case termbox.KeyTab:
			self.changeWindow <- struct{}{}
		case termbox.KeyEnter:
			self.completeChanges <- struct{}{}
		}
	}
}

func(self *ConsoleUI) fillNameBlock(names chan []byte) {
	namesBlock := self.getNamesBlock()
	for name := range names {
		namesBlock.addRow(name)
		fmt.Println(len(self.getNamesBlock().rows))
	}
	self.setNamesBlock(namesBlock)
}

func(self *ConsoleUI) fillEmailBlock(emails chan []byte) {
	emailsBlock := self.getEmailsBlock()
	for email := range emails {
		emailsBlock.addRow(email)
	}
	self.setEmailsBlock(emailsBlock)
}

func (self *ConsoleUI) runEventLoop() {
loop:
	for {
		select {
		case <- self.finishChan:
			break loop
		case <- self.changeWindow:
			self.changeBlock()
		case <- self.arrowUpChan:
			block := self.getSelectedBlock()
			block.handleArrowUp()
			self.setSelectedBlock(block)
		case <- self.arrowDownChan:
			block := self.getSelectedBlock()
			block.handleArrowDown()
			self.setSelectedBlock(block)
		case <- self.selectChan:
			block := self.getSelectedBlock()
			block.handleBackSpace()
		case <- self.completeChanges:
			// namesBlock := self.getNamesBlock()
			// emailsBlock := self.getNamesBlock()
			// name := namesBlock.getSelectedRow()
			// email := emailsBlock.getSelectedRow()
			// user := User{
			// 	UserName: name,
			// 	UserEmail: email,
			// }
			// todo: handle errors and imagine another way for passing path for .git/config
			// UpdateUserInfo("./.git/config", []byte(user.UserRepresentation()))
			break loop
		default:
			time.Sleep(time.Millisecond*10)
		}
	}
}

