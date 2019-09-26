package main

import (
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func main() {
	s := loginMenu()
	defer s.Close()
	for {
		run(s)
	}
}

func run(s *discordgo.Session) {
	running = false
	cguild = ""
	cchan = ""
	text, guilds := serverMenu(s)
	if text == "q" {
		callClear()
		s.Close()
		os.Exit(0)
	}
	if text == "d" {
		text, users := dmMenu(s, guilds)
		if text == "b" {
			return
		}
		if text == "q" {
			callClear()
			s.Close()
			os.Exit(0)
		}
		selc, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal("Invalid Selection")
		}
		if selc > len(users) || selc < 1 {
			log.Fatal("Invalid Selection")
		}
		selc = selc - 1
		user = users[selc]
		channel, err = s.UserChannelCreate(user.ID)
		if err != nil {
			log.Fatal(err)
		}
		cguild = "DM"
		cchan = channel.ID
	} else {
		selc, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal("Invalid Selection")
		}
		if selc > len(guilds) || selc < 1 {
			log.Fatal("Invalid Selection")
		}
		selc = selc - 1
		guild = guilds[selc]
		text, txtChannels := channelMenu(s)
		if text == "b" {
			return
		}
		if text == "q" {
			callClear()
			s.Close()
			os.Exit(0)
		}
		if text == "c" {
			nicknameMenu(s)
			return
		}
		selc, err = strconv.Atoi(text)
		if err != nil {
			log.Fatal("Invalid Selection")
		}
		if selc > len(txtChannels) || selc < 1 {
			log.Fatal("Invalid Selection")
		}
		selc = selc - 1
		channel = txtChannels[selc]
		cguild = guild.ID
		cchan = channel.ID
	}
	chatHandler(s)
}
