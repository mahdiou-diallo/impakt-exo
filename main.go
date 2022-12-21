package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IntString struct {
	// regexp=^\\d+(\\,\\d+)*$
	Data string `json:"values" validate:"required"`
}

type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type StringWithIDArray struct {
	Values []StringWithID `validate:"required,dive"`
}

func main() {
	fmt.Println("starting")
	app := fiber.New()
	validate := validator.New()

	app.Post("/count", func(ctx *fiber.Ctx) error {
		var data string

		if err := ctx.BodyParser(&data); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		if err := validate.Var(data, "required"); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		count, err := CountIncreases(data)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		return ctx.Status(fiber.StatusOK).JSON(count)
	})

	app.Post("/count2", func(ctx *fiber.Ctx) error {
		var payload []StringWithID

		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		if len(payload) != 2 {
			return ctx.Status(fiber.StatusBadRequest).JSON("wrong array size")
		}

		if err := validate.Struct(StringWithIDArray{payload}); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		co := make(chan *CountWithIDAndError, 2)
		go CountIncreases2(&payload[0], co)
		go CountIncreases2(&payload[1], co)

		result1 := <-co
		result2 := <-co
		if err := result1.Err; err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}
		if err := result2.Err; err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		return ctx.Status(fiber.StatusOK).JSON([]CountWithID{*result1.Count, *result2.Count})
	})

	log.Fatal(app.Listen(":10000"))
}
