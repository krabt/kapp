package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/mobile"
)

const (
	defaultSidebarWidth float32 = 120
	minSidebarWidth     float32 = 60
)

var (
	overviewIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="m3 10 9-7 9 7v10a1 1 0 0 1-1 1h-5v-7H9v7H4a1 1 0 0 1-1-1z" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"/></svg>`))
	counterIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 5v14M5 12h14" fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2"/></svg>`))
)

type App struct {
	name         string
	count        int
	page         string
	mobilePage   int
	sidebarWidth float32
}

func New() *App {
	return &App{page: "home", sidebarWidth: defaultSidebarWidth}
}

func (a *App) View(c *ui.Context) {
	width, _ := c.Size()
	if width < 720 {
		a.mobileView(c)
		return
	}
	a.desktopView(c)
}

func (a *App) desktopView(c *ui.Context) {
	split := ui.Split(c, &a.sidebarWidth, func() {
		a.sidebarView(c)
	}, func() {
		a.mainView(c)
	}).SplitReset(defaultSidebarWidth).Fill()
	if split.Changed() && a.sidebarWidth < minSidebarWidth {
		a.sidebarWidth = minSidebarWidth
	}
}

func (a *App) mobileView(c *ui.Context) {
	if a.page == "counter" {
		a.mobilePage = 1
	} else {
		a.mobilePage = 0
	}
	ui.Column(c).Fill().Children(func() {
		a.mainView(c).Grow(1)
		bar := mobile.BottomBar(c, &a.mobilePage,
			mobile.BarButton{Label: "概览", Icon: mobileIcon(overviewIcon)},
			mobile.BarButton{Label: "计数器", Icon: mobileIcon(counterIcon)},
		)
		if bar.Changed() {
			if a.mobilePage == 1 {
				a.page = "counter"
			} else {
				a.page = "home"
			}
		}
	})
}

func mobileIcon(icon *ui.SVG) func(*ui.Painter, ui.Rect, ui.Color) {
	return func(p *ui.Painter, bounds ui.Rect, color ui.Color) {
		p.Icon(icon, bounds, color)
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

func (a *App) mainView(c *ui.Context) ui.Element {
	t := c.Theme()
	width, _ := c.Size()
	padding := float32(28)
	if width < 720 {
		padding = 20
	}
	return ui.Column(c).Grow(1).Gap(16).Padding(padding).Children(func() {
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
