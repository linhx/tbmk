# Tbmk - Terminal bookmarker

![](https://github.com/linhx/tbmk/actions/workflows/go.yml/badge.svg)

A commands bookmark for terminal

![demo](./tbmk.gif)

Support placeholders and multilines command

![Add command](https://github.com/user-attachments/assets/0f501fbf-2963-484e-8b9d-1510b4c087c9)

![Preview command](https://github.com/user-attachments/assets/7601267b-430e-4abc-bbd5-d9507051cdda)

![Edit placeholders](https://github.com/user-attachments/assets/d2c46031-f661-43ae-bfb9-c539d7665b10)


## Worked on

- Linux: bash, zsh, fish
- Mac zsh

## How to install

1. Download built file on release page
2. Extract the file. e.g. /somepath/tbmk
3. Run `cd /somepath/tbmk`
4. Run `./install` (don't install by execute `/absolute-path/install`), it will appends keybinding to `~/.bashrc`, `~/.zsh` and `~/.config/fish/config.fish`
5. Restart your shell or reload config file:
    - `source ~/.bashrc # bash`
    - `source ~/.zshrc # zsh`
    - `source ~/.config/fish/config.fish #fish`

## How to use

1. Search: type and `ctrl + space`
2. Delete: in the result screen, select the item then press `ctrl + d`
3. Add: `ctrl + t`. you can type the command first then press `ctrl + t`
4. Edit: Override the old one by add new command with the same title.

The data are stored in `~/.tbmk`. You can backup or edit it directly.

TODO

- [ ] Windows

## Development

### Run

```shell
APP_ENV=dev go run .
```

### Build

```shell
go build .
```
