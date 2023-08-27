package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func main() {
	// 创建一个新的应用程序
	app := tview.NewApplication()

	// 创建登陆界面的组件
	usernameInput := tview.NewInputField().SetLabel("Username: ").SetFieldWidth(20)
	passwordInput := tview.NewInputField().SetLabel("Password: ").SetFieldWidth(20).SetMaskCharacter('*')
	loginButton := tview.NewButton("Login").SetSelectedFunc(func() {
		// 获取输入的用户名和密码
		usernameInput.GetText()
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
						newPage.SetText(fmt.Sprintf("Welcome to the new page!\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nCurrent Time: %s", time.Now().Format("15:04:05")))
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
		}
	})

	// 创建布局，并将组件添加到布局中
	loginLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Welcome to the Login Page").SetTextAlign(tview.AlignCenter), 0, 1, false).
		AddItem(usernameInput, 0, 1, false).
		AddItem(passwordInput, 0, 1, false).
		AddItem(loginButton, 0, 1, false)

	// 将登陆页面布局设置为根组件，并启动应用程序
	if err := app.SetRoot(loginLayout, true).Run(); err != nil {
		fmt.Println(err)
	}
}
