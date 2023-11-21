package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	//middleware
	app.Use((func(c *fiber.Ctx) error {
		c.Locals("name", "john")
		fmt.Println("Before middleware")
		err := c.Next()
		fmt.Println("After middleware")
		return err
	}))

	app.Use(requestid.New())

	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		name := c.Locals("name")
		fmt.Println("Hello")
		return c.SendString(fmt.Sprintf("GET Hello, World! %v", name))
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("POST Hello, World!")
	})

	app.Get("/:name/:surname", func(c *fiber.Ctx) error {
		name := c.Params("name")
		return c.SendString("Hello, " + name + "!")
	})

	//query
	app.Get("/query", func(c *fiber.Ctx) error {
		person := Person{}
		c.QueryParser(&person)

		return c.JSON(person)
	})

	//windcard
	app.Get("/windcards/*", func(c *fiber.Ctx) error {
		windCard := c.Params("*")
		return c.SendString("Windcard: " + windCard)
	})

	//error
	app.Get("/error", func(c *fiber.Ctx) error {
		fmt.Println("error")
		return fiber.NewError(fiber.StatusNotFound, "content not found")
	})

	//group
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("version", "1")
		return c.Next()
	})
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello v1")
	})

	v2 := app.Group("/v2", func(c *fiber.Ctx) error {
		c.Set("version", "2")
		return c.Next()
	})
	v2.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello v2")
	})

	//mount
	userApp := fiber.New()
	userApp.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login")
	})

	//server
	app.Server().MaxConnsPerIP = 1
	app.Get("/server", func(c *fiber.Ctx) error {
		time.Sleep(time.Second * 30)
		return c.SendString("Server")
	})

	app.Mount("/user", userApp)

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"BaseURL":     c.BaseURL(),
			"Hostname":    c.Hostname(),
			"IP":          c.IP(),
			"IPs":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomains":  c.Subdomains(),
		})
	})

	//body
	app.Post("/body", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson %v\n", c.Is("json"))
		fmt.Println(string(c.Body()))

		person := Person{}
		err := c.BodyParser(&person)
		if err != nil {
			return err
		}
		fmt.Println(person)
		return nil
	})

	app.Post("/body2", func(c *fiber.Ctx) error {
		fmt.Printf("IsJson %v\n", c.Is("json"))
		// fmt.Println(string(c.Body()))

		data := map[string]interface{}{}
		err := c.BodyParser(&data)
		if err != nil {
			return err
		}
		fmt.Println(data)
		return nil
	})

	app.Listen(":3000")
}

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
