package main

// type TimedError struct {
// 	err  error
// 	time time.Time
// }

// func (this TimedError) Error() string {
// 	return this.Error()
// }

// func newErrorModel(errs []error) (tea.Model, tea.Model) {
// 	err2row := func(i int, err error) table.Row {
// 		return table.Row{
// 			fmt.Sprintf("%d", i),
// 			err.Error(),
// 		}
// 	}

// 	columns := []table.Column{
// 		{Title: "ID", Width: 10},
// 		{Title: "Error Text", Width: 120},
// 	}

// 	return widget.Widget("errors", columns, err2row, errs)
// }

// func createFakeErrors() []error {
// 	return []error{}
// }
