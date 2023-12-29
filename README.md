# Tbmk - Terminal Bookmarker

A commands bookmark for terminal

![Demo](./tbmk.gif)

## Worked on

- Linux bash
- Linux zsh
- Mac zsh

## Installation

1. Download latest version from the [release page](https://github.com/linhx/tbmk/releases).
2. Extract the file, e.g., `unzip tbmk*.zip`.
3. Run `cd tbmk-VERSION`.
4. Run `./install` (avoid using `absolute-path/install`). It will add keybindings to `~/.bashrc` and `~/.zsh`.

## Usage

1. **Search:** Type and hit `ctrl + space`.
2. **Delete:** In the result screen, select the item and press `ctrl + d`.
3. **Add:** Press `ctrl + t`. You can type the command first and then press `ctrl + t`.
4. **Edit:** Override the old one by adding a new command with the same title.

Data is stored in `~/.tbmk`. Feel free to backup or edit it directly.

**TODO:**

- [ ] Windows support

## Development

### Build

```shell
git clone https://github.com/linhx/tbmk.git
cd tbmk
go build .
./install
```
