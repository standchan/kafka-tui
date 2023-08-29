package main

import (
	"fmt"
	"github.com/rivo/tview"
	"time"
)

var app = tview.NewApplication()

func main() {
	CreateLoginPanel()
}

// todo:连接失败等错误可以采用统一样式的model进行展示
func CreateLoginPanel() {
	securityProto := tview.NewDropDown().SetLabel("SecurityProtocol: ").SetOptions([]string{"SASL/PLAIN", "SSL/TLS", "SASL/SCRAM"}, nil)
	brokerInput := tview.NewInputField().SetLabel("Broker: ").SetFieldWidth(20)
	portInput := tview.NewInputField().SetLabel("Port: ").SetFieldWidth(20).SetText("9092")
	//brokers := tview.NewFlex().
	//	SetDirection(tview.FlexRow).
	//	AddItem(brokerInput, 0, 1, false).
	//	AddItem(portInput, 0, 1, false).SetBorder(false)
	userInput := tview.NewInputField().SetLabel("User: ").SetFieldWidth(20)
	passwordInput := tview.NewInputField().SetLabel("Password: ").SetFieldWidth(20).SetMaskCharacter('*')
	tview.NewForm().AddFormItem(brokerInput)
	loginModel := tview.NewForm().
		AddFormItem(securityProto).
		AddFormItem(brokerInput).
		AddFormItem(portInput).
		AddFormItem(userInput).
		AddFormItem(passwordInput).
		AddButton("Enter", func() {
			broker := brokerInput.GetText()
			port := portInput.GetText()
			user := userInput.GetText()
			password := passwordInput.GetText()
			if mockKafka(broker, port, user, password) {
				fmt.Println("Connect success!")
			} else {
				errorMessage := tview.NewModal().
					SetText("Invalid username or password!").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						app.SetFocus(brokerInput)
					})
				app.SetRoot(errorMessage, true).SetFocus(errorMessage)
				app.GetFocus()
			}
		})
	loginFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(loginModel, 0, 1, true)
	loginPage := tview.NewPages().AddPage("login", loginFlex, true, true)
	if err := app.SetRoot(loginPage, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}

func mockKafka(broker, port, user, password string) bool {
	return true
}

func login() {
	// 创建登陆界面的组件
	usernameInput := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	passwordInput := tview.NewInputField().SetLabel("Password: ").SetFieldWidth(20).SetMaskCharacter('*')
	loginButton := tview.NewButton("Login").SetSelectedFunc(func() {
		// 获取输入的用户名和密码
		username := usernameInput.GetText()
		password := passwordInput.GetText()

		// 检查用户名和密码是否正确
		if username == "admin" && password == "123456" {
			// 如果正确，创建新页面的组件并设置为根组件
			newPage := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Welcome to the new page!\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
			app.SetRoot(newPage, true).SetFocus(newPage)

			// 创建新页面中显示系统时间的定时器
			go func() {
				for {
					time.Sleep(time.Second) // 每秒更新一次时间
					app.QueueUpdateDraw(func() {
						newPage.SetText(fmt.Sprintf("Welcome to the new page!\nCurrent Time: %s", time.Now().Format("15:04:05")))
					})
				}
			}()
		} else {
			// 如果不正确，显示错误消息
			errorMessage := tview.NewModal().
				SetText("Invalid username or password!").
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetFocus(usernameInput)
				})
			app.SetRoot(errorMessage, true).SetFocus(errorMessage)
			app.GetFocus()
		}
	})

	// 创建布局，并将组件添加到布局中
	loginLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Welcome to the Login Page").SetTextAlign(tview.AlignCenter), 0, 1, false).
		AddItem(usernameInput, 0, 1, false).
		AddItem(passwordInput, 0, 1, false).
		AddItem(loginButton, 0, 1, false)

	// 将登陆页面布局设置为根组件，并启动应用程序
	if err := app.SetRoot(loginLayout, true).EnableMouse(true).Run(); err != nil {
		fmt.Println(err)
	}
}
