# Wallt

> This is a work in progress. Use at your own risk.

Wallt is a wallpaper changer written in Go. It is cabable of changing wallpaper on most GNU/Linux DE's and WM's (feh needed) at specific time, at specific interval of time or it divides the number of pictures by the day time.

## Installation
To build from the source:
    `go get -u github.com/fiplox/wallt`

## Usage
`wallt -h` is not implemented yet!

To run wallt with given interval, run:

`wallt /path/to/wallpapers -i HH:mm`

or

`wallt /path/to/wallpapers --set-interval HH:mm`

To run wallt at predefined time, run:

`wallt /path/to/wallpapers -t HH HH HH` as many as you wish

or

`wallt /path/to/wallpapers --set-time HH HH`

And if you run without options, like so:

`wallt /path/to/wallpapers` it will calculate the delay of changing wallpaper by the number of pictures in the directory.
