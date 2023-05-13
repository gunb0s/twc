package telegrams

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"telegramInsiderBot/db"
	"telegramInsiderBot/insiders"
	"time"
)

func ExecuteBot() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		panic("TELEGRAM_TOKEN environment variable is empty")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: ext.DefaultMaxRoutines,
		}),
	})
	dispatcher := updater.Dispatcher

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	// /source command to send the bot source code
	dispatcher.AddHandler(handlers.NewCommand("crawl", insiderTrades))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func insiderTrades(b *gotgbot.Bot, ctx *ext.Context) error {
	var iths []insiders.InsiderTableHeader
	db.DB.Where("filing_date LIKE ?", "2023-05-12%").Find(&iths)

	today := time.Now().Format("2006-01-02")

	templateStr := `
		<b>
			<i>Today {{.Today}} - Insider Sales $100K++</i>
			{{range .Iths}}
			<b>
				<b>{{.InsiderName}}({{.Title}}) - {{.CompanyName}}({{.Ticker}}, {{.Price}})</b>
				<u>{{.Value}}</u>
			</b>
			{{end}}
		</b>
	`

	tmpl := template.Must(template.New("myTemplate").Parse(templateStr))
	data := map[string]interface{}{
		"Today": today,
		"Iths":  iths,
	}

	var result strings.Builder

	err := tmpl.ExecuteTemplate(&result, "myTemplate", data)
	if err != nil {
		panic(err)
	}

	_, err = ctx.EffectiveMessage.Reply(b, result.String(), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	return nil
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s. I <b>repeat</b> all your messages.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
