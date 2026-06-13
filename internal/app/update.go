package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/ui"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.currentView == viewLive || m.currentView == viewFinished {
			m.resizeList()
		}
		if m.currentView == viewWorldCup && m.wcSubView == wcSubViewGroups {
			m.wcGroupsList.SetSize(m.width, max(m.height-6, 8))
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading || m.wcLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case matchesMsg:
		m.loading = false
		if msg.err != nil {
			m.lastError = msg.err.Error()
			m.matches = nil
			m.matchList.SetItems(nil)
			return m, nil
		}
		m.lastError = ""
		m.matches = msg.matches
		m.matchList.SetItems(ui.ToListItems(msg.matches))
		if len(msg.matches) > 0 {
			m.matchList.Select(0)
			cmds = append(cmds, fetchMatchDetails(m.client, msg.matches[0].ID, m.useMock))
			m.loading = true
		}
		return m, tea.Batch(cmds...)

	case matchDetailsMsg:
		m.loading = false
		if msg.err != nil {
			m.lastError = msg.err.Error()
			m.matchDetails = nil
			return m, nil
		}
		m.lastError = ""
		m.matchDetails = msg.details
		if m.currentView == viewLive && msg.details != nil && msg.details.Status == api.MatchStatusLive {
			m.polling = true
			m.pollGen++
			cmds = append(cmds, schedulePoll(msg.details.ID, m.pollGen))
		}
		return m, tea.Batch(cmds...)

	case wcDataMsg:
		return m.handleWCData(msg)

	case pollTickMsg:
		if !m.polling || msg.gen != m.pollGen || m.matchDetails == nil || m.matchDetails.ID != msg.matchID {
			return m, nil
		}
		m.loading = true
		cmds = append(cmds, fetchMatchDetails(m.client, msg.matchID, m.useMock), schedulePoll(msg.matchID, msg.gen))
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		if isEsc(msg) {
			if m.handleEscFilter(msg) {
				return m, nil
			}
			// WC sub-views (not grid) use Esc to step back internally.
			if m.currentView == viewWorldCup && m.wcSubView != wcSubViewGroupGrid {
				return m.handleWorldCupKeys(msg)
			}
			if m.currentView != viewMenu {
				return m.resetToMainView()
			}
			return m, nil
		}
		if m.currentView == viewWorldCup {
			return m.handleWorldCupKeys(msg)
		}
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "j", "down":
		if m.currentView == viewMenu {
			if m.menuIndex < 2 {
				m.menuIndex++
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.matchList, cmd = m.matchList.Update(msg)
		if id, ok := m.selectedMatchID(); ok {
			m.loading = true
			return m, tea.Batch(cmd, fetchMatchDetails(m.client, id, m.useMock))
		}
		return m, cmd

	case "k", "up":
		if m.currentView == viewMenu {
			if m.menuIndex > 0 {
				m.menuIndex--
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.matchList, cmd = m.matchList.Update(msg)
		if id, ok := m.selectedMatchID(); ok {
			m.loading = true
			return m, tea.Batch(cmd, fetchMatchDetails(m.client, id, m.useMock))
		}
		return m, cmd

	case "enter":
		if m.currentView == viewMenu {
			m.loading = true
			m.lastError = ""
			m.matchDetails = nil
			m.polling = false
			m.pollGen++
			switch m.menuIndex {
			case 0:
				m.currentView = viewLive
				return m, tea.Batch(m.spinner.Tick, fetchLiveMatches(m.client, m.useMock))
			case 1:
				m.currentView = viewFinished
				return m, tea.Batch(m.spinner.Tick, fetchFinishedMatches(m.client, m.useMock))
			case 2:
				m.currentView = viewWorldCup
				m.wcData = nil
				m.wcLoading = true
				m.wcLastError = ""
				m.wcSubView = wcSubViewGroupGrid
				m.wcGridSelectedIdx = 0
				return m, tea.Batch(m.spinner.Tick, fetchWorldCupData(m.client, m.useMock))
			}
		}
	}

	if m.currentView == viewLive || m.currentView == viewFinished {
		var cmd tea.Cmd
		m.matchList, cmd = m.matchList.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) selectedMatchID() (int, bool) {
	item := m.matchList.SelectedItem()
	if item == nil {
		return 0, false
	}
	if mi, ok := item.(ui.MatchItem); ok {
		return mi.ID, true
	}
	return 0, false
}

func (m *model) resizeList() {
	if m.width == 0 || m.height == 0 {
		return
	}
	leftW := max(m.width*38/100, 28)
	h := max(m.height-6, 8)
	m.matchList.SetSize(leftW-4, h)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	switch m.currentView {
	case viewMenu:
		return ui.RenderMenu(m.width, m.height, m.menuIndex, m.useMock, m.loading, m.spinner)
	case viewWorldCup:
		return m.renderWorldCupView()
	default:
		m.resizeList()
		title := "Live Matches"
		if m.currentView == viewFinished {
			title = "Finished Today"
		}
		return ui.RenderMatchView(m.width, m.height, title, m.matchList, m.matchDetails, m.spinner, m.loading, m.lastError, m.polling)
	}
}

func (m model) renderWorldCupView() string {
	switch m.wcSubView {
	case wcSubViewGroupDetail:
		return ui.RenderWorldCupGroupDetail(m.width, m.height, m.wcData, m.wcSelectedGroup)
	case wcSubViewBracket:
		return ui.RenderWorldCupBracket(m.width, m.height, m.wcData, m.wcBracketScroll)
	case wcSubViewGroups:
		return ui.RenderWorldCupGroups(m.width, m.height, m.wcData, m.wcGroupsList, m.wcLoading, m.wcLastError)
	default:
		if m.wcLoading && m.wcData == nil {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
				ui.DimStyle().Render("Loading World Cup 2026…")+" "+m.spinner.View())
		}
		return ui.RenderWorldCupGroupGrid(m.width, m.height, m.wcData, m.wcGridSelectedIdx)
	}
}

func isEsc(msg tea.KeyMsg) bool {
	return msg.Type == tea.KeyEsc || msg.String() == "esc"
}

// handleEscFilter returns true if Esc was consumed to cancel an active list filter.
func (m model) handleEscFilter(msg tea.KeyMsg) bool {
	var filtering bool
	switch m.currentView {
	case viewLive, viewFinished:
		filtering = m.matchList.FilterState() == list.Filtering ||
			m.matchList.FilterState() == list.FilterApplied
	case viewWorldCup:
		if m.wcSubView == wcSubViewGroups {
			filtering = m.wcGroupsList.FilterState() == list.Filtering ||
				m.wcGroupsList.FilterState() == list.FilterApplied
		}
	}
	if !filtering {
		return false
	}
	switch m.currentView {
	case viewLive, viewFinished:
		m.matchList, _ = m.matchList.Update(msg)
	case viewWorldCup:
		m.wcGroupsList, _ = m.wcGroupsList.Update(msg)
	}
	return true
}

func (m model) resetToMainView() (tea.Model, tea.Cmd) {
	m.polling = false
	m.pollGen++
	m.matchDetails = nil
	m.matches = nil
	m.wcData = nil
	m.wcLoading = false
	m.wcSubView = wcSubViewGroupGrid
	m.wcGridSelectedIdx = 0
	m.wcSelectedGroup = 0
	m.wcBracketScroll = 0
	m.lastError = ""
	m.wcLastError = ""
	m.loading = false
	m.currentView = viewMenu
	return m, tea.ClearScreen
}
