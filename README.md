<div align="right">

</div>

<div align="center">

# ▗( ◕ ̬̫ ◕ )▖ puing

`ping` command implementation in Go but with colorful output and puing ascii art


![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)

</div>

<div align="center">

<img src="https://github.com/taris-samusnow/puing/blob/1387e1848abe34a1d4addf793fa3430b14ea22c0/images/alice-puing.png" width="70%" alt="screenshot" />

</div>

## Features
- [x] Colorful and fun output.
- [x] support platform : Windows
- [x] It works with a single executable file, so it can be installed easily.
- [x] Surpports IPv4 and IPv6.

## Usage

Simply specify the target host name or IP address in the first argument e.g. `puing github.com` or `puing 13.114.40.48`.
You can change the number of transmissions by specifying the `-c` option.

```
Usage:
  puing [OPTIONS] HOST

`ping` command but with puing

Application Options:
  -c, --count=     Stop after <count> replies (default: 20)
  -P, --privilege  Enable privileged mode
  -V, --version    Show version

Help Options:
  -h, --help       Show this help message
```

## Installation

### Download executable binaries

You can download executable binaries from the latest release page.

> [![Latest Release](https://img.shields.io/github/v/release/sheepla/puing?style=flat-square)](https://github.com/sheepla/puing/releases/latest)

### Build from source

To build from source, clone this repository then run `make build`. Develo*ping* on `go1.19 windows/amd64`.


## LICENSE

[MIT](./LICENSE)

## Author

[Taris](https://github.com/taris-samusnow)