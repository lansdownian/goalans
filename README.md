# Goalans

A terminal app for following football (soccer) matches — live scores, today's finished results, and match details with events and statistics.

## What it provides

- **Live Matches** — currently in-play games across top European leagues and World Cup
- **Recent Results** — finished matches from the last 14 days (includes World Cup)
- **World Cup 2026** — group standings and knockout bracket (USA / Mexico / Canada)
- **Match details** — score, venue, referee, goal/card events, key stats
- **Live polling** — selected live matches refresh every 90 seconds
- **Filter** — press `/` to search the match list by team or league
- **Mock mode** — run fully offline with sample data (`--mock`)

### Leagues covered (default)

Premier League, La Liga, Bundesliga, Serie A, Ligue 1, FIFA World Cup

Data is fetched from public FotMob league and match pages.

## Requirements

- **Go 1.24+**
- A terminal with reasonable size (80×24 or larger recommended)
- Network access (unless using `--mock`)

## How to run

### Build and run

```bash
cd Goalans
go mod tidy
go build -o goalans .
./goalans
```

### Run without building

```bash
go run . 
```

### Mock mode (no network)

```bash
./goalans --mock
```

Useful for trying the UI when no matches are live or when offline.

## Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Move selection |
| `Enter` | Open a menu view; select a match in list |
| `↑/↓/←/→` or `hjkl` | In World Cup grid: navigate groups |
| `Enter` | In World Cup grid: open group detail |
| `t` | In World Cup: table/list view of groups |
| `b` | In World Cup: knockout bracket |
| `/` | Filter matches |
| `Esc` | Back to main menu |
| `q` | Quit |

## Project layout

```
Goalans/
├── main.go              Entry point
├── cmd/                 CLI (Cobra root command)
├── internal/
│   ├── api/             Shared types (Match, MatchDetails, …)
│   ├── app/             Bubble Tea application (model, update, commands)
│   ├── data/            Default leagues and mock fixtures
│   ├── fotmob/          FotMob page fetcher and parser
│   └── ui/              Terminal rendering (lipgloss + bubbles)
└── README.md
```

## Notes

- Match details for live/network mode require selecting a match from the list first (the app stores the FotMob page URL from the list fetch).
- If the match list is empty, there may simply be no live or recent finished games in the tracked leagues.
- FotMob HTML structure can change; if fetches fail, use `--mock` to verify the app itself is working.

## License

This project is your own work. Third-party libraries used are open source (Bubble Tea, Cobra, Lip Gloss — see `go.mod`).
