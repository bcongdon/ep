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

func main() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true).
		SetFixed(1, 1)

	emojis, _ := getEmojis()
	justEmojis := []string{}
	for _, e := range emojis {
		for _, emoji := range e {
			justEmojis = append(justEmojis, emoji.Char)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(justEmojis)))

	cols, rows := 10, 200
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			table.SetCell(r, c,
				tview.NewTableCell(" "+justEmojis[word]+" "))
			word = (word + 1) % len(justEmojis)
		}
	}
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		clipboard.WriteAll(cell.Text)
		app.Stop()
	})
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
