# 🧬 **golcli**
**A basic CLI implementation of Conway's Game of Life.**

---

## 🌱 **Setup**
Download and compile from sources:
```
go get github.com/branogarbo/golcli
```
Install just the binary with Go:
```
go install github.com/branogarbo/golcli@latest
```

Or get the pre-compiled binaries for your platform on the [releases page](https://github.com/branogarbo/golcli/releases)


## 🌳 **CLI usage**
```
golcli

A basic CLI implementation of Conway's Game of Life.

Usage:
  golcli [flags]

Examples:
golcli -c 500 -l "##" -d "--" ./pattern.txt

Flags:
  -c, --count int          The number of frames displayed before exiting (-1 : infinite loop) (default -1)
  -d, --dead-char string   The character(s) that represent a live cell (default "  ")
  -H, --height int         The height of the frames (default 30)
  -h, --help               help for golcli
  -i, --interval int       The number of milliseconds between frames (default 50)
  -l, --live-char string   The character(s) that represent a live cell (default "██")
  -x, --pattern-x int      The x offset of the initial pattern (default 12)
  -y, --pattern-y int      The y offset of the initial pattern (default 8)
  -W, --width int          The width of the frames (default 40)
```