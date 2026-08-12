package main

// func (this *tui) keyHandle(event *tcell.EventKey) *tcell.EventKey {
// 	// bounded-docker.application.quit
// 	if event.Key() == tcell.KeyEscape || event.Rune() == 'q' || event.Rune() == 'Q' {
// 		this.cancel()
// 		return nil
// 	}

// 	// bounded-docker.application.errors.clear
// 	if event.Rune() == 'c' || event.Rune() == 'C' {
// 		this.errClear()
// 		return nil
// 	}

// 	if event.Key() == tcell.KeyEnter {
// 		// row, col := this.table.GetSelection()
// 		// fmt.Println(row, col)
// 	}

// 	if event.Rune() == '+' || event.Rune() == '=' {
// 		current := this.options
// 		current.cellPadding += 1
// 		this.sendOptions(current)
// 		return nil
// 	}

// 	if event.Rune() == '-' || event.Rune() == '_' {
// 		current := this.options
// 		current.cellPadding -= 1
// 		this.sendOptions(current)
// 		return nil
// 	}

// 	return event
// }
