package main

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/marcusolsson/tui-go"
)

var memberCache []*discordgo.Member
var clear map[string]func()
var ready = make(chan bool)
var ui tui.UI
var loadUI tui.UI
var history *tui.Box
var running = false
var cchan string
var cguild string
var guild *discordgo.UserGuild
var user *discordgo.User
var channel *discordgo.Channel
var t *tui.Theme
var typing bool

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	t = tui.NewTheme()
	t.SetStyle("normal", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorWhite})
	t.SetStyle("label.magenta", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorMagenta})
	t.SetStyle("label.red", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})
	t.SetStyle("label.green", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorGreen})
	t.SetStyle("label.cyan", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorCyan})
	t.SetStyle("label.yellow", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorYellow})
	t.SetStyle("label.blue", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorBlue})
}

func callClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

func insertInto(s string, interval int, sep rune) string {
	var buffer bytes.Buffer
	before := interval - 1
	last := len(s) - 1
	for i, char := range s {
		buffer.WriteRune(char)
		if i%interval == before && i != last {
			buffer.WriteRune(sep)
		}
	}
	return buffer.String()
}
