//go:generate esc -o static.go -pkg main static
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/atotto/clipboard"
	ordering "github.com/bcongdon/emoji-ordering"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const numCols int = 10

var (
	outputFlag         = flag.String("output", "clipboard", "The output of ep. Choices: clipboard, stdout")
	noninteractiveFlag = flag.Bool("noninteractive", false, "If set, doesn't display emoji picker -- instead just outputting the first selection for the provided query.")
	emojis             = mustGetEmojis()
)

var usageFunc = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "%s [QUERY]\n", os.Args[0])
	flag.PrintDefaults()
}

type Emoji struct {
	Keywords []string `json:"keywords"`
	Char     string   `json:"char"`
}

func mustGetEmojis() map[string][]Emoji {
	emojis, err := getEmojis()
	if err != nil {
		panic(err)
	}
	return emojis
}

func getEmojis() (map[string][]Emoji, error) {
	raw := FSMustByte(false, "/static/emojis.json")

	nameMap := make(map[string]Emoji)
	err := json.Unmarshal(raw, &nameMap)
	if err != nil {
		return nil, err
	}

	keywordMap := make(map[string][]Emoji)
	for name, emoji := range nameMap {
		keywordMap[name] = append(keywordMap[name], emoji)
		for _, keyword := range emoji.Keywords {
			keywordMap[keyword] = append(keywordMap[keyword], emoji)
		}
	}

	return keywordMap, nil
}

func filterEmojis(emojis map[string][]Emoji, query string) []string {
	justEmojis := []string{}
	for key, e := range emojis {
		if !strings.Contains(key, query) {
			continue
		}
		for _, emoji := range e {
			justEmojis = append(justEmojis, emoji.Char)
		}
	}
	sort.Sort(ordering.EmojiSlice(justEmojis))
	return justEmojis
}

func drawEmojis(table *tview.Table, emojis map[string][]Emoji, query string) {
	filteredEmojis := filterEmojis(emojis, query)
	used := make(map[string]bool)
	gridIdx := 0
	for idx := 0; idx < len(filteredEmojis); idx++ {
		r, c := gridIdx/numCols, gridIdx%numCols
		emoji := filteredEmojis[idx]

		if _, alreadyUsed := used[emoji]; !alreadyUsed {
			table.SetCell(r, c, tview.NewTableCell(emoji))
			used[emoji] = true
			gridIdx++
		}
	}

	table.ScrollToBeginning()
	table.Select(0, 0)
}

func validateFlags() {
	if *outputFlag != "clipboard" && *outputFlag != "stdout" {
		log.Panicf("Invalid output method: %s\n", *outputFlag)
	}
}

func outputEmoji(emoji string) {
	switch *outputFlag {
	case "clipboard":
		clipboard.WriteAll(emoji)
	case "stdout":
		fmt.Println(emoji)
	default:
		log.Panicf("Unknown output method: %s", *outputFlag)
	}
}

func getEmojiName(emoji string) (string, bool) {
	raw := FSMustByte(false, "/static/emojis.json")

	nameMap := make(map[string]Emoji)
	err := json.Unmarshal(raw, &nameMap)
	if err != nil {
		return "", false
	}

	for name, e := range nameMap {
		if e.Char == emoji {
			return name, true
		}
	}
	return "", false
}

func runNoninterativeMode(emojis map[string][]Emoji, query string) {
	if len(query) == 0 {
		log.Panicln("A query must be specified in noninteractive mode.")
	}

	filteredEmojis := filterEmojis(emojis, query)
	if len(filteredEmojis) > 0 {
		outputEmoji(filteredEmojis[0])
	}
}

func main() {
	flag.Usage = usageFunc
	flag.Parse()
	validateFlags()

	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, true).
		SetFixed(0, 0)

	emojiNameField := tview.NewInputField()

	initialQuery := strings.Join(flag.Args(), " ")
	inputField := tview.NewInputField().
		SetDoneFunc(func(key tcell.Key) {
			app.SetFocus(table)
			table.SetSelectable(true, true)
		}).SetText(initialQuery)

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
		AddItem(table, 2, 0, 9, 3, 0, 0, false).
		AddItem(emojiNameField, 11, 0, 1, 3, 0, 0, false)
	grid.SetBorder(true).SetRect(0, 0, 60, 25)
	grid.SetTitle("Emoji Picker")

	if *noninteractiveFlag {
		runNoninterativeMode(emojis, initialQuery)
		return
	}
	drawEmojis(table, emojis, initialQuery)

	inputField.SetChangedFunc(func(text string) {
		table.Clear()
		app.Draw()
		drawEmojis(table, emojis, text)
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		app.Stop()
		outputEmoji(cell.Text)
	}).SetSelectable(false, false)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := table.GetSelection()
		if event.Key() == tcell.KeyUp && row == 0 {
			app.SetFocus(inputField)
			table.SetSelectable(false, false)
		} else if event.Key() == tcell.KeyRune {
			inputField.SetText(inputField.GetText() + string(event.Rune()))
			app.SetFocus(inputField)
			table.SetSelectable(false, false)
		}
		return event
	})

	table.SetSelectionChangedFunc(func(row, col int) {
		if r, c := table.GetSelectable(); !r && !c {
			emojiNameField.SetText("")
			return
		}
		cell := table.GetCell(row, col)
		if name, ok := getEmojiName(cell.Text); ok {
			emojiNameField.SetText(name)
		}
	})

	if err := app.SetRoot(grid, false).Run(); err != nil {
		panic(err)
	}
}
