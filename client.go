package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"github.com/marcusolsson/tui-go"
)

var clear map[string]func()
var ready = make(chan bool)
var ui tui.UI
var history *tui.Box
var running = false
var cchan string
var cguild string

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
}

func callClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

func main() {
	callClear()
	color.Magenta.Print("      Discord Bot TUI - By Xnopyt\n\n")
	color.Blue.Println("         # #                   # #")
	color.Blue.Println("       # #     # # # # # # #     # #")
	color.Blue.Println("     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("     # # # # # # # # # # # # # # # # #")
	color.Blue.Println("   # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("   # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println("   # # # # #     # # # # #     # # # # #")
	color.Blue.Println("   # # # #         # # #         # # # #")
	color.Blue.Println(" # # # # #         # # #         # # # # #")
	color.Blue.Println(" # # # # # #     # # # # #     # # # # # #")
	color.Blue.Println(" # # # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println(" # # # # # # # # # # # # # # # # # # # # #")
	color.Blue.Println(" # # # # #     # # # # # # #     # # # # #")
	color.Blue.Println("     # # # #                   # # # #")
	color.Blue.Println("       # # # #               # # # #")
	fmt.Print("\n\n\n\n\n")
	reader := bufio.NewReader(os.Stdin)
	color.Cyan.Print("Enter Token: ")
	text, _ := reader.ReadString('\n')
	token := text[:len(text)-1]
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
	if m.ChannelID == cchan {
		var cname string
		var ctime string
		member, err := s.GuildMember(cguild, m.Author.ID)
		if err != nil {
			cname = m.Author.Username
		} else {
			if member.Nick == "" {
				cname = m.Author.Username
			} else {
				cname = member.Nick
			}
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
		history.Append(tui.NewHBox(
			tui.NewLabel(ctime),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", cname))),
			tui.NewLabel(m.Content),
			tui.NewSpacer(),
		))
		ui.Repaint()
	}
}

func run(s *discordgo.Session) {
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
	fmt.Println("q: Quit")
	fmt.Print("\n\n>")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-1]
	if text == "q" {
		callClear()
		os.Exit(0)
	}
	selc, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal("Invalid Selection")
	}
	if selc > len(guilds) || selc < 1 {
		log.Fatal("Invalid Selection")
	}
	selc = selc - 1
	guild := guilds[selc]
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
	text = text[:len(text)-1]
	selc, err = strconv.Atoi(text)
	if err != nil {
		log.Fatal("Invalid Selection")
	}
	if selc > len(txtChannels) || selc < 1 {
		log.Fatal("Invalid Selection")
	}
	selc = selc - 1
	channel := txtChannels[selc]
	callClear()
	color.Magenta.Println("Now Loading...")
	sidebar := tui.NewVBox(
		tui.NewLabel("Discord Bot TUI"),
		tui.NewLabel("By Xnopyt"),
		tui.NewLabel(s.State.User.Username+"#"+s.State.User.Discriminator),
		tui.NewLabel("Server: "+guild.Name),
		tui.NewLabel("Channel: "+channel.Name),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history = tui.NewVBox()

	msgs, _ := s.ChannelMessages(channel.ID, 25, "", "", "")
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	for _, v := range msgs {
		var cname string
		var ctime string
		member, err := s.GuildMember(guild.ID, v.Author.ID)
		if err != nil {
			cname = v.Author.Username
		} else {
			if member.Nick == "" {
				cname = v.Author.Username
			} else {
				cname = member.Nick
			}
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
		}
		history.Append(tui.NewHBox(
			tui.NewLabel(ctime),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", cname))),
			tui.NewLabel(v.Content),
			tui.NewSpacer(),
		))
	}

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		s.ChannelMessageSend(channel.ID, e.Text())
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err = tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() {
		ui.Quit()
		callClear()
	})

	cguild = guild.ID
	cchan = channel.ID
	running = true
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
