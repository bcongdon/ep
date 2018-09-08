# ep
⛏ Emoji Picker

`ep` is an emoji picker for the CLI.

### Demo:

![](demo.gif)

## Installation

```
go install github.com/bcongdon/ep
```

## Usage

```
ep
```

Navigation can be done with the arrow keys. Pressing `Enter` copies the selected emoji to the clipboard.

## Troubleshooting

- If you see blank squares in the emoji grid, these emojis cannot be rendered by your terminal's font.
- If you see composite emojis (i.e. `👨‍👨‍👧`) rendered as multiple emojis (i.e. `👨👨👧`), this is a [known issue](https://github.com/rivo/tview/issues/161).

## Acknowledgements

Emoji list sourced from [emojilib](https://github.com/muan/emojilib)

### Prior Art

* [Mojibar](https://github.com/muan/mojibar)