package app

import (
	"context"
	"time"

	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/data"
	"github.com/lansdownian/goalans/internal/fotmob"
	tea "github.com/charmbracelet/bubbletea"
)

func fetchWorldCupData(client *fotmob.Client, useMock bool) tea.Cmd {
	return func() tea.Msg {
		if useMock {
			return wcDataMsg{data: data.MockWorldCupData()}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()
		wc, err := client.WorldCupData(ctx, api.WCSeason2026)
		return wcDataMsg{data: wc, err: err}
	}
}
