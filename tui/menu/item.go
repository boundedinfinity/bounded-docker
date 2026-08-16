package menu

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type MenuItem interface {
	Title() string
	Count() int
	tea.Model
}

// ////////////////////////////////////////////////////////////////////////////////////////////////

func Find[T any](items []MenuItem, title string) (MenuItem, bool) {
	for _, item := range items {
		if item.Title() == title {
			return item, true
		}
	}

	return nil, false
}

func Item[T any](title string) MenuItem {
	return &countedItem[T]{
		title: title,
		count: 0,
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////

var _ MenuItem = &countedItem[any]{}

type countedItem[T any] struct {
	title string
	count int
}

func (this countedItem[T]) Title() string {
	return this.title
}

func (this countedItem[T]) Count() int {
	return this.count
}

func (this countedItem[T]) Init() tea.Cmd {
	return nil
}

func (this countedItem[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case T:
		this.count += 1
	case []T:
		this.count = len(msg)
	}
	return this, nil
}

func (this countedItem[T]) View() tea.View {
	text := fmt.Sprintf("%s [%d]", this.title, this.count)
	return tea.NewView(text)
}
