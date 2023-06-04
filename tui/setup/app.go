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
	"github.com/kovansky/wp-battleships/tui/board"
	"github.com/kovansky/wp-battleships/tui/common"
	"github.com/mbndr/figlet4go"
	"github.com/rs/zerolog"
)

type Setup struct {
	ctx context.Context
	log zerolog.Logger

	theme battleships.Theme

	subcomponents map[string]tea.Model

	board board.NewSingle
	input textinput.Model
	ships map[int][]parts.Ship

	errorText       string
	protectedFields map[string]parts.State

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
		shipsLimit:        shipsLimit,
		shipsLimitPerSize: perSize,
		errorText:         "",
		protectedFields:   map[string]parts.State{},
		asciiRender:       asciiRender,
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

			newC := c.addShip(value)
			newC.board.SetBoard(newC.protectedFields)

			newC.input.SetValue("")

			return newC, nil
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
		c.board.View(),
		fmt.Sprintf("Ships count: %d", c.countShips()),
		c.input.View(),
	)

	if len(c.errorText) > 0 {
		layout = lipgloss.JoinVertical(lipgloss.Center,
			layout,
			c.errorText,
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, layout)
}

func (c Setup) countShips() int {
	var count int

	for _, ships := range c.ships {
		count += len(ships)
	}

	return count
}

func (c Setup) addShip(value string) Setup {
	var (
		ship parts.Ship
		err  error
	)
	if len(c.ships[0]) == 0 {
		if c.countShips() >= c.shipsLimit {
			c.errorText = "you have reached the limit of ships you can place"
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
		} else if errors.As(err, nonAdjacent) {
			c.errorText = "new field has to touch the current ship"
		} else if errors.As(err, malformed) || errors.As(err, adjacentMalformed) {
			c.errorText = "field input incorrect"
		}

		return c
	}

	var stateShip parts.State = parts.FieldShip
	shipPriority := stateShip.Priority()
	for f := range ship.Ship() {
		// If field is currently denoted empty
		// Or if it's here as a protected (edge) - overwrite
		if current, contains := c.protectedFields[f]; !contains || current.Priority() > shipPriority {
			c.protectedFields[f] = stateShip
			if err != nil {
				c.errorText = "could not add field to ship"
				return c
			}
		}
	}
	protected, err := ship.Protected()
	if err != nil {
		c.errorText = "could not add field to ship"
		return c
	}
	for identifier, f := range protected {
		// If field is currently denoted empty
		// Or if it's there as a protected (edge/corner), but the type differ (i.e. was edge, is corner) - overwrite
		if current, contains := c.protectedFields[identifier]; !contains ||
			(current.Priority() == f.State.Priority() && current != f.State) {
			c.protectedFields[identifier] = f.State
		}
	}

	if len(c.ships[0]) == 0 {
		c.ships[0] = append(c.ships[0], ship)
	} else {
		c.ships[0][0] = ship
	}

	return c
}
