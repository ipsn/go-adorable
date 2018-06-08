// go-adorable - Adorable Avatars from Go
// Copyright (c) 2018 Péter Szilágyi. All rights reserved.

// +build ignore

package main

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/ipsn/go-adorable"
)

// This is a sample program that generates an HTML file with a lot of embedded
// random avatars. Its sole purpose is to allow quickly samplig the generator
// for debugging purposes.
func main() {
	// Generate the sample HTML file
	html, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(html, "<html><body>")

	// Generate the random images
	fmt.Println("Generating fully random avatars...")
	fmt.Fprintf(html, "<h3>Fully random</h3>")
	for i := 0; i < 125; i++ {
		uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString(adorable.Random())
		fmt.Fprintf(html, "<img src='%s' width='64px' height='64px' style='margin: 4px'/>", uri)
	}
	fmt.Println("Generating random avatars with color...")
	fmt.Fprintf(html, "<h3>Random with color</h3>")
	for i := 0; i < 125; i++ {
		uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString(adorable.RandomWithColor(color.RGBA{R: 96, G: byte(2 * i), B: 192, A: 255}))
		fmt.Fprintf(html, "<img src='%s' width='64px' height='64px' style='margin: 4px'/>", uri)
	}
	// Generate the pseudorandom images
	fmt.Println("Generating fully pseudorandom avatars...")
	fmt.Fprintf(html, "<h3>Fully pseudorandom</h3>")
	for i := 0; i < 125; i++ {
		uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString(adorable.PseudoRandom([]byte{byte(i)}))
		fmt.Fprintf(html, "<img src='%s' width='64px' height='64px' style='margin: 4px'/>", uri)
	}
	fmt.Println("Generating pseudorandom avatars with color...")
	fmt.Fprintf(html, "<h3>Pseudorandom with color</h3>")
	for i := 0; i < 125; i++ {
		uri := "data:image/png;base64," + base64.StdEncoding.EncodeToString(adorable.PseudoRandomWithColor([]byte{byte(i)}, color.RGBA{R: 96, G: byte(2 * i), B: 192, A: 255}))
		fmt.Fprintf(html, "<img src='%s' width='64px' height='64px' style='margin: 4px'/>", uri)
	}
	fmt.Fprintf(html, "</body></html>")
	html.Close()

	// Assemble the list of commad we can use to open a browser
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	default:
		cmds = append(cmds, []string{"xdg-open"})
	}
	cmds = append(cmds,
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"chromium"},
		[]string{"firefox"},
	)
	// Open the sample HTML page, wait a bit and delete it
	for _, args := range cmds {
		cmd := exec.Command(args[0], append(args[1:], html.Name())...)
		if cmd.Start() == nil {
			time.Sleep(3 * time.Second)
			os.Remove(html.Name())
			break
		}
	}
}
