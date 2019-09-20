package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
)

func serverMenu(s *discordgo.Session) (string, []*discordgo.UserGuild) {
	callClear()
	fmt.Print("Logged in as ")
	color.Red.Print(s.State.User.Username)
	color.Magenta.Print("#" + s.State.User.Discriminator + "\n\n")
	color.Cyan.Println("Select a Server:")
	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range guilds {
		fmt.Println(strconv.Itoa(i+1) + ": " + v.Name)
	}
	fmt.Println("\nd: DM User")
	fmt.Println("q: Quit")
	fmt.Print("\n\n>")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text[:len(text)-1], "\r")
	return text, guilds
}

func channelMenu(s *discordgo.Session) (string, []*discordgo.Channel) {
	callClear()
	channels, err := s.GuildChannels(guild.ID)
	if err != nil {
		log.Fatal(err)
	}
	var txtChannels []*discordgo.Channel
	for _, v := range channels {
		if !(v.Type > 1) {
			txtChannels = append(txtChannels, v)
		}
	}
	fmt.Print("Logged in as ")
	color.Red.Print(s.State.User.Username)
	color.Magenta.Println("#" + s.State.User.Discriminator)
	fmt.Print("Server: ")
	color.Green.Print(guild.Name, "\n\n")
	color.Cyan.Println("Select a Channel:")
	for i, v := range txtChannels {
		fmt.Println(strconv.Itoa(i+1) + ": " + v.Name)
	}
	fmt.Print("\n\nc: Change Nickname\n\n>")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text[:len(text)-1], "\r")
	return text, txtChannels
}

func nicknameMenu(s *discordgo.Session) {
	callClear()
	fmt.Print("Logged in as ")
	color.Red.Print(s.State.User.Username)
	color.Magenta.Print("#" + s.State.User.Discriminator + "\n\n")
	fmt.Print("Server: ")
	color.Green.Println(guild.Name)
	member, err := s.GuildMember(guild.ID, s.State.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Current nickname: ")
	color.Cyan.Print(member.Nick + "\n\n")
	color.Magenta.Print("Enter New nickname:\n >")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text[:len(text)-1], "\r")
	s.GuildMemberNickname(guild.ID, "@me", text)
	return
}

func dmMenu(s *discordgo.Session, guilds []*discordgo.UserGuild) (string, []*discordgo.User) {
	callClear()
	color.Magenta.Println("Now Loading...")
	var users []*discordgo.User
	for _, v := range guilds {
		g, err := s.Guild(v.ID)
		if err == nil {
			for _, x := range g.Members {
				if !x.User.Bot {
					z := false
					for _, y := range users {
						if y.ID == x.User.ID {
							z = true
						}
					}
					if z {
						continue
					}
					users = append(users, x.User)
				}
			}
		}
	}
	callClear()
	fmt.Print("Logged in as ")
	color.Red.Print(s.State.User.Username)
	color.Magenta.Println("#" + s.State.User.Discriminator + "\n\n")
	color.Cyan.Println("Select a User:")
	for i, v := range users {
		fmt.Println(strconv.Itoa(i+1) + ": " + v.Username + "#" + v.Discriminator)
	}
	fmt.Print("\n\n>")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text[:len(text)-1], "\r")
	return text, users
}
