package setup

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	battleships "github.com/kovansky/wp-battleships"
	"github.com/kovansky/wp-battleships/parts"
	"github.com/kovansky/wp-battleships/tui"
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/kovansky/wp-battleships/tui/common"
	"github.com/kovansky/wp-battleships/tui/lobby"
	"github.com/kovansky/wp-battleships/tui/wait"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
	"strings"
)

const SizeMax = 4

type Setup struct {
	ctx context.Context
	log zerolog.Logger

	theme battleships.Theme

	subcomponents map[string]tea.Model

	board board.NewSingle
	input textinput.Model
	ships map[int][]parts.Ship

	errorText              string
	protectedFields        map[string]parts.State
	currentProtectedFields map[string]parts.State

	shipsLimitPerSize map[int]int
	shipsLimit        int

	asciiRender *figlet4go.AsciiRender
}

func Create(ctx context.Context, theme battleships.Theme) Setup {
	log := ctx.Value(battleships.ContextKeyLog).(zerolog.Logger)

	asciiRender := figlet4go.NewAsciiRender()

	header := common.CreateHeader("Battleships", theme, asciiRender)
	b := board.InitNewSingle(battleships.Themes.Player, map[string]parts.State{})

	input := textinput.New()
	input.Placeholder = "Place next ship part"
	input.CharLimit = 5
	input.Width = 25
	input.Focus()

	perSize := map[int]int{
		1: 4,
		2: 3,
		3: 2,
		4: 1,
	}
	shipsLimit := 0
	for _, lim := range perSize {
		shipsLimit += lim
	}

	return Setup{
		ctx:   ctx,
		log:   log,
		theme: theme,
		subcomponents: map[string]tea.Model{
			"header": header,
		},
		board: b,
		input: input,
		ships: map[int][]parts.Ship{
			0: make([]parts.Ship, 0, 1),
			1: make([]parts.Ship, 0, perSize[1]),
			2: make([]parts.Ship, 0, perSize[2]),
			3: make([]parts.Ship, 0, perSize[3]),
			4: make([]parts.Ship, 0, perSize[4]),
		},
		shipsLimit:             shipsLimit,
		shipsLimitPerSize:      perSize,
		errorText:              "",
		protectedFields:        map[string]parts.State{},
		currentProtectedFields: map[string]parts.State{},
		asciiRender:            asciiRender,
	}
}

func (c Setup) Init() tea.Cmd {
	return textinput.Blink
}

