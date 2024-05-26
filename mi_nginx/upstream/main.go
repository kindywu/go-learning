package main

import (
	"log"
	"main/mi_nginx/common"

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
		return c.SendString("Hello, World!")
	})

	log.Printf("Start upstream service")

	log.Fatal(app.Listen(config.UpstreamAddr))
}
