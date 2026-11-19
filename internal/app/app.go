package app

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"kapp/internal/data"
	"kapp/internal/data/sql"
	"kapp/internal/store"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/mobile"
)

const (
	defaultSidebarWidth float32 = 120
	minSidebarWidth     float32 = 60
)

var (
	overviewIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="m3 10 9-7 9 7v10a1 1 0 0 1-1 1h-5v-7H9v7H4a1 1 0 0 1-1-1z" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"/></svg>`))
	counterIcon  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 5v14M5 12h14" fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2"/></svg>`))
	settingsIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><circle cx="12" cy="12" r="3" fill="none" stroke="currentColor" stroke-width="2"/><path d="M19.4 15a1.7 1.7 0 0 0 .3 1.9l.1.1-1.8 1.8-.1-.1a1.7 1.7 0 0 0-1.9-.3 1.7 1.7 0 0 0-1 1.6v.2h-2.6V20a1.7 1.7 0 0 0-1-1.6 1.7 1.7 0 0 0-1.9.3l-.1.1L7.6 17l.1-.1a1.7 1.7 0 0 0 .3-1.9A1.7 1.7 0 0 0 6.4 14H6v-2.6h.4A1.7 1.7 0 0 0 8 10.4a1.7 1.7 0 0 0-.3-1.9L7.6 8.4 9.4 6l.1.1A1.7 1.7 0 0 0 11.4 6a1.7 1.7 0 0 0 1-1.6V4h2.6v.4A1.7 1.7 0 0 0 16 6a1.7 1.7 0 0 0 1.9-.3l.1-.1L19.8 7l-.1.1a1.7 1.7 0 0 0-.3 1.9 1.7 1.7 0 0 0 1.6 1h.2v2.6H21a1.7 1.7 0 0 0-1.6 1.4z" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"/></svg>`))
)

type App struct {
	name         string
	count        int
	page         string
	mobilePage   int
	sidebarWidth float32
	store        *store.Store
	cleanup      func()
	dbError      string
}

func New() *App {
	dataDir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return &App{page: "home", sidebarWidth: defaultSidebarWidth, dbError: err.Error()}
	}

	dbPath := filepath.Join(dataDir, "kapp.db")
	if err := data.GenTable(dbPath, sql.Files, []string{"up"}, ""); err != nil {
		return &App{page: "home", sidebarWidth: defaultSidebarWidth, dbError: err.Error()}
	}
	db, cleanup, err := data.NewDB(dbPath, nil)
	if err != nil {
		return &App{page: "home", sidebarWidth: defaultSidebarWidth, dbError: err.Error()}
	}

	a := NewWithStore(store.NewStore(db))
	a.cleanup = cleanup
	return a
}

func NewWithStore(store *store.Store) *App {
	a := &App{page: "home", sidebarWidth: defaultSidebarWidth, store: store}
	if store != nil {
		values, err := store.GetKV(context.Background(), []string{"name", "count"})
		if err != nil {
			a.dbError = err.Error()
		} else {
			a.name = values["name"]
			if count, err := strconv.Atoi(values["count"]); err == nil {
				a.count = count
			}
		}
	}
	return a
}

func (a *App) saveSettings() {
	if a.store == nil {
		return
	}
	if err := a.store.SetKV(context.Background(), map[string]string{
		"name":  a.name,
		"count": strconv.Itoa(a.count),
	}); err != nil {
		a.dbError = err.Error()
	} else {
		a.dbError = ""
	}
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
	switch a.page {
	case "counter":
		a.mobilePage = 1
	case "settings":
		a.mobilePage = 2
	default:
		a.mobilePage = 0
	}
	ui.Column(c).Fill().Children(func() {
		a.mainView(c).Grow(1)
		bar := mobile.BottomBar(c, &a.mobilePage,
			mobile.BarButton{Label: "概览", Icon: mobileIcon(overviewIcon)},
			mobile.BarButton{Label: "计数器", Icon: mobileIcon(counterIcon)},
			mobile.BarButton{Label: "设置", Icon: mobileIcon(settingsIcon)},
		)
		if bar.Changed() {
			switch a.mobilePage {
			case 1:
				a.page = "counter"
			case 2:
				a.page = "settings"
			default:
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
			ui.SidebarItem(c, "settings", nil, "设置")
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
					a.saveSettings()
				}
				ui.Textf(c, "%d", a.count).FontSize(20).Width(48).TextAlign(ui.Center)
				if ui.PrimaryButton(c, "+").Width(44).Clicked() {
					a.count++
					a.saveSettings()
				}
			})
		case "settings":
			ui.Text(c, "设置").FontSize(26).Bold()
			ui.Text(c, "应用偏好会保存在本机 SQLite 数据库中。").TextColor(t.TextMuted)
			ui.TextInput(c, &a.name).Placeholder("输入你的姓名").Label("显示姓名").Width(280).OnChange(a.saveSettings)
			if ui.Button(c, "重置计数器").Clicked() {
				a.count = 0
				a.saveSettings()
			}
			if a.dbError == "" {
				ui.Text(c, "数据库状态：正常").TextColor(t.TextMuted)
			} else {
				ui.Text(c, "数据库错误："+a.dbError).TextColor(t.TextMuted)
			}
		default:
			ui.Text(c, "概览").FontSize(26).Bold()
			ui.Text(c, "欢迎使用 kapp").TextColor(t.TextMuted)
			ui.TextInput(c, &a.name).Placeholder("输入你的姓名").Label("姓名").Width(280).OnChange(a.saveSettings)
			greeting := "你好！"
			if a.name != "" {
				greeting = fmt.Sprintf("你好，%s！", a.name)
			}
			ui.Text(c, greeting)
		}
	})
}
