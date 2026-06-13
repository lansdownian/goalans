package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/lansdownian/goalans/internal/api"
	"github.com/lansdownian/goalans/internal/ui/worldcup"
)

type WCGroupItem = worldcup.WCGroupItem

func NewWCGroupDelegate() list.DefaultDelegate {
	return worldcup.NewWCGroupDelegate()
}

func RenderWorldCupGroups(width, height int, wcData *api.WorldCupData, groupsList list.Model, loading bool, lastErr string) string {
	return worldcup.RenderGroupsList(width, height, wcData, groupsList, loading, lastErr, "")
}

func RenderWorldCupGroupDetail(width, height int, wcData *api.WorldCupData, groupIdx int) string {
	return worldcup.RenderGroupDetail(width, height, wcData, groupIdx, "")
}

func RenderWorldCupGroupGrid(width, height int, wcData *api.WorldCupData, selectedGroupIdx int) string {
	return worldcup.RenderGroupGrid(width, height, wcData, selectedGroupIdx, "")
}

func RenderWorldCupBracket(width, height int, wcData *api.WorldCupData, scrollOffset int) string {
	return worldcup.RenderBracket(width, height, wcData, scrollOffset, "")
}

func RenderWorldCupUpcoming(width, height int, matches []api.Match, loading bool, lastErr string) string {
	return worldcup.RenderUpcoming(width, height, matches, loading, lastErr, "")
}
