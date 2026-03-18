package main

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func testModel() model {
	rng := rand.New(rand.NewSource(42))
	return model{
		systems:      GenerateSystems(rng),
		rng:          rng,
		tickInterval: 5 * time.Second,
	}
}

func sendKey(m model, key string) model {
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	return updated.(model)
}

func sendSpecialKey(m model, keyType tea.KeyType) model {
	updated, _ := m.Update(tea.KeyMsg{Type: keyType})
	return updated.(model)
}

func TestModelInit(t *testing.T) {
	m := testModel()
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init should return a tick command")
	}
}

func TestModelNavigateDown(t *testing.T) {
	m := testModel()
	m = sendKey(m, "j")
	if m.cursor != 1 {
		t.Errorf("cursor after j: got %d, want 1", m.cursor)
	}
	m = sendKey(m, "j")
	if m.cursor != 2 {
		t.Errorf("cursor after jj: got %d, want 2", m.cursor)
	}
}

func TestModelNavigateUp(t *testing.T) {
	m := testModel()
	m.cursor = 3
	m = sendKey(m, "k")
	if m.cursor != 2 {
		t.Errorf("cursor after k: got %d, want 2", m.cursor)
	}
}

func TestModelCursorBounds(t *testing.T) {
	m := testModel()
	// Can't go below 0
	m = sendKey(m, "k")
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0, got %d", m.cursor)
	}
	// Can't go past last system
	m.cursor = len(m.systems) - 1
	m = sendKey(m, "j")
	if m.cursor != len(m.systems)-1 {
		t.Errorf("cursor should stay at %d, got %d", len(m.systems)-1, m.cursor)
	}
}

func TestModelEnterDetail(t *testing.T) {
	m := testModel()
	m = sendSpecialKey(m, tea.KeyEnter)
	if m.view != viewDetail {
		t.Errorf("expected viewDetail, got %d", m.view)
	}
}

func TestModelEscBackToDashboard(t *testing.T) {
	m := testModel()
	m.view = viewDetail
	m = sendSpecialKey(m, tea.KeyEsc)
	if m.view != viewDashboard {
		t.Errorf("expected viewDashboard after esc, got %d", m.view)
	}
}

func TestModelQFromDetailGoesBack(t *testing.T) {
	m := testModel()
	m.view = viewDetail
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = updated.(model)
	if m.view != viewDashboard {
		t.Errorf("q from detail should go to dashboard, got view %d", m.view)
	}
	if cmd != nil {
		t.Error("q from detail should not quit")
	}
}

func TestModelQFromDashboardQuits(t *testing.T) {
	m := testModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Error("q from dashboard should produce quit command")
	}
}

func TestModelHelpToggle(t *testing.T) {
	m := testModel()
	m = sendKey(m, "?")
	if m.view != viewHelp {
		t.Errorf("expected viewHelp, got %d", m.view)
	}
	m = sendKey(m, "?")
	if m.view != viewDashboard {
		t.Errorf("expected viewDashboard after help toggle, got %d", m.view)
	}
}

func TestModelTickRefresh(t *testing.T) {
	m := testModel()
	updated, cmd := m.Update(tickMsg(time.Now()))
	m = updated.(model)
	if cmd == nil {
		t.Error("tick should schedule next tick")
	}
	if len(m.systems) != 6 {
		t.Errorf("expected 6 systems after tick, got %d", len(m.systems))
	}
}

func TestModelNavigationIgnoredInDetail(t *testing.T) {
	m := testModel()
	m.view = viewDetail
	m.cursor = 0
	m = sendKey(m, "j")
	if m.cursor != 0 {
		t.Errorf("j should not move cursor in detail view, got %d", m.cursor)
	}
}

func TestDashboardViewContainsAllSystems(t *testing.T) {
	m := testModel()
	view := m.View()
	for _, name := range []string{"Traffic", "Energy", "Water", "Waste", "Emergency", "Air Quality"} {
		if !strings.Contains(view, name) {
			t.Errorf("dashboard should contain %q", name)
		}
	}
}

func TestDashboardViewContainsHeader(t *testing.T) {
	m := testModel()
	view := m.View()
	if !strings.Contains(view, "SmartCity Orchestrator Dashboard") {
		t.Error("dashboard should contain title")
	}
	if !strings.Contains(view, "SYSTEM") {
		t.Error("dashboard should contain SYSTEM header")
	}
}

func TestDetailViewContainsMetrics(t *testing.T) {
	m := testModel()
	m.view = viewDetail
	m.cursor = 0
	view := m.View()
	if !strings.Contains(view, "Traffic") {
		t.Error("detail view should contain system name")
	}
	if !strings.Contains(view, "Congestion Level") {
		t.Error("detail view should contain metric names")
	}
}

func TestHelpViewContainsKeys(t *testing.T) {
	m := testModel()
	m.view = viewHelp
	view := m.View()
	if !strings.Contains(view, "Help") {
		t.Error("help view should contain Help title")
	}
	if !strings.Contains(view, "enter") {
		t.Error("help view should describe enter key")
	}
}

func TestNewModel(t *testing.T) {
	m := newModel(3 * time.Second)
	if m.tickInterval != 3*time.Second {
		t.Errorf("expected 3s tick, got %s", m.tickInterval)
	}
	if len(m.systems) != 6 {
		t.Errorf("expected 6 systems, got %d", len(m.systems))
	}
	if m.rng == nil {
		t.Error("rng should not be nil")
	}
	if m.view != viewDashboard {
		t.Error("initial view should be dashboard")
	}
	if m.cursor != 0 {
		t.Error("initial cursor should be 0")
	}
}
