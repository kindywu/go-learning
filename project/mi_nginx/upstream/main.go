package main

import (
	"log"
	"main/project/mi_nginx/common"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config, err := common.ReadConfig("../config.yml")
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Config %v \n", config)

	startUpstream(config)

}

func startUpstream(config *common.Config) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!" + time.Now().Format("2006-01-02 15:04:05"))
	})

	log.Printf("Start upstream service")

	log.Fatal(app.Listen(config.UpstreamAddr))
}
