# Tbmk - Terminal bookmarker

A commands bookmark for terminal

![demo](./tbmk.gif)

## Worked on

- Linux bash
- Linux zsh

## How to install

1. Download built file on release page
2. Extract the file. e.g. /somepath/tbmk
3. Run `cd /somepath/tbmk`
4. Run `./install`, it will appends keybinding to `~/.bashrc` and `~/.zsh`

## How to use

1. Search: type and `ctrl + space`
2. Delete: in the result screen, select the item then press `ctrl + d`
3. Add: `ctrl + t`. you can type the command first then press `ctrl + t`
4. Edit: Override the old one by add new command with the same title.

The data are stored in `~/.tbmk/BookmarkItem`. You can backup or edit it directly.

TODO

- [ ] Windows
- [ ] Mac

## Develop

### Build

```shell
go build .
```
