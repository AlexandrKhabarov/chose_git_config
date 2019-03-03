package cli

import (
	"github.com/nsf/termbox-go"
	"github.com/sachez/chose_git_config/config"
	"time"
)

type ConsoleUI struct {
	Errors             chan error
	blocks             []*Block
	selectedBlockIndex int
	namesBlockIndex    int
	emailsBlockIndex   int
	selectionColor     termbox.Attribute
	selectedBlockColor termbox.Attribute
}

func NewConsoleUI() ConsoleUI {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

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
		Errors: make(chan error),
		blocks: []*Block{
			&namesBlock,
			&emailsBlock,
		},
		selectedBlockIndex: 0,
		namesBlockIndex:    0,
		emailsBlockIndex:   1,
		selectionColor:     termbox.ColorRed,
	}
	return ui
}

func (ui *ConsoleUI) RunUI(names, email chan []byte) {
	ui.renderBlocks()
	ui.selectBlock()
	go ui.fillNameBlock(names)
	go ui.fillEmailBlock(email)
	ui.runEventLoop()
}

func (ui *ConsoleUI) getNamesBlock() *Block {
	return ui.blocks[ui.namesBlockIndex]
}

func (ui *ConsoleUI) getEmailsBlock() *Block {
	return ui.blocks[ui.emailsBlockIndex]
}

func (ui *ConsoleUI) renderBlocks() {
	for _, block := range ui.blocks {
		block.drawBlock()
	}
}

func (ui *ConsoleUI) getSelectedBlock() *Block {
	return ui.blocks[ui.selectedBlockIndex]
}

func (ui *ConsoleUI) changeBlock() {
	selectedBlockIndex := ui.selectedBlockIndex
	ui.unSelectBlock()
	nextBlockIndex := (selectedBlockIndex + 1) % len(ui.blocks)
	ui.selectedBlockIndex = nextBlockIndex
	ui.selectBlock()
}

func (ui *ConsoleUI) unSelectBlock() {
	selectedBlock := ui.getSelectedBlock()
	selectedBlock.blockColor = ui.selectedBlockColor
	selectedBlock.drawBlock()
}

func (ui *ConsoleUI) selectBlock() {
	selectedBlock := ui.getSelectedBlock()
	ui.selectedBlockColor = selectedBlock.blockColor
	selectedBlock.blockColor = ui.selectionColor
	selectedBlock.drawBlock()
}

func (ui *ConsoleUI) fillNameBlock(names chan []byte) {
	namesBlock := ui.getNamesBlock()
	for name := range names {
		namesBlock.addRow(name)
	}
}

func (ui *ConsoleUI) fillEmailBlock(emails chan []byte) {
	emailsBlock := ui.getEmailsBlock()
	for email := range emails {
		emailsBlock.addRow(email)
	}
}

func (ui *ConsoleUI) runEventLoop() {
loop:
	for {
		switch e := termbox.PollEvent(); e.Key {
		case termbox.KeyEsc:
			termbox.Close()
			break loop
		case termbox.KeyArrowDown:
			block := ui.getSelectedBlock()
			block.handleArrowDown()
		case termbox.KeyArrowUp:
			block := ui.getSelectedBlock()
			block.handleArrowUp()
		case termbox.KeySpace:
			block := ui.getSelectedBlock()
			block.handleBackSpace()
		case termbox.KeyTab:
			ui.changeBlock()
		case termbox.KeyEnter:
			namesBlock := ui.getNamesBlock()
			emailsBlock := ui.getEmailsBlock()
			name := namesBlock.getSelectedRow()
			email := emailsBlock.getSelectedRow()
			user := config.User{
				UserName:  name,
				UserEmail: email,
			}
			//TODO: handle errors and imagine another way for passing path for .git/config
			config.UpdateUserInfo(".git/config", []byte(user.UserRepresentation()))
			break loop
		default:
			time.Sleep(time.Microsecond * 10)
		}
	}
}
