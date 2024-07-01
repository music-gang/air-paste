package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	airpaste "github.com/music-gang/air-paste"
	"github.com/music-gang/air-paste/crypto"
	"github.com/music-gang/air-paste/kv"
)

// errLogger is the logger for errors
var errLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)

func main() {

	listeningPort := os.Getenv("PORT")
	if listeningPort == "" {
		listeningPort = "8080"
	}

	e := echo.New()

	gateway := airpaste.NewGateway(kv.NewSyncedDatastore())

	airpaste.RandomString = crypto.RandomString

	e.GET("/", func(c echo.Context) error {
		return c.String(200, welcome)
	})

	e.POST("/air-copy", func(c echo.Context) error {

		value := c.FormValue("value")

		rawTTL, err := strconv.ParseInt(c.FormValue("ttl"), 10, 64)
		if err != nil {
			return handleErr(c, err)
		}

		ttl := time.Duration(rawTTL) * time.Second

		key, err := gateway.SetHandler(value, airpaste.SetOptions{
			TTL: &ttl,
		})
		if err != nil {
			return handleErr(c, err)
		}
		return c.String(200, key)
	})
	e.GET("/air-copy", func(c echo.Context) error {
		value := c.QueryParam("value")

		rawTTL, err := strconv.ParseInt(c.QueryParam("ttl"), 10, 64)
		if err != nil {
			return handleErr(c, err)
		}

		ttl := time.Duration(rawTTL) * time.Second

		key, err := gateway.SetHandler(value, airpaste.SetOptions{
			TTL: &ttl,
		})
		if err != nil {
			return handleErr(c, err)
		}
		return c.String(200, key)
	})

	e.GET("/air-paste/:key", func(c echo.Context) error {
		key := c.Param("key")
		value, ok := gateway.GetHandler(key)
		if !ok {
			return c.String(404, "Key not found")
		}
		return c.String(200, value)
	})

	e.Logger.Fatal(e.Start(":" + listeningPort))
}

func handleErr(c echo.Context, err error) error {
	msg := "And error occured while processing the request, please try again later"
	errLogger.Println(err)
	return c.String(500, msg)
}

var welcome = `Welcome to AirPaste!
Tired of write bunch of text in a tunneled ssh session or a remote desktop and copy-pasting not working? AirPaste is here to help!

How to "air copy":
- POST /air-copy with form value "value" to copy the value and get the key, optionally you can set the time to live (ttl) in seconds with form value "ttl"
	OR
- GET /air-copy with query param "value" to copy the value and get the key, optionally you can set the time to live (ttl) in seconds with query param "ttl"

How to "air paste":
- GET /air-paste/:key to paste the value with the key

Happy pasting! 🎉
`
