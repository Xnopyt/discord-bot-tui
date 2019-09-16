package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"github.com/marcusolsson/tui-go"
)

var memberCache []*discordgo.Member
var clear map[string]func()
var ready = make(chan bool)
var ui tui.UI
var history *tui.Box
var running = false
var cchan string
var cguild string
var guild *discordgo.UserGuild
var user *discordgo.User
var channel *discordgo.Channel
var t *tui.Theme

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
}

func callClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

func main() {
	callClear()
	color.Magenta.Print("      			Discord Bot TUI - By Xnopyt\n\n")
	color.Blue.Println("		         # #                   # #")
	color.Blue.Println("		       # #     # # # # # # #     # #")
	color.Blue.Println("		     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		   # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		   # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		   # # # # #     # # # # #     # # # # #")
	color.Blue.Println("		   # # # #         # # #         # # # #")
	color.Blue.Println("		 # # # # #         # # #         # # # # #")
	color.Blue.Println("		 # # # # # #     # # # # #     # # # # # #")
	color.Blue.Println("		 # # # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		 # # # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("		 # # # # #     # # # # # # #     # # # # #")
	color.Blue.Println("		     # # # #                   # # # #")
	color.Blue.Println("		       # # # #               # # # #")
	fmt.Print("\n\n\n\n\n")
	reader := bufio.NewReader(os.Stdin)
	color.Cyan.Print("Enter Token: ")
	text, _ := reader.ReadString('\n')
	token := strings.TrimSuffix(text[:len(text)-1], "\r")
	callClear()
	color.Magenta.Println("Connecting to Discord....")
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		ready <- true
	})
	s.AddHandler(recvMsg)
	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	<-ready
	for {
		run(s)
	}
}

func recvMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !running {
		return
	}
	if m.ChannelID == cchan {
		var cname string
		var ctime string
		var err error
		var member *discordgo.Member
		member = nil
		err = nil
		if cguild != "DM" {
			for _, z := range memberCache {
				if z.User.ID == m.Author.ID {
					member = z
				}
			}
			if member == nil {
				member, err = s.GuildMember(cguild, m.Author.ID)
				if err == nil {
					memberCache = append(memberCache, member)
				}
			}
			if err != nil {
				cname = m.Author.Username
			} else {
				if member.Nick == "" {
					cname = m.Author.Username
				} else {
					cname = member.Nick
				}
			}
		} else {
			cname = m.Author.Username
		}
		times, err := m.Timestamp.Parse()
		if err != nil {
			ctime = "00:00"
		} else {
			hr, mi, _ := times.Clock()
			var min string
			if mi < 10 {
				min = strconv.Itoa(mi)
				min = "0" + min
			} else {
				min = strconv.Itoa(mi)
			}
			ctime = strconv.Itoa(hr) + ":" + min
		}
		for _, z := range m.Attachments {
			if m.Content == "" {
				m.Content = z.URL
			} else {
				m.Content += "\n" + z.URL
			}
		}
		for _, z := range m.Embeds {
			if m.Content != "" {
				m.Content += "\n" + "Embed:"
			}
			if z.Title != "" {
				m.Content += "\n" + z.Title
			}
			if z.Description != "" {
				m.Content += "\n" + z.Description
			}
			if z.URL != "" {
				m.Content += "\n" + z.URL
			}
			if z.Description != "" {
				m.Content += "\n" + z.Description
			}
			if z.Image != nil {
				m.Content += "\n" + z.Image.URL
			}
			if z.Thumbnail != nil {
				m.Content += "\n" + z.Thumbnail.URL
			}
			if z.Video != nil {
				m.Content += "\n" + z.Video.URL
			}
			for _, f := range z.Fields {
				m.Content += "\n" + f.Name + ": " + f.Value
			}
			if z.Provider != nil {
				m.Content += "\n" + "Provider: " + z.Provider.Name + " (" + z.Provider.URL + ")"
			}
			if z.Footer != nil {
				m.Content += "\n" + z.Footer.Text + " " + z.Footer.IconURL
			}
		}
		l1 := tui.NewLabel(ctime)
		l1.SetStyleName("red")
		l2 := tui.NewLabel(fmt.Sprintf("<%s>", cname))
		l2.SetStyleName("green")
		l3 := tui.NewLabel(m.Content)
		l3.SetStyleName("cyan")
		history.Append(tui.NewHBox(
			l1,
			tui.NewPadder(1, 0, l2),
			l3,
			tui.NewSpacer(),
		))
		ui.Repaint()
	}
}

