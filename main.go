package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/reujab/wallpaper"
)

func main() {
	var wallpaperPath, mode, HHMM string

	flag.StringVar(&wallpaperPath, "path", "", "Specify absolute path to wallpaper folder. Required.")
	flag.StringVar(&mode, "mode", "auto", "Specify mode of use. Mods available: auto (default), interval, time.")
	flag.StringVar(&HHMM, "time", "", "Specify the interval of time to change wallpaper. Required if -mode=interval.")

	flag.Parse()

	if wallpaperPath == "" || (mode == "interval" && HHMM == "") || (mode == "auto" && HHMM != "") || (mode == "time" && HHMM != "") {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if wallpaperPath[len(wallpaperPath)-1] != '/' {
		wallpaperPath += "/"
	}

	files, err := ioutil.ReadDir(wallpaperPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get file names.
	var fileNames []string
	var i int
	for _, f := range files {
		if f.Name()[0] == '.' || f.IsDir() {
			continue
		}
		fileNames = append(fileNames, f.Name())
		i++
	}

	i = 0
	conf, err := os.OpenFile(wallpaperPath+".index", os.O_RDWR, 0755)
	if err != nil {
		conf, err := os.Create(wallpaperPath + ".index")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conf.WriteString("0")
	} else {
		_, err = fmt.Fscanf(conf, "%d", &i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	conf.Close()

	switch mode {
	case "auto":
		var delay int = 86400 / len(fileNames)
		for {
			now := time.Now()
			if i >= len(fileNames) {
				i -= len(fileNames)
			}

			next := now.Add(time.Second * time.Duration(delay))
			err := wallpaper.SetFromFile(wallpaperPath + fileNames[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("next in:", time.Until(next))
			i++
			ioutil.WriteFile(wallpaperPath+".index", []byte(strconv.Itoa(i)), 644)
			time.Sleep(time.Until(next))
		}
	case "interval":
		re := regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
		if !re.MatchString(HHMM) {
			fmt.Printf("%s: Wrong format.", HHMM)
			os.Exit(1)
		}

		h, _ := strconv.Atoi(string(HHMM[0]) + string(HHMM[1]))
		m, _ := strconv.Atoi(string(HHMM[3]) + string(HHMM[4]))
		m += h * 60

		for {
			if i > len(fileNames) {
				i = 0
			}
			now := time.Now()
			next := now.Add(time.Minute * time.Duration(m))
			err := wallpaper.SetFromFile(wallpaperPath + fileNames[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("next in:", time.Until(next))
			i++
			ioutil.WriteFile(wallpaperPath+".index", []byte(strconv.Itoa(i)), 644)
			time.Sleep(time.Until(next))
		}
	case "time":
		times := make([]int, len(flag.Args()))
		for I, t := range flag.Args() {
			times[I], err = strconv.Atoi(t)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		now := time.Now()
		i = 0
		for I := 0; I < len(times)-1; I++ {
			if now.Hour() > times[I] && now.Hour() < times[I+1] {
				i = I + 1
				break
			}
		}

		for {
			now = time.Now()
			if i >= len(times) {
				i = 0
			}

			next := time.Date(now.Year(), now.Month(), now.Day(), times[i], 0, 0, 0, now.Location())
			if time.Until(next) < time.Duration(0) {
				next = next.AddDate(0, 0, 1)
			}
			err := wallpaper.SetFromFile(wallpaperPath + fileNames[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("next in:", time.Until(next))
			i++
			time.Sleep(time.Until(next))
		}
	}
}
