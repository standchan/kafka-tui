package core

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
	tview.Box	盒子、外壳
	tview.TextView	文本输出？
	tview.TextArea	文本输入
	tview.Table	表格
	tview.TreeView 树形文件夹结构
	tview.Grid	网格布局
	tview.Flex	弹性布局
	tview.List	列表
	tview.InputField
	tview.Primitive
	这几种的区别在哪里
*/

type KafkaTUI struct {
	enterPanel       *tview.Form
	versionPanel     *tview.Flex
	versionInfoPanel *tview.TextView

	leftPanel  *tview.Flex
	rightPanel *tview.Flex
	// metaPanel    *tview.TextView
	// commandPanel *tview.InputField
	// resultsPanel *tview.TextView
	// outputPanel  *tview.List
	searchPanel *tview.InputField
	topicsPanel *tview.List
	infoPanel   *tview.TextView
	// helpPanel    *tview.Flex

	version   string
	gitCommit string
	// config    config.Config

	app *tview.Application
}

//kafClient KafClient, conf config.Config, version string, gitCommit string

func NewKafkaTUI() *KafkaTUI {
	tui := &KafkaTUI{}
	tui.app = &tview.Application{}
	tui.enterPanel = tui.CreateEnterPanel()
	tui.versionPanel = tui.CreateVersionPanel()
	tui.leftPanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.searchPanel, 3, 0, false).
		AddItem(tui.topicsPanel, 0, 1, false).
		AddItem(tui.infoPanel, 3, 1, false)
	return tui
}

func (ui *KafkaTUI) CreateVersionPanel() *tview.Flex {
	versionPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	versionPanel.SetBorder(true).SetTitle(fmt.Sprintf(" Version: %s (%s) ", ui.version, ui.gitCommit))
	ui.versionInfoPanel = tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	versionPanel.AddItem(ui.versionInfoPanel, 2, 1, false)
	return versionPanel
}

func (ui *KafkaTUI) CreateMetaPanel() *tview.TextView {
	return nil
}

func (ui *KafkaTUI) CreateResultsPanel() *tview.InputField {
	return nil
}

func (ui *KafkaTUI) CreateCommandPanel() *tview.TextView {
	return nil
}

func (ui *KafkaTUI) CreateOutputPanel() *tview.List {
	return nil
}

func (ui *KafkaTUI) CreateSearchPanel() *tview.InputField {
	searchArea := tview.NewInputField().SetLabel(" TOPIC ")
	searchArea.SetBorder(true).SetTitle(" Search (%s) ")
	return searchArea
}

func (ui *KafkaTUI) CreateTopicsPanel() *tview.List {
	topicList := tview.NewList().ShowSecondaryText(false)
	topicList.SetBorder(true).SetTitle(" Keys (%s) ")
	return topicList
}

func (ui *KafkaTUI) CreateInfoPanel() *tview.TextView {
	infoArea := tview.NewTextView()
	return infoArea
}

func (ui *KafkaTUI) CreateWelcomePanel() {
	welcomeScreen := tview.NewTextView().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(welcomeScreen, true).Run(); err != nil {
		panic(err)
	}
}

func (ui *KafkaTUI) CreateEnterPanel() *tview.Form {
	form := tview.NewForm().
		AddDropDown("SecurityProtocal", []string{"SASL/PLAINTEXT"}, 0, nil).
		AddInputField("BrokerList", "", 20, nil, nil).
		AddInputField("UserName", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Enter", nil).
		AddButton("Quit", func() {
			ui.app.Stop()
		}).SetButtonBackgroundColor(tcell.ColorRed)
	form.SetBorder(true).SetTitle("Enter KafkaConn Info").SetTitleAlign(tview.AlignLeft)
	return form
}

func (ui *KafkaTUI) Start() error {
	return ui.app.SetRoot(ui.enterPanel, true).Run()
}
