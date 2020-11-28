# ep

⛏ Emoji Picker

`ep` is an emoji picker for the CLI.

### Demo:

![](demo.gif)

## Installation

```
go install github.com/bcongdon/ep
```

### Via HomeBrew

```
brew tap bcongdon/ep
brew install ep
```

(Brew definition is at
[bcongdon/homebrew-ep](https://github.com/bcongdon/homebrew-ep))

## Usage

```
Usage of ./ep:
./ep [QUERY]
  -noninteractive
    	If set, doesn't display emoji picker -- instead just outputting the first selection for the provided query.
  -output string
    	The output of ep. Choices: clipboard, stdout (default "clipboard")
```

Navigation can be done with the arrow keys. Pressing `Enter` copies the selected
emoji to the clipboard.

### Examples

- `ep` - Opens the default emoji picker.
- `ep wink` - Opens the default emoji picker with the query "wink" already
  entered
- `ep -noninteractive wink` - Copies the first search result for "wink" to the
  clipboard
- `ep -noninteractive -output=stdout wink` - Outputs the first search result for
  "wink" to stdout

## Troubleshooting

- If you see blank squares in the emoji grid, these emojis cannot be rendered by
  your terminal's font.

## Acknowledgements

Emoji list sourced from [emojilib](https://github.com/muan/emojilib)

### Prior Art

- [Mojibar](https://github.com/muan/mojibar)
