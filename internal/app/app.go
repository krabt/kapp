package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

const (
	defaultSidebarWidth float32 = 120
	minSidebarWidth     float32 = 60
)

type App struct {
	name         string
	count        int
	page         string
	sidebarWidth float32
}

func New() *App {
	return &App{page: "home", sidebarWidth: defaultSidebarWidth}
}

func (a *App) View(c *ui.Context) {
	split := ui.Split(c, &a.sidebarWidth, func() {
		a.sidebarView(c)
	}, func() {
		a.mainView(c)
	}).SplitReset(defaultSidebarWidth).Fill()
	if split.Changed() && a.sidebarWidth < minSidebarWidth {
		a.sidebarWidth = minSidebarWidth
	}
}

func (a *App) sidebarView(c *ui.Context) {
	ui.Sidebar(c, &a.page, func() {
		ui.SidebarSection(c, "工作区", nil, func() {
			ui.SidebarItem(c, "home", nil, "概览")
			ui.SidebarItem(c, "counter", nil, "计数器")
		})
	}).Fill()
}

func (a *App) mainView(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Gap(16).Padding(28).Children(func() {
		switch a.page {
		case "counter":
			ui.Text(c, "计数器").FontSize(26).Bold()
			ui.Text(c, "使用按钮调整当前数值。").TextColor(t.TextMuted)
			ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
				if ui.Button(c, "−").Width(44).Clicked() {
					a.count--
				}
				ui.Textf(c, "%d", a.count).FontSize(20).Width(48).TextAlign(ui.Center)
				if ui.PrimaryButton(c, "+").Width(44).Clicked() {
					a.count++
				}
			})
		default:
			ui.Text(c, "概览").FontSize(26).Bold()
			ui.Text(c, "欢迎使用 kapp").TextColor(t.TextMuted)
			ui.TextInput(c, &a.name).Placeholder("输入你的姓名").Label("姓名").Width(280)
			greeting := "你好！"
			if a.name != "" {
				greeting = fmt.Sprintf("你好，%s！", a.name)
			}
			ui.Text(c, greeting)
		}
	})
}
