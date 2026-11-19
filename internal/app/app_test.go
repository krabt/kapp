package app

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestView(t *testing.T) {
	a := New()
	tt := ui.NewTester(a.View, 760, 520)
	if err := tt.Click("计数器"); err != nil {
		t.Fatal(err)
	}
	if a.page != "counter" {
		t.Fatalf("page %q after selecting the counter", a.page)
	}
	if err := tt.Click("+"); err != nil {
		t.Fatal(err)
	}
	if a.count != 1 || !tt.HasText("1") {
		t.Errorf("count %d after a click, texts %q", a.count, tt.Texts())
	}
	if err := tt.Click("概览"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("姓名"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Ada")
	if !tt.HasText("你好，Ada！") {
		t.Errorf("texts %q after typing a name", tt.Texts())
	}
}

func TestSidebarCanResize(t *testing.T) {
	a := New()
	tt := ui.NewTester(a.View, 760, 520)

	dragSplitter := func(delta float32) {
		splitter, ok := tt.AccessByRole(ui.RoleSplitter, "")
		if !ok {
			t.Fatal("sidebar splitter not found")
		}
		x := splitter.Bounds.X + splitter.Bounds.W/2
		y := splitter.Bounds.Y + splitter.Bounds.H/2
		tt.Press(x, y)
		tt.Move(x+delta, y)
		tt.Release(x+delta, y)
	}

	initialWidth := a.sidebarWidth
	dragSplitter(64)
	if a.sidebarWidth <= initialWidth {
		t.Fatalf("sidebar width %v after dragging wider from %v", a.sidebarWidth, initialWidth)
	}

	widerWidth := a.sidebarWidth
	dragSplitter(-140)
	if a.sidebarWidth != minSidebarWidth {
		t.Fatalf("sidebar width %v after dragging below minimum %v from %v", a.sidebarWidth, minSidebarWidth, widerWidth)
	}
}