package main

import (
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	// Create a new i18n bundle with default language.
	bundle := i18n.NewBundle(language.English)

	// Register a toml unmarshal function for i18n bundle.
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load translations from toml files for non-default languages.
	bundle.MustLoadMessageFile("./lang/active.es.toml")
	bundle.MustLoadMessageFile("./lang/active.ru.toml")

	// Create a new engine by passing the template folder
	// and template extension.
	engine := html.New("./templates", ".html")

	// Reload the templates on each render, good for development.
	// Optional, default is false.
	engine.Reload(true)

	// After you created your engine, you can pass it
	// to Fiber's Views Engine.
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Register a new route.
	app.Get("/", func(c *fiber.Ctx) error {
		lang := c.Query("lang")            // parse language from query
		accept := c.Get("Accept-Language") // or, parse from Header

		// Create a new localizer.
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		// Set title message.
		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",     // set translation ID
				Other: "Hello {{.Name}}", // set default translation
			},
			TemplateData: &fiber.Map{
				"Name": "John",
			},
		})

		// Parse and set unread count of emails.
		unreadEmailCount, _ := strconv.ParseInt(c.Query("unread"), 10, 64)

		// Config for translation of email count.
		unreadEmailConfig := &i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "MyUnreadEmails",
				One:   "You have {{.PluralCount}} unread email.",
				Other: "You have {{.PluralCount}} unread emails.",
			},
			PluralCount: unreadEmailCount,
		}

		// Set localizer for unread emails.
		unreadEmails := localizer.MustLocalize(unreadEmailConfig)

		// Return data as JSON.
		if c.Query("format") == "json" {
			return c.JSON(&fiber.Map{
				"name":          helloPerson,
				"unread_emails": unreadEmails,
			})
		}

		// Return rendered template.
		return c.Render("index", fiber.Map{
			"Title":        helloPerson,
			"UnreadEmails": unreadEmails,
		})
	})

	// Start server on port 3000.
	log.Fatal(app.Listen(":3000"))
}
