package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"ivs-net-winclinet/configure"
	"ivs-net-winclinet/login"
	"ivs-net-winclinet/netLink"
	"os"
	"sync"
)

type App struct {
	app     fyne.App
	output  *widget.Label
	email   *widget.Entry
	pass    *widget.Entry
	code    *widget.Entry
	setting *widget.Button
	win     fyne.Window
	ui      fyne.Window
	set     fyne.Window
}

var my App

// UI 图形界面
func UI() {
	os.Setenv("FYNE_FONT", "conf/simkai.ttf") //设置env环境
	a := app.New()
	my.app = a
	//
	/*loginUi := a.NewWindow("VPN登录")
	ui := a.NewWindow("VPN工具")
	loginUi.Resize(fyne.NewSize(430, 300))
	logo, _ := fyne.LoadResourceFromPath("conf/link.png")
	loginUi.SetIcon(logo)
	loginUi.CenterOnScreen()
	loginUi.SetFixedSize(true)
	loginUi.SetMaster()
	loginUi.Show()

	a.Run()*/
	//
	w := a.NewWindow("VPN登录")
	my.ui = a.NewWindow("VPN工具")
	my.set = a.NewWindow("VPN设置")
	my.win = w

	label, email, pass, submit, setting := my.makeUI()
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), label, layout.NewSpacer())
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(430, 60)),
		centered, email, pass, submit, setting)
	//w.SetContent(container.NewVBox(label, email, pass, submit, setting))
	w.SetContent(grid)
	makeTray(a)
	w.Resize(fyne.NewSize(430, 300))
	logo, _ := fyne.LoadResourceFromPath("conf/link.png")
	w.SetIcon(logo)
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.SetMaster()
	w.Show()
	a.Run()
	os.Unsetenv("FYNE_FONT")

}

type Input struct {
	email *widget.Entry
	pass  *widget.Entry
}

var input Input

func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Entry, *widget.Button, *widget.Button) {
	label := widget.NewLabel("欢迎使用VPN工具")
	input.email = widget.NewEntry()
	input.email.SetPlaceHolder("邮箱")
	input.pass = widget.NewEntry()
	input.pass.SetPlaceHolder("密码")

	submit := widget.NewButton("登录", submit)

	setting := widget.NewButton("设置", setting)
	logo, _ := fyne.LoadResourceFromPath("conf/setting.png")
	setting.SetIcon(logo)

	app.output = label
	return label, input.email, input.pass, submit, setting
}

var uiFlag bool

var wg sync.WaitGroup

func makeMasterUi() {

	my.win.Hide()
	uiFlag = true
	w := my.ui
	//页面内容
	label := widget.NewLabel("VPN连接")
	label.Resize(fyne.NewSize(830, 100))
	labelCenter := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), label, layout.NewSpacer())

	linked := widget.NewLabel("未选择连接")
	linkedCenter := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), linked, layout.NewSpacer())

	flag, body := login.AddressIp(UserEmail, 0)
	address := netLink.Long2IPString(body.Ip)

	var options []string
	if flag {
		options = []string{address}
	} else {
		options = []string{"无ip"}
	}
	newSelect := widget.NewSelect(options, func(s string) {
	})
	newSelect.SetSelected(options[0])
	newSelect.Resize(fyne.NewSize(830, 200))
	//selectCenter := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), newSelect, layout.NewSpacer())

	refresh := widget.NewButton("刷新", func() {
		flag, body = login.AddressIp(UserEmail, 0)
		if flag {
			address = netLink.Long2IPString(body.Ip)
			newSelect.Options = []string{address}
			newSelect.ClearSelected()
			newSelect.SetSelectedIndex(0)
		} else {
			linked.SetText("请刷新重试：" + body.Mes)
			newSelect.Options = []string{body.Mes}
			newSelect.ClearSelected()
			newSelect.SetSelectedIndex(0)
		}
	})
	logo, _ := fyne.LoadResourceFromPath("conf/refresh.png")
	refresh.SetIcon(logo)
	refreshCenter := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), refresh, layout.NewSpacer())
	refresh.Hide()
	if !flag {
		refresh.Show()
		newSelect.Options = []string{"无ip"}
		newSelect.SetSelected(options[0])
	}
	isLink := false
	var button *widget.Button
	wg.Add(2)
	button = widget.NewButton(" 连  接 ", func() {
		if flag {
			if isLink {
				netLink.UnNet()
				linked.SetText("未选择连接")
				button.SetText(" 连  接 ")
				refresh.Show()
				isLink = false
			} else {
				ip := body.Ip
				flag, body = login.Dhcp(UserEmail, ip)
				netLink.Em = UserEmail
				netLink.IIp = ip
				linked.SetText("已连接：" + newSelect.Selected)
				button.SetText("断开连接")
				refresh.Hide()
				go netLink.Net(body)
				isLink = true
			}
		} else {
			linked.SetText("请刷新重试：" + body.Mes)
			newSelect.Options = []string{"无ip"}
			newSelect.ClearSelected()
			newSelect.SetSelectedIndex(0)
			refresh.Show()
		}

	})
	button.Importance = widget.HighImportance
	button.Resize(fyne.NewSize(530, 200))
	buttonCenter := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), button, layout.NewSpacer())

	newLabel := widget.NewLabel("")

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(830, 50)),
		newLabel, labelCenter, newLabel, newSelect, newLabel, linkedCenter, newLabel, buttonCenter, newLabel, refreshCenter)
	w.SetContent(grid)
	//w.SetContent(container.NewVBox(labelCenter, selectCenter, buttonCenter))

	w.Resize(fyne.NewSize(830, 600))
	linkLogo, _ := fyne.LoadResourceFromPath("conf/link.png")
	w.SetIcon(linkLogo)
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.SetCloseIntercept(func() {
		newForm := dialog.NewConfirm("是否退出主程序?", "取消即为最小化", func(b bool) {
			if b {
				w.SetMaster()
				w.Close()
				netLink.UnNet()
			} else {
				w.Hide()
			}
		}, w)
		newForm.Show()
	})
	w.Show()
	wg.Wait()
}
func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("显示主界面", func() {
		})
		menu := fyne.NewMenu("hello world", h)
		h.Action = func() {
			if uiFlag {
				my.ui.Show()
			}
		}
		logo, _ := fyne.LoadResourceFromPath("conf/link.png")
		desk.SetSystemTrayIcon(logo)
		desk.SetSystemTrayMenu(menu)
	}

}

