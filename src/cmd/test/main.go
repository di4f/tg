package main

import (
    "log"
    "os"

    //tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "boteval/src/behx"
)

var rootKbd = behx.NewKeyboard(
	behx.NewButtonRow(
		behx.NewButton(
			"1",
			behx.NewCustomAction(func(c *behx.Context){
				log.Println("pressed the button!")
			}),
		),
		behx.NewButton("PRESS ME 2", behx.NewCustomAction(func(c *behx.Context){
			log.Println("pressed another button!")
		})),
	),
	behx.NewButtonRow(
		behx.NewButton("To second screen",  behx.NewScreenChange("second")),
	),
)

var secondKbd = behx.NewKeyboard(
	behx.NewButtonRow(
		behx.NewButton(
			"❤",
			behx.NewScreenChange("start"),
		),
	),
)

var inlineKbd = behx.NewKeyboard(
	behx.NewButtonRow(
		behx.NewButton(
			"INLINE PRESS ME",
			behx.NewCustomAction(func(c *behx.Context){
				log.Println("INLINE pressed the button!")
			}),
		),
		behx.NewButton("INLINE PRESS ME 2", behx.NewCustomAction(func(c *behx.Context){
			log.Println("INLINE pressed another button!")
		})),
	),
	behx.NewButtonRow(
		behx.NewButton(
			"INLINE PRESS ME 3",
			behx.ScreenChange("second"),
		),
	),
)


var startScreen = behx.NewScreen(
	"Hello, World!",
	"inline",
	"root",
)

var secondScreen = behx.NewScreen(
	"Second screen!",
	"",
	"second",
)

var behaviour = behx.NewBehaviour(
	behx.NewScreenChange("start"),
	behx.ScreenMap{
		"start": startScreen,
		"second": secondScreen,
	},
	behx.KeyboardMap{
		"root": rootKbd,
		"inline": inlineKbd,
		"second": secondKbd,
	},
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	
    bot, err := behx.NewBot(token, behaviour, nil)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)
    bot.Run()
}

