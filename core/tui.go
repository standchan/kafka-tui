package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
	tview.TextView
	tview.Flex
	tview.List
	tview.InputField
	tview.Primitive
	这几种的区别在哪里
*/

type KafkaTUI struct {
	enterPanel *tview.Form
	// versionPanel *tview.Flex
	// metaPanel    *tview.TextView
	// commandPanel *tview.InputField
	// resultsPanel *tview.TextView
	// outputPanel  *tview.List
	// searchPanel  *tview.InputField
	// topicsPanel  *tview.List
	// infoPanel    *tview.TextArea
	// helpPanel    *tview.Flex

	// version   string
	// gitCommit string
	// config    config.Config

	app *tview.Application
}

//kafClient KafClient, conf config.Config, version string, gitCommit string

func NewKafkaTUI() *KafkaTUI {
	tui := &KafkaTUI{}
	tui.app = &tview.Application{}
	tui.enterPanel = tui.CreateEnterPanel()
	return tui
}

func (ui *KafkaTUI) CreateVersionPanel() *tview.Flex {
	return nil
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
	return nil
}

func (ui *KafkaTUI) CreateTopicsPanel() *tview.List {
	return nil
}

func (ui *KafkaTUI) CreateInfoPanel() *tview.TextArea {
	return nil
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