type config struct {
	EditWidget  *widget.Entry
	CurrentFile fyne.URI
}

var cfg config

func setting() {
	my.win.Hide()
	w := my.set
	//设置页面内容
	cfg.EditWidget = widget.NewMultiLineEntry()
	cfg.EditWidget.Resize(fyne.NewSize(430, 250))
	//菜单栏
	openPrivate := fyne.NewMenuItem("修改私钥", open("privateKey.txt"))
	openPublic := fyne.NewMenuItem("修改公钥", open("publicKey.txt"))
	openConn := fyne.NewMenuItem("修改服务器地址", editConn())
	fileMenu := fyne.NewMenu("设置", openPrivate, openPublic, openConn)
	menu := fyne.NewMainMenu(fileMenu)
	w.SetMainMenu(menu)

	fmt.Println("设置中")

	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(430, 250)), cfg.EditWidget)
	w.SetContent(grid)

	w.Resize(fyne.NewSize(430, 300))
	logo, _ := fyne.LoadResourceFromPath("conf/link.png")
	w.SetIcon(logo)
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.SetCloseIntercept(func() {
		newForm := dialog.NewConfirm("是否保存设置?", "", func(b bool) {
			if b {
				save(w)
			}
			w.Hide()
			my.win.Show()
		}, w)
		newForm.Show()
	})
	w.Show()
}
func editConn() func() {
	return func() {
		data, _ := configure.Config.Get("servers.auth_address").(string)
		cfg.CurrentFile = storage.NewFileURI("./conf/client.toml")
		cfg.EditWidget.SetText(data)
	}
}
func open(u string) func() {
	return func() {
		// 读取文件内容
		data, err := ioutil.ReadFile("./conf/" + u)
		if err != nil {
			fmt.Println("读取文件失败:", err)
			return
		}
		cfg.CurrentFile = storage.NewFileURI("./conf/" + u)
		cfg.EditWidget.SetText(string(data))
	}
}
func save(win fyne.Window) {

	if cfg.CurrentFile != nil {
		writer, err := storage.Writer(cfg.CurrentFile)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		var c string
		if cfg.CurrentFile.Path() == "./conf/client.toml" {
			c = "client_name = \"client conf\"\n\n[servers]\nauth_address = \"" +
				cfg.EditWidget.Text + "\"\nauth_port = \"8080\"\n\ntcp_address = \"" +
				cfg.EditWidget.Text + "\"\ntcp_port = \"9000\""
		} else {
			c = cfg.EditWidget.Text
		}

		writer.Write([]byte(c))
		defer writer.Close()
	}
	configure.ReInit()
	login.GetUrl()
}

/*

var filter = storage.NewExtensionFileFilter([]string{".txt", ".TXT"})

func open(win fyne.Window, u string) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if read == nil {
				return
			}

			defer read.Close()
			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			cfg.EditWidget.SetText(string(data))
			cfg.CurrentFile = read.URI()

			win.SetTitle(win.Title() + " - " + read.URI().Name())
		}, win)
		uri, _ := storage.ListerForURI(storage.NewFileURI("./conf/" + u))
		openDialog.SetLocation(uri)
		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}*/

var UserEmail string

func submit() {

	str1 := login.Login1(input.email.Text)
	//dhcp获得用户名
	UserEmail = input.email.Text
	if len(str1) != 0 {
		my.output.SetText(str1)
	} else {
		str2 := login.Login2(input.pass.Text)
		if len(str2) != 0 {
			my.output.SetText(str2)
		} else {
			entry := widget.NewEntry()
			code := widget.NewFormItem(entry.Text, entry)
			form := widget.NewForm(code)
			newForm := dialog.NewForm("请输入验证码", "登录", "取消", form.Items, func(b bool) {
				if b {
					str3 := login.Login3(entry.Text)
					if str3 == "登录成功！" {
						makeMasterUi()
					}
				}
			}, my.win)
			newForm.Show()
		}
	}
}
