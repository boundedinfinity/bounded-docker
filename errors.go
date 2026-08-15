package main

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

func newErrorModel(errs []error) tea.Model {
	err2row := func(i int, err error) table.Row {
		return table.Row{
			fmt.Sprintf("%d", i),
			err.Error(),
		}
	}

	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Error Text", Width: 120},
	}

	return newTableModel("Errors", columns, err2row, errs)
}

func createFakeErrors() []error {
	return []error{
		errors.New("Error 1"),
		errors.New("Error 2"),
		errors.New("Error 3"),
	}
}
