package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gdamore/tcell/v2"
	"github.com/krallistic/kazoo-go"
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
	kafkaCli         sarama.Client
	zooClient        *kazoo.Kazoo
	enterPanel       *tview.Form
	versionPanel     *tview.Flex
	versionInfoPanel *tview.TextView

	metaPanel    *tview.TextView
	commandPanel *tview.TextView
	resultsPanel *tview.TextView
	outputPanel  *tview.List

	searchPanel *tview.InputField
	topicsPanel *tview.List
	infoPanel   *tview.TextView
	helpPanel   *tview.Flex
	hintPanel   *tview.Modal

	leftPanel  *tview.Flex
	rightPanel *tview.Flex
	layout     *tview.Flex
	pages      *tview.Pages
	app        *tview.Application

	version   string
	gitCommit string
	config    Config
}

func NewKafkaTUI() *KafkaTUI {
	tui := &KafkaTUI{}
	tui.app = &tview.Application{}
	tui.enterPanel = tui.CreateEnterPanel()
	tui.versionPanel = tui.CreateVersionPanel()

	tui.searchPanel = tui.CreateSearchPanel()
	tui.topicsPanel = tui.CreateTopicsPanel()
	tui.infoPanel = tui.CreateInfoPanel()

	tui.leftPanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.searchPanel, 3, 0, false).
		AddItem(tui.topicsPanel, 0, 1, false).
		AddItem(tui.infoPanel, 3, 1, false)

	tui.versionPanel = tui.CreateVersionPanel()
	tui.metaPanel = tui.CreateMetaPanel()
	tui.commandPanel = tui.CreateCommandPanel()
	tui.resultsPanel = tui.CreateResultsPanel()
	tui.outputPanel = tui.CreateOutputPanel()

	tui.rightPanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.versionPanel, 3, 0, false).
		AddItem(tui.metaPanel, 0, 1, false).
		AddItem(tui.commandPanel, 3, 1, false).
		AddItem(tui.resultsPanel, 3, 1, false).
		AddItem(tui.outputPanel, 3, 1, false)

	tui.layout = tview.NewFlex().
		AddItem(tui.leftPanel, 0, 3, false).
		AddItem(tui.rightPanel, 0, 8, false)
	tui.pages = tview.NewPages()
	tui.pages.AddPage("base", tui.layout, true, true)

	return tui
}

func (k *KafkaTUI) CreateVersionPanel() *tview.Flex {
	versionPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	versionPanel.SetBorder(true).SetTitle(fmt.Sprintf(" Version: %s (%s) ", k.version, k.gitCommit))
	k.versionInfoPanel = tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	versionPanel.AddItem(k.versionInfoPanel, 2, 1, false)
	return versionPanel
}

func (k *KafkaTUI) CreateMetaPanel() *tview.TextView {
	return nil
}

func (k *KafkaTUI) CreateResultsPanel() *tview.TextView {
	return nil
}

func (k *KafkaTUI) CreateCommandPanel() *tview.TextView {
	return nil
}

func (k *KafkaTUI) CreateOutputPanel() *tview.List {
	return nil
}

func (k *KafkaTUI) CreateSearchPanel() *tview.InputField {
	searchArea := tview.NewInputField().SetLabel(" TOPIC ")
	searchArea.SetBorder(true).SetTitle(" Search (%s) ")
	return searchArea
}

func (k *KafkaTUI) CreateTopicsPanel() *tview.List {
	topicList := tview.NewList().ShowSecondaryText(false)
	topicList.SetBorder(true).SetTitle(" Keys (%s) ")
	return topicList
}

func (k *KafkaTUI) CreateInfoPanel() *tview.TextView {
	infoArea := tview.NewTextView()
	return infoArea
}

func (k *KafkaTUI) CreateWelcomePanel() {
	welcomeScreen := tview.NewTextView().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(welcomeScreen, true).Run(); err != nil {
		panic(err)
	}
}

func (k *KafkaTUI) CreateEnterPanel() *tview.Form {
	form := tview.NewForm().
		AddDropDown("SecurityProtocol", []string{"SASL/PLAINTEXT"}, 0, nil).
		AddInputField("BrokerList", "", 20, nil, nil).
		AddInputField("UserName", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Enter", nil).
		AddButton("Quit", func() {
			k.app.Stop()
		}).SetButtonBackgroundColor(tcell.ColorRed)
	form.SetBorder(true).SetTitle("Enter KafkaConn Info").SetTitleAlign(tview.AlignLeft)
	return form
}

func (k *KafkaTUI) CreateHintPanel() *tview.Modal {
	return nil
}

func (k *KafkaTUI) Start() error {
	return k.app.SetRoot(k.layout, false).Run()
}

var outputMsgs = make(chan []OutputMsg, 0)

type OutputMsg struct {
	Color   tcell.Color
	Message string
}