func run(s *discordgo.Session) {
	running = false
	cguild = ""
	cchan = ""
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
	if text == "q" {
		callClear()
		os.Exit(0)
	}
	if text == "d" {
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
		fmt.Print("\n\n>")
		text, _ = reader.ReadString('\n')
		text = strings.TrimSuffix(text[:len(text)-1], "\r")
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
	callClear()
	color.Magenta.Println("Now Loading...")
	l1 := tui.NewLabel("  Discord Bot TUI  ")
	l1.SetStyleName("magenta")
	l2 := tui.NewLabel("    By Xnopyt\n\n")
	l2.SetStyleName("red")
	l3 := tui.NewLabel(s.State.User.Username + "#" + s.State.User.Discriminator)
	l3.SetStyleName("cyan")
	var l4 *tui.Label
	var l5 *tui.Label
	if cguild == "DM" {
		l4 = tui.NewLabel("\nDirect Message\n")
		l5 = tui.NewLabel("\nUser:\n" + " " + user.Username + "#" + user.Discriminator)

	} else {
		l4 = tui.NewLabel("\nServer:\n" + " " + guild.Name)
		l5 = tui.NewLabel("\nChannel:\n" + " " + channel.Name)

	}
	l4.SetStyleName("green")
	l5.SetStyleName("green")
	sidebar := tui.NewVBox(
		l1,
		l2,
		l3,
		l4,
		l5,
		tui.NewLabel("\n\n\nPress 'ESC' to exit"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history = tui.NewVBox()

	msgs, _ := s.ChannelMessages(channel.ID, 100, "", "", "")
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	color.Red.Printf("Processing Channel history: 0/" + strconv.Itoa(len(msgs)))
	memberCache = []*discordgo.Member{}
	x := 1
	for _, v := range msgs {
		color.Red.Printf("\rProcessing Channel history: %d/"+strconv.Itoa(len(msgs)), x)
		x++
		var cname string
		var ctime string
		var err error
		var member *discordgo.Member
		member = nil
		err = nil
		if cguild != "DM" {
			for _, z := range memberCache {
				if z.User.ID == v.Author.ID {
					member = z
				}
			}
			if member == nil {
				member, err = s.GuildMember(guild.ID, v.Author.ID)
				if err == nil {
					memberCache = append(memberCache, member)
				}
			}
			if err != nil {
				cname = v.Author.Username
			} else {
				if member.Nick == "" {
					cname = v.Author.Username
				} else {
					cname = member.Nick
				}
			}
		} else {
			cname = v.Author.Username
		}
		times, err := v.Timestamp.Parse()
		if err != nil {
			ctime = "00:00"
		} else {
			hr, mi, _ := times.Clock()
			var min string
			if mi < 10 {
				min = strconv.Itoa(mi)
				min = "0" + min
			} else {
				min = strconv.Itoa(mi)
			}
			ctime = strconv.Itoa(hr) + ":" + min
			y, m, d := times.Date()
			cy, cm, cd := time.Now().Date()
			im := int(m)
			icm := int(cm)
			if y != cy || im != icm || d != cd {
				ctime = strconv.Itoa(d) + "/" + strconv.Itoa(im) + "/" + strconv.Itoa(y)[2:]
			}
		}
		for _, z := range v.Attachments {
			if v.Content == "" {
				v.Content = z.URL
			} else {
				v.Content += "\n" + z.URL
			}
		}
		for _, z := range v.Embeds {
			if v.Content != "" {
				v.Content += "\n" + "Embed:"
			}
			if z.Title != "" {
				v.Content += "\n" + z.Title
			}
			if z.Description != "" {
				v.Content += "\n" + z.Description
			}
			if z.URL != "" {
				v.Content += "\n" + z.URL
			}
			if z.Description != "" {
				v.Content += "\n" + z.Description
			}
			if z.Image != nil {
				v.Content += "\n" + z.Image.URL
			}
			if z.Thumbnail != nil {
				v.Content += "\n" + z.Thumbnail.URL
			}
			if z.Video != nil {
				v.Content += "\n" + z.Video.URL
			}
			for _, f := range z.Fields {
				v.Content += "\n" + f.Name + ": " + f.Value
			}
			if z.Provider != nil {
				v.Content += "\n" + "Provider: " + z.Provider.Name + " (" + z.Provider.URL + ")"
			}
			if z.Footer != nil {
				v.Content += "\n" + z.Footer.Text + " " + z.Footer.IconURL
			}
		}
		l6 := tui.NewLabel(ctime)
		l6.SetStyleName("red")
		l7 := tui.NewLabel(fmt.Sprintf("<%s>", cname))
		l7.SetStyleName("green")
		l8 := tui.NewLabel(v.Content)
		l8.SetStyleName("cyan")
		history.Append(tui.NewHBox(
			l6,
			tui.NewPadder(1, 0, l7),
			l8,
			tui.NewSpacer(),
		))
	}

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetText(">")
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		s.ChannelMessageSend(channel.ID, e.Text()[1:])
		input.SetText(">")
	})

	input.OnChanged(func(e *tui.Entry) {
		if e.Text() == "" {
			input.SetText(">")
		}
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err = tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetTheme(t)

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		callClear()
	})

	running = true
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
