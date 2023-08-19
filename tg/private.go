package tg

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type context struct {
	Session *Session
	Bot     *Bot
	updates chan *Update
	// Is true if currently reading the Update.
	readingUpdate bool

	curScreen, prevScreen *Screen
}

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.

// Goroutie function to handle each user.
func (c *context) handleUpdateChan(updates chan *Update) {
	beh := c.Bot.behaviour

	if beh.Init != nil {
		c.run(beh.Init, nil)
	}
	for u := range updates {
		var act Action
		screen := c.curScreen
		// The part is added to implement custom update handling.
		if u.Message != nil {
			if u.Message.IsCommand() && !c.readingUpdate {
				cmdName := CommandName(u.Message.Command())
				cmd, ok := beh.Commands[cmdName]
				if ok {
					act = cmd.Action
				} else {
				}
			} else {
				kbd := screen.Keyboard
				if kbd == nil {
					if c.readingUpdate {
						c.updates <- u
					}
					continue
				}
				btns := kbd.buttonMap()
				text := u.Message.Text
				btn, ok := btns[text]
				if !ok {
					if u.Message.Location != nil {
						for _, b := range btns {
							if b.SendLocation {
								btn = b
								ok = true
							}
						}
					} else if c.readingUpdate {
						// Skipping the update sending it to
						// the reading goroutine.
						c.updates <- u
						continue
					}
				}

				if ok {
					act = btn.Action
				}
			}
		} else if u.CallbackQuery != nil {
			cb := tgbotapi.NewCallback(
				u.CallbackQuery.ID,
				u.CallbackQuery.Data,
			)
			data := u.CallbackQuery.Data

			_, err := c.Bot.Api.Request(cb)
			if err != nil {
				panic(err)
			}
			kbd := screen.InlineKeyboard
			if kbd == nil {
				if c.readingUpdate {
					c.updates <- u
				}
				continue
			}

			btns := kbd.buttonMap()
			btn, ok := btns[data]
			if !ok && c.readingUpdate {
				c.updates <- u
				continue
			}
			if !ok {
				c.Sendf("%q", btns)
				continue
			}
			act = btn.Action
		}
		if act != nil {
			c.run(act, u)
		}
	}
}

func (c *context) run(a Action, u *Update) {
	go a.Act(&Context{
		context: c,
		Update:  u,
	})
}

// Returns the next update ignoring current screen.
func (c *context) ReadUpdate() (*Update, error) {
	c.readingUpdate = true
	u := <-c.updates
	c.readingUpdate = false
	if u == nil {
		return nil, NotAvailableErr
	}

	return u, nil
}

// Returns the next text message that the user sends.
func (c *context) ReadTextMessage() (string, error) {
	u, err := c.ReadUpdate()
	if err != nil {
		return "", err
	}
	if u.Message == nil {
		return "", WrongUpdateType{}
	}

	return u.Message.Text, nil
}

// Sends to the user specified text.
func (c *context) Send(v any) (*Message, error) {
	return c.Bot.Send(c.Session.Id, v)
}

// Sends the formatted with fmt.Sprintf message to the user.
func (c *context) Sendf(format string, v ...any) (*Message, error) {
	msg, err := c.Send(fmt.Sprintf(format, v...))
	if err != nil {
		return nil, err
	}
	return msg, err
}

// Interface to interact with the user.
type Context struct {
	*context
	*Update
}

// Customized actions for the bot.
type Action interface {
	Act(*Context)
}

type ActionFunc func(*Context)

func (af ActionFunc) Act(c *Context) {
	af(c)
}

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Context) {
	if !c.Bot.behaviour.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

type C = Context

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
	if !c.Bot.behaviour.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Stop the reading by sending the nil,
	// since we change the screen and
	// current goroutine needs to be stopped.
	if c.readingUpdate {
		c.updates <- nil
	}

	// Getting the screen and changing to
	// then executing its action.
	screen := c.Bot.behaviour.Screens[screenId]
	c.prevScreen = c.curScreen
	c.curScreen = screen
	screen.Render(c.Session.Id, c.Bot)
	if screen.Action != nil {
		c.run(screen.Action, c.Update)
	}

	return nil
}