func (c Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			c.errorText = ""
			value := c.input.Value()

			if len(value) < 2 {
				c.errorText = "too short value"
				return c, nil
			}

			switch strings.ToLower(value) {
			case "ok", "next":
				newC := c.finishShip()
				newC.board.SetBoard(newC.protectedFields)

				newC.input.SetValue("")

				return newC, nil
			case "start":
				return c.finishSetup()
			default:
				newC := c.addShip(value)

				fullBoard := map[string]parts.State{}
				for k, v := range newC.protectedFields {
					fullBoard[k] = v
				}
				for k, v := range newC.currentProtectedFields {
					fullBoard[k] = v
				}

				newC.board.SetBoard(fullBoard)

				return newC, nil
			}
		case "ctrl+c":
			return c, tea.Quit
		}
	}

	c.board, cmd = c.board.Update(msg)
	cmds = append(cmds, cmd)
	c.input, cmd = c.input.Update(msg)
	cmds = append(cmds, cmd)

	for name, cmp := range c.subcomponents {
		c.subcomponents[name], cmd = cmp.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

func (c Setup) View() string {
	layout := lipgloss.JoinVertical(lipgloss.Center,
		c.subcomponents["header"].View(),

		lipgloss.JoinHorizontal(lipgloss.Top,
			c.board.View(),

			lipgloss.NewStyle().MarginLeft(2).Render(
				fmt.Sprintf("Ships count:\n"+
					"* One-masted: %d/%d\n"+
					"* Two-masted: %d/%d\n"+
					"* Three-masted: %d/%d\n"+
					"* Four-masted: %d/%d\n",
					len(c.ships[1]), cap(c.ships[1]),
					len(c.ships[2]), cap(c.ships[2]),
					len(c.ships[3]), cap(c.ships[3]),
					len(c.ships[4]), cap(c.ships[4]),
				),
			),
		),
		"",
		c.input.View(),
	)

	if len(c.errorText) > 0 {
		layout = lipgloss.JoinVertical(lipgloss.Center,
			layout,
			c.errorText,
		)
	}

	layout = lipgloss.JoinVertical(lipgloss.Center,
		layout,
		"\n\n\n",
		"type in field identifiers to place ships",
		"ok/next to submit ship",
		"start to save board",
	)

	return lipgloss.JoinHorizontal(lipgloss.Center, layout)
}

func (c Setup) countShips() int {
	var count int

	for _, ships := range c.ships {
		count += len(ships)
	}

	return count
}

func (c Setup) finishSetup() (Setup, tea.Cmd) {
	var shipsBoard []string

	for i, shipsOfSize := range c.ships {
		if i == 0 {
			continue
		}

		if len(shipsOfSize) != cap(shipsOfSize) {
			c.errorText = "you don't have enough ships!"
			return c, nil
		}

		for _, ship := range shipsOfSize {
			for f := range ship.Ship() {
				shipsBoard = append(shipsBoard, f)
			}
		}
	}

	battleships.PlayerData.Board = shipsBoard

	var targetStage tui.Stage
	targetStage = tui.StageWait
	if battleships.PlayerData.PlayMode == battleships.PlayModeChallenge {
		targetStage = tui.StageLobby
	}

	var app tea.Model
	switch targetStage {
	case tui.StageLobby:
		players, err := battleships.ServerClient.ListPlayers()
		if err != nil {
			c.log.Fatal().Err(err).Msg("failed to list players")
		}

		app = lobby.Create(c.ctx, battleships.Themes.Global, players)
		break
	default:
		var (
			game     battleships.Game
			gamePost battleships.GamePost
			err      error
		)

		gamePost = battleships.GamePost{
			Wpbot:  false,
			Coords: battleships.PlayerData.Board,
		}
		if len(battleships.PlayerData.Nick) > 0 {
			gamePost.Nick = battleships.PlayerData.Nick
			gamePost.Desc = battleships.PlayerData.Description
		}

		game, err = battleships.ServerClient.InitGame(gamePost)
		if err != nil {
			c.log.Fatal().Err(err).Msg("Couldn't initialize game")
		}

		battleships.GameInstance = game

		app = wait.Create(c.ctx, c.theme)
	}

	return c, func() tea.Msg {
		return tui.ApplicationStageChangeMsg{
			From:  tui.StageSetup,
			Stage: targetStage,
			Model: app,
		}
	}
}

func (c Setup) finishShip() Setup {
	var (
		ship parts.Ship
		err  error
	)
	if len(c.ships[0]) == 0 {
		c.errorText = "there is no ship being currently built"
		return c
	}
	ship = c.ships[0][0]

	shipCategory := c.ships[ship.Size()]
	if len(shipCategory) >= cap(shipCategory) {
		c.errorText = "you have reached the limit of ships of that size"
		c.ships[0] = make([]parts.Ship, 0, 1)
		c.currentProtectedFields = map[string]parts.State{}
		return c
	}

	ship, err = ship.Finish()
	if err != nil {
		c.errorText = "the ship size is wrong. Acceptable ship sizes: 1, 2, 3, 4"
		return c
	}

	shipCategory = append(shipCategory, ship)
	c.ships[ship.Size()] = shipCategory
	c.ships[0] = make([]parts.Ship, 0, 1)
	c.currentProtectedFields = map[string]parts.State{}

	var (
		stateHit    parts.State = parts.FieldHit
		hitPriority             = stateHit.Priority()
	)
	for f := range ship.Ship() {
		// If field is currently denoted empty
		// Or if it's here as a protected (edge) - overwrite
		if current, contains := c.protectedFields[f]; !contains || current.Priority() < hitPriority {
			c.protectedFields[f] = stateHit
			if err != nil {
				c.errorText = "error while saving ship"
				return c
			}
		}
	}
	protected, err := ship.Protected()
	if err != nil {
		return c
	}
	for identifier, f := range protected {
		// If field is currently denoted empty
		// Or if it's there as a protected (edge/corner), but the type differ (i.e. was edge, is corner) - overwrite
		if _, contains := c.protectedFields[identifier]; !contains {
			c.protectedFields[identifier] = f.State
		}
	}

	return c
}

func (c Setup) addShip(value string) Setup {
	var (
		ship parts.Ship
		err  error
	)

	value = strings.ToUpper(value)

	if len(c.ships[0]) == 0 {
		if c.countShips() >= c.shipsLimit {
			c.errorText = "you have reached the limit of ships you can place"
			c.input.SetValue("")
			return c
		}

		ship = parts.NewShip()
	} else {
		ship = c.ships[0][0]
	}

	if _, contains := ship.Ship()[value]; contains {
		c.errorText = "this field is already selected"
		return c
	}

	if _, contains := c.protectedFields[value]; contains {
		c.errorText = "you cannot place a ship on that field"
		return c
	}

	ship, err = ship.Add(value)
	if err != nil {
		var (
			size              = &parts.ErrShipSize{}
			nonAdjacent       = &parts.ErrFieldNonadjacent{}
			malformed         = &parts.ErrFieldMalformed{}
			adjacentMalformed = &parts.ErrAdjacentFieldMalformed{}
		)

		c.errorText = "could not add field to ship"
		if errors.As(err, size) {
			c.errorText = "ship has reached it's maximum size"
			c.input.SetValue("")
		} else if errors.As(err, nonAdjacent) {
			c.errorText = "new field has to touch the current ship"
		} else if errors.As(err, malformed) || errors.As(err, adjacentMalformed) {
			c.errorText = "field input incorrect"
		}

		return c
	}

	c.currentProtectedFields = map[string]parts.State{}

	var stateShip parts.State = parts.FieldShip
	shipPriority := stateShip.Priority()
	for f := range ship.Ship() {
		if current, contains := c.currentProtectedFields[f]; !contains || current.Priority() < shipPriority {
			c.currentProtectedFields[f] = stateShip
		}
	}
	protected, err := ship.Protected()
	if err != nil {
		return c
	}
	for identifier, f := range protected {
		c.currentProtectedFields[identifier] = f.State

		if _, fullContains := c.protectedFields[identifier]; ship.Size() < SizeMax &&
			f.State == parts.FieldProtected &&
			!fullContains {
			c.currentProtectedFields[identifier] = parts.FieldPotential
		}
	}

	if len(c.ships[0]) == 0 {
		c.ships[0] = append(c.ships[0], ship)
	} else {
		c.ships[0][0] = ship
	}

	c.input.SetValue("")

	return c
}
