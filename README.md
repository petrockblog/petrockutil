# petrockutil

[![GitHub release](https://img.shields.io/github/release/petrockblog/petroutil.svg)](https://github.com/petrockblog/petrockutil/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/petrockblog/petrockutil)](https://goreportcard.com/report/github.com/petrockblog/petrockutil) [![AUR](https://img.shields.io/aur/license/yaourt.svg)]()

Petrockutil is a command line utility for interacting with petrockblock.com devices.

Petrockutil is written in the Go programming language (http://golang.org) for maximum speed and portability.

# Getting Started

Download the binaries for your platform, unzip the archive and run the command line tool without any arguments to get usage information.

A typical series of steps would be to first

1. scan for serial ports, `petrockutil scan serial`, then
2. read the firmware version, `petrockutil gamepadblock readversion <port>`, and then
3. if you do not have AVRDude installed, install it, `petrockutil gamepadblock prepare`,
3. update the firmware, `petrockutil gamepadblock update <port>`.


# Download

We provide binaries for OSX, Windows, and Linux. You find the latest releases on our web site at [https://blog.petrockblock.com/gamepadblock-downloads/](https://blog.petrockblock.com/gamepadblock-downloads/).


# How To Use

Calling `petrockutil` without any arguments prints __usage information__ to the console:
```
 $ ./petrockutil 
NAME:
   petrockutil - Command Line Utility for petrockblock.com gadgets

USAGE:
   petrockutil [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     scan          Scan for connected devices on Serial, USB, or Bluetooth ports
     gamepadblock  Installs avrdude, read firmware version, and update firmware of your GamepadBlock
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version   
```

Usually, the first thing to do is to __look for the serial com port__ of the device that that you want to interact with:
```
$ ./petrockutil scan serial
/dev/cu.Bluetooth-Incoming-Port     /dev/tty.Bluetooth-Incoming-Port
/dev/cu.usbmodem1421            /dev/tty.usbmodem1421

```