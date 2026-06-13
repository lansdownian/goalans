package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lansdownian/goalans/internal/ui"
)

func (m model) handleWorldCupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Allow Esc even while loading (central handler already routed sub-view Esc).
	if isEsc(msg) && m.wcSubView != wcSubViewGroupGrid {
		switch m.wcSubView {
		case wcSubViewGroupDetail, wcSubViewBracket, wcSubViewGroups:
			m.wcSubView = wcSubViewGroupGrid
			return m, tea.ClearScreen
		}
	}
	if m.wcLoading {
		return m, nil
	}
	switch m.wcSubView {
	case wcSubViewGroups:
		return m.handleWCGroupsKeys(msg)
	case wcSubViewGroupDetail:
		return m.handleWCGroupDetailKeys(msg)
	case wcSubViewBracket:
		return m.handleWCBracketKeys(msg)
	default:
		return m.handleWCGroupGridKeys(msg)
	}
}

func (m model) handleWCGroupsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.wcData == nil {
		return m, nil
	}
	switch msg.String() {
	case "esc":
		m.wcSubView = wcSubViewGroupGrid
		return m, tea.ClearScreen
	case "enter":
		if item, ok := m.wcGroupsList.SelectedItem().(ui.WCGroupItem); ok {
			for i, g := range m.wcData.Groups {
				if g.Letter == item.Group.Letter {
					m.wcSelectedGroup = i
					break
				}
			}
			m.wcSubView = wcSubViewGroupDetail
			return m, tea.ClearScreen
		}
	case "b":
		if len(m.wcData.KnockoutRounds) > 0 {
			m.wcBracketScroll = 0
			m.wcSubView = wcSubViewBracket
			return m, tea.ClearScreen
		}
	default:
		var cmd tea.Cmd
		m.wcGroupsList, cmd = m.wcGroupsList.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) handleWCGroupDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.wcSubView = wcSubViewGroupGrid
		return m, tea.ClearScreen
	}
	return m, nil
}

func (m model) handleWCBracketKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.wcSubView = wcSubViewGroupGrid
		return m, tea.ClearScreen
	case "j", "down":
		if m.wcBracketLines > 0 && m.wcBracketScroll < m.wcBracketLines-1 {
			m.wcBracketScroll++
		}
	case "k", "up":
		if m.wcBracketScroll > 0 {
			m.wcBracketScroll--
		}
	}
	return m, nil
}

func (m model) handleWCGroupGridKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.wcData == nil {
		return m, nil
	}
	n := len(m.wcData.Groups)
	if n == 0 {
		return m, nil
	}
	cols := 2
	if m.width > 120 {
		cols = 4
	} else if m.width > 80 {
		cols = 3
	}
	switch msg.String() {
	case "esc":
		m.wcSubView = wcSubViewGroupGrid
		return m, tea.ClearScreen
	case "enter":
		m.wcSelectedGroup = m.wcGridSelectedIdx
		m.wcSubView = wcSubViewGroupDetail
		return m, tea.ClearScreen
	case "t":
		m.wcSubView = wcSubViewGroups
		return m, tea.ClearScreen
	case "b":
		if len(m.wcData.KnockoutRounds) > 0 {
			m.wcBracketScroll = 0
			m.wcSubView = wcSubViewBracket
			return m, tea.ClearScreen
		}
	case "right", "l":
		if m.wcGridSelectedIdx < n-1 {
			m.wcGridSelectedIdx++
		}
	case "left", "h":
		if m.wcGridSelectedIdx > 0 {
			m.wcGridSelectedIdx--
		}
	case "down", "j":
		if m.wcGridSelectedIdx+cols < n {
			m.wcGridSelectedIdx += cols
		}
	case "up", "k":
		if m.wcGridSelectedIdx-cols >= 0 {
			m.wcGridSelectedIdx -= cols
		}
	}
	return m, nil
}

func (m model) handleWCData(msg wcDataMsg) (tea.Model, tea.Cmd) {
	m.wcLoading = false
	if msg.err != nil {
		m.wcLastError = "Failed to load World Cup data"
		return m, nil
	}
	m.wcData = msg.data
	m.wcLastError = ""
	if msg.data != nil {
		m.wcBracketLines = msg.data.BracketLineCount()
		items := make([]list.Item, len(msg.data.Groups))
		for i, g := range msg.data.Groups {
			items[i] = ui.WCGroupItem{Group: g}
		}
		m.wcGroupsList.SetItems(items)
	}
	return m, nil
}
