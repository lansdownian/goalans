package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/fotmob"
	"github.com/lansdownian/goalans/internal/ui"
)

type view int

const (
	viewMenu view = iota
	viewLive
	viewFinished
	viewWorldCup
)

type wcSubView int

const (
	wcSubViewGroupGrid wcSubView = iota
	wcSubViewGroups
	wcSubViewGroupDetail
	wcSubViewBracket
)

type model struct {
	width  int
	height int

	currentView view
	menuIndex   int

	useMock    bool
	client     *fotmob.Client
	loading    bool
	lastError  string
	polling    bool
	pollGen    int

	matches      []api.Match
	matchDetails *api.MatchDetails
	matchList    list.Model
	spinner      spinner.Model

	// World Cup
	wcData            *api.WorldCupData
	wcLoading         bool
	wcLastError       string
	wcSubView         wcSubView
	wcSelectedGroup   int
	wcGridSelectedIdx int
	wcGroupsList      list.Model
	wcBracketScroll   int
	wcBracketLines    int
}

func New(useMock bool) model {
	d := list.NewDefaultDelegate()
	d.SetHeight(2)
	l := list.New([]list.Item{}, d, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	cursor, prompt := ui.FilterInputStyles()
	l.Styles.FilterCursor = cursor
	l.FilterInput.PromptStyle = prompt
	l.FilterInput.Cursor.Style = cursor

	wcList := list.New([]list.Item{}, ui.NewWCGroupDelegate(), 0, 0)
	wcList.SetShowTitle(false)
	wcList.SetShowStatusBar(true)
	wcList.SetFilteringEnabled(true)
	wcList.SetShowFilter(true)
	wcList.Styles.FilterCursor = cursor
	wcList.FilterInput.PromptStyle = prompt
	wcList.FilterInput.Cursor.Style = cursor

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = ui.AccentStyle()

	return model{
		currentView: viewMenu,
		useMock:     useMock,
		client:      fotmob.NewClient(),
		matchList:   l,
		wcGroupsList: wcList,
		spinner:     s,
		wcSubView:   wcSubViewGroupGrid,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}
