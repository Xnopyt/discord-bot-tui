package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
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
var loadUI tui.UI
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
	t.SetStyle("label.blue", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorBlue})
}

func callClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}

func main() {
	callClear()
	tokenFlag := flag.String("token", "", "Discord Bot Token")
	flag.Parse()
	tokenFlags := *tokenFlag
	var token string
	logol := tui.NewLabel("        # #                   # #\n      # #     # # # # # # #     # #\n    # # # # # # # # # # # # # # # # #\n    # # # # # # # # # # # # # # # # #\n    # # # # # # # # # # # # # # # # #\n  # # # # # # # # # # # # # # # # # # #\n  # # # # # # # # # # # # # # # # # # #\n  # # # # #     # # # # #     # # # # #\n  # # # #         # # #         # # # #\n# # # # #         # # #         # # # # #\n# # # # # #     # # # # #     # # # # # #\n# # # # # # # # # # # # # # # # # # # # #\n# # # # # # # # # # # # # # # # # # # # #\n# # # # #     # # # # # # #     # # # # #\n    # # # #                   # # # #\n      # # # #               # # # #\n\n")
	logol.SetStyleName("blue")
	titlel := tui.NewLabel("Discord Bot TUI")
	titlel.SetStyleName("magenta")
	authorl := tui.NewLabel("By Xnopyt\n")
	authorl.SetStyleName("red")
	authorBox := tui.NewPadder(36, 0, authorl)
	logoBox := tui.NewPadder(20, 0, logol)
	titleBox := tui.NewPadder(33, 0, titlel)
	logo := tui.NewVBox(
		titleBox,
		authorBox,
		logoBox,
		tui.NewSpacer(),
	)
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.SetEchoMode(tui.EchoModePassword)

	tokenText := tui.NewLabel("Enter Token > ")
	tokenText.SetStyleName("cyan")
	tokenText.SetSizePolicy(tui.Minimum, tui.Minimum)

	inputBox := tui.NewHBox(tokenText, input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	menuA := tui.NewVBox(
		tui.NewSpacer(),
		logo,
		inputBox,
		tui.NewSpacer(),
	)

	menu := tui.NewHBox(
		tui.NewSpacer(),
		menuA,
		tui.NewSpacer(),
	)

	mui, err := tui.New(menu)
	if err != nil {
		log.Fatal(err)
	}

	input.OnSubmit(func(e *tui.Entry) {
		token = e.Text()
		mui.Quit()
	})

	mui.SetTheme(t)

	mui.SetKeybinding("Esc", func() {
		mui.Quit()
		callClear()
		os.Exit(0)
	})

	if tokenFlags == "" {
		if err := mui.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		token = tokenFlags
	}
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

func appendToHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	var cname string
	var ctime string
	var err error
	var member *discordgo.Member
	member = nil
	m.Content, err = m.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		m.Content = m.ContentWithMentionsReplaced()
		err = nil
	}
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
		y, m, d := times.Date()
		cy, cm, cd := time.Now().Date()
		im := int(m)
		icm := int(cm)
		if y != cy || im != icm || d != cd {
			ctime = strconv.Itoa(d) + "/" + strconv.Itoa(im) + "/" + strconv.Itoa(y)[2:]
		}
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
	if running {
		ui.Repaint()
	}
}

func recvMsg(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !running {
		return
	}
	if m.ChannelID == cchan {
		appendToHistory(s, m)
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
	l1 := tui.NewLabel("Discord Bot TUI")
	l1.SetStyleName("magenta")
	l2 := tui.NewLabel("By Xnopyt\n\n")
	l2.SetStyleName("red")
	h1 := tui.NewHBox(
		tui.NewSpacer(),
		l1,
		tui.NewSpacer(),
	)
	h2 := tui.NewHBox(
		tui.NewSpacer(),
		l2,
		tui.NewSpacer(),
	)
	l3 := tui.NewLabel(s.State.User.Username + "#" + s.State.User.Discriminator)
	l3.SetStyleName("cyan")
	var l4 *tui.Label
	var l5 *tui.Label
	if cguild == "DM" {
		l4 = tui.NewLabel("\nDirect Message\n")
		l5 = tui.NewLabel("\nUser:\n" + " " + user.Username + "#" + user.Discriminator)

	} else {
		if len(guild.Name) > 19 {
			l4 = tui.NewLabel("\nServer:\n" + insertInto(guild.Name, 19, '\n'))
		} else {
			l4 = tui.NewLabel("\nServer:\n" + guild.Name)

		}
		if len(channel.Name) > 19 {
			l5 = tui.NewLabel("\nChannel:\n" + insertInto(channel.Name, 19, '\n'))
		} else {
			l5 = tui.NewLabel("\nChannel:\n" + channel.Name)
		}

	}
	l4.SetStyleName("green")
	l5.SetStyleName("green")
	sidebar := tui.NewVBox(
		h1,
		h2,
		l3,
		l4,
		l5,
		tui.NewLabel("\n\n\nPress 'ESC' to exit"),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history = tui.NewVBox()

	msgs, _ := s.ChannelMessages(channel.ID, 100, "", "", "")

	loadl := tui.NewLabel("Processing Channel History\n\n\n")
	loadl.SetStyleName("magenta")
	loadText := tui.NewHBox(
		tui.NewSpacer(),
		loadl,
		tui.NewSpacer(),
	)
	percent := tui.NewLabel("0%")
	percent.SetStyleName("cyan")
	percentText := tui.NewHBox(
		tui.NewSpacer(),
		percent,
		tui.NewSpacer(),
	)
	progress := tui.NewProgress(len(msgs))
	progress.SetCurrent(0)
	loadBox := tui.NewVBox(
		tui.NewSpacer(),
		loadText,
		percentText,
		progress,
		tui.NewSpacer(),
	)

	loadUI, err = tui.New(loadBox)
	if err != nil {
		log.Fatal(err)
	}
	loadUI.SetTheme(t)
	go func() {
		if err := loadUI.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	memberCache = []*discordgo.Member{}
	x := 1
	for _, v := range msgs {
		percent.SetText(strconv.Itoa(int(math.Floor(float64(x)/float64(len(msgs))*float64(100)))) + "%")
		progress.SetCurrent(x)
		loadUI.Repaint()
		x++
		appendToHistory(s, &discordgo.MessageCreate{Message: v})
	}
	time.Sleep(time.Second)
	loadUI.Quit()

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
