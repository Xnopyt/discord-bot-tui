package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/marcusolsson/tui-go"
)

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
		times = times.Local()
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

func processChannelHistory(s *discordgo.Session) *tui.Box {
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
	progressPad := tui.NewPadder(2, 0, progress)
	progressBox := tui.NewHBox(
		tui.NewSpacer(),
		progressPad,
		tui.NewSpacer(),
	)
	loadBox := tui.NewVBox(
		tui.NewSpacer(),
		loadText,
		percentText,
		progressBox,
		tui.NewSpacer(),
	)

	loadUI, err := tui.New(loadBox)
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
	return history
}

func chatHandler(s *discordgo.Session) {
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
	var l3 *tui.Label
	if len(s.State.User.Username) > 14 {
		l3 = tui.NewLabel(insertInto(s.State.User.Username+"#"+s.State.User.Discriminator, 19, '\n'))
	} else {
		l3 = tui.NewLabel(s.State.User.Username + "#" + s.State.User.Discriminator)
	}
	l3.SetStyleName("cyan")
	var l4 *tui.Label
	var l5 *tui.Label
	if cguild == "DM" {
		l4 = tui.NewLabel("\nDirect Message\n")
		if len(user.Username) > 14 {
			l5 = tui.NewLabel("\nUser:\n" + " " + insertInto(user.Username+"#"+user.Discriminator, 19, '\n'))
		} else {
			l5 = tui.NewLabel("\nUser:\n" + " " + user.Username + "#" + user.Discriminator)
		}

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

	processChannelHistory(s)

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
		} else if e.Text() != ">" {
			go func(id string) {
				if typing {
					return
				}
				typing = true
				s.ChannelTyping(id)
				time.Sleep(3 * time.Second)
				typing = false
			}(channel.ID)
		}
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
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
