// Demo code for the Table primitive.
package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Emoji struct {
	Keywords []string `json:"keywords"`
	Char     string   `json:"char"`
}

func getEmojis() (map[string][]Emoji, error) {
	raw, err := ioutil.ReadFile("emojis.json")
	if err != nil {
		return nil, err
	}

	nameMap := make(map[string]Emoji)
	err = json.Unmarshal(raw, &nameMap)
	if err != nil {
		return nil, err
	}

	keywordMap := make(map[string][]Emoji)
	for _, emoji := range nameMap {
		for _, keyword := range emoji.Keywords {
			if _, ok := keywordMap[keyword]; ok {
				keywordMap[keyword] = []Emoji{}
			}

			keywordMap[keyword] = append(keywordMap[keyword], emoji)
		}
	}

	return keywordMap, nil
}

func drawEmojis(table *tview.Table, emojis []string) {
	cols, rows := 10, 200
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			table.SetCell(r, c,
				tview.NewTableCell(" "+emojis[word]+" "))
			word = (word + 1) % len(emojis)
		}
	}
}

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true).
		SetFixed(1, 1)

	inputField := tview.NewInputField().
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyDown {
			app.SetFocus(table)
			table.SetSelectable(true, true)
		}
		return event
	})

	grid := tview.NewGrid().
		SetRows(1, 1).
		SetColumns(1, 1).
		AddItem(inputField, 0, 0, 1, 3, 0, 0, true).
		AddItem(table, 2, 0, 1, 3, 0, 0, false)

	emojis, _ := getEmojis()
	justEmojis := []string{}
	for _, e := range emojis {
		for _, emoji := range e {
			justEmojis = append(justEmojis, emoji.Char)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(justEmojis)))
	drawEmojis(table, justEmojis)

	table.SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		clipboard.WriteAll(cell.Text)
		app.Stop()
	}).SetSelectable(false, false)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		if event.Key() == tcell.KeyUp && row == 0 {
			app.SetFocus(inputField)
			table.SetSelectable(false, false)
		}
		return event
	})

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
