package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"github.com/marcusolsson/tui-go"
)

const logo = `
        # #                   # #
      # #     # # # # # # #     # #
    # # # # # # # # # # # # # # # # #
    # # # # # # # # # # # # # # # # #
    # # # # # # # # # # # # # # # # #
  # # # # # # # # # # # # # # # # # # #
  # # # # # # # # # # # # # # # # # # #
  # # # # #     # # # # #     # # # # #
  # # # #         # # #         # # # #
# # # # #         # # #         # # # # #
# # # # # #     # # # # #     # # # # # #
# # # # # # # # # # # # # # # # # # # # #
# # # # # # # # # # # # # # # # # # # # #
# # # # #     # # # # # # #     # # # # #
    # # # #                   # # # #
      # # # #               # # # #

`

func loginMenu(i int) *discordgo.Session {
	callClear()
	var tokenFlags string
	if i != 1 {
		tokenFlag := flag.String("token", "", "Discord Bot Token")
		flag.Parse()
		tokenFlags = *tokenFlag
	}
	var token string
	logol := tui.NewLabel(logo)
	logol.SetStyleName("blue")
	titlel := tui.NewLabel("Discord Bot TUI")
	titlel.SetStyleName("magenta")
	authorl := tui.NewLabel("By Xnopyt")
	authorl.SetStyleName("red")
	authorBox := tui.NewPadder(36, 0, authorl)
	logopadder := tui.NewPadder(20, 0, logol)
	titleBox := tui.NewPadder(33, 0, titlel)
	logoBox := tui.NewVBox(
		titleBox,
		authorBox,
		logopadder,
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
		logoBox,
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

	if runtime.GOOS == "windows" {
		tokenText.SetText("		  Enter Token > \nPress TAB to paste")
		mui.SetKeybinding("TAB", func() {
			clip, err := clipboard.ReadAll()
			if err != nil {
				return
			}
			input.SetText(clip)
		})
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
	<-ready
	return s
}
