package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"charm.land/bubbles/v2/table"
	"github.com/boundedinfinity/docker-tui/docker"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/moby/moby/api/types/container"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	ctx         context.Context
	cancel      context.CancelFunc
	table       table.Model
	columns     []table.Column
	rows        []table.Row
	windowWidth int
	em          errModel
}

func (this model) Init() tea.Cmd {
	return nil
}

func (this model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	select {
	case <-this.ctx.Done():
		return this, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		this.windowWidth = msg.Width
		this.table = createContainerTable(msg.Width, this.columns, this.rows)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if this.table.Focused() {
				this.table.Blur()
			} else {
				this.table.Focus()
			}
		case "q", "ctrl+c":
			this.cancel()
			return this, nil
		case "enter":
			return this, tea.Batch(
				tea.Printf("Let's go to %s!", this.table.SelectedRow()[1]),
			)
		}
	case error:
		return this, nil
	case []container.Summary:
		this.rows = make([]table.Row, 0, len(msg))

		for i := range msg {
			this.rows = append(this.rows, table.Row{
				msg[i].ID,
				msg[i].Image,
				msg[i].Command,
				msg[i].Status,
			})
		}
	}

	this.table, cmd = this.table.Update(msg)
	return this, cmd
}

func (this model) View() tea.View {
	join := lipgloss.JoinHorizontal(
		0,
		this.table.View(),
		this.table.HelpView(),
		this.em,
	)

	return tea.NewView(baseStyle.Render(join))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	sigCh := make(chan os.Signal, 1)
	d, err := docker.NewDocker(wg, ctx)

	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		cancel()
	}()

	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		os.Exit(1)
	}

	rows := []table.Row{}

	m := model{
		columns: []table.Column{
			{Title: "ID"},
			{Title: "Image"},
			{Title: "Command"},
			{Title: "Status"},
		},
		rows: rows,
		em:   createErrorView(),
	}
	m.table = createContainerTable(0, m.columns, m.rows)

	p := tea.NewProgram(m)

	wg.Go(func() {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	})

	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-d.ErrCh:
				p.Send(err)
			case summary := <-d.SummaryCh:
				p.Send(summary)
			}
		}
	})

	d.GetSummary()

	wg.Wait()
}
