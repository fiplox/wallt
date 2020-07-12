package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

var Desktop = os.Getenv("XDG_CURRENT_DESKTOP")

var DesktopSession = os.Getenv("DESKTOP_SESSION")

// ErrUnsupportedDE is thrown when Desktop is not a supported desktop environment.
var ErrUnsupportedDE = errors.New("your desktop environment is not supported")

func main() {
	args := os.Args[1:]

	files, err := ioutil.ReadDir(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var fileNames []string
	var i int
	for _, f := range files {
		if f.Name()[0] == '.' || f.IsDir() {
			continue
		}
		fileNames = append(fileNames, f.Name())
		i++
	}

	// Divide day by number of pictures in the folder.
	if len(args) == 1 {
		var delay int = 86400 / len(files)
		h := (delay / 3600)
		m := (delay - (3600 * h)) / 60
		s := (delay - (3600 * h) - (m * 60))
		var clock [3][10]int
		clock[0][1] += h
		clock[1][1] += m
		clock[2][1] += s

		for i := 2; i < len(files); i++ {
			clock[0][i] = h + clock[0][i-1]
			clock[1][i] = m + clock[1][i-1]
			clock[2][i] = s + clock[2][i-1]
			if clock[2][i] >= 60 {
				clock[2][i] -= 60
				clock[1][i] += 1
			}
			if clock[1][i] >= 60 {
				clock[1][i] -= 60
				clock[0][i] += 1
			}
			if clock[0][i] >= 24 {
				clock[0][i] -= 24
			}

		}
		now := time.Now()
		var i int
		for I, t := range clock[0] {
			if now.Hour() > t && now.Hour() < t {
				// set wallpaper to files[i] and calculate next sleep time.
				fmt.Println(files[I].Name())
				i = I + 1
				break
			}
		}

		if i == 0 {
			// set wallpaper to files[0]
			fmt.Println(args[0] + files[i].Name())
		}

		for {
			if i >= len(files) {
				i -= len(files)
			}
			next := time.Date(now.Year(), now.Month(), now.Day(), clock[0][i], clock[1][i], clock[2][i], 0, now.Location())
			//set wallpaper to files[i].
			fmt.Println(time.Until(next))
			i++
			time.Sleep(time.Until(next))
		}
	}

	// Set time manually.
	if len(args) > 1 {
		if args[1] == "-t" || args[1] == "--set-time" {

			times := make([]int, len(args[2:]))
			for i, t := range args[2:] {
				times[i], err = strconv.Atoi(t)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			now := time.Now()
			var i int
			for I, t := range times {
				if now.Hour() > t && now.Hour() < t {
					// set wallpaper to files[i] and calculate next sleep time.
					fmt.Println(files[I].Name())
					i = I + 1
					break
				}
			}
			if i == 0 {
				// set wallpaper to files[0]
				fmt.Println(args[0] + files[i].Name())
			}
			for {
				if i >= len(files) {
					i -= len(files)
				}
				next := time.Date(now.Year(), now.Month(), now.Day(), times[i], 0, 0, 0, now.Location())
				//set wallpaper to files[i].
				fmt.Println(time.Until(next))
				i++
				time.Sleep(time.Until(next))
			}
		}
		if args[1] == "-i" || args[1] == "--set-interval" {
			re := regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
			if !re.MatchString(args[2]) {
				fmt.Println(args[2], ": Wrong format.")
				os.Exit(1)
			}
			var i int
			conf, err := os.OpenFile(args[0]+".index", os.O_RDWR, 0644)
			if err != nil {
				conf, err := os.Create(args[0] + ".index")
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				conf.Write([]byte("0"))
			} else {
				_, err = fmt.Fscanf(conf, "%d", &i)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			h, _ := strconv.Atoi(string(args[2][0]) + string(args[2][1]))
			m, _ := strconv.Atoi(string(args[2][3]) + string(args[2][4]))
			m += h * 60
			for {
				if i > len(fileNames) {
					i = 0
				}
				now := time.Now()
				next := now.Add(time.Minute * time.Duration(m))
				// set wallpaper to fileNames[i]
				fmt.Println("next in:", time.Until(next), len(fileNames))
				i++
				time.Sleep(time.Until(next))
			}
		}
	}
}
