package handlers

import (
	"fmt"
	"github.com/dharun/poc/internals/controllers"
	"github.com/gofiber/fiber/v2"
)

func GetPlanner(c *fiber.Ctx) error {
	data, err := controllers.GetPlannerDataV1(c)
	if err != nil {
		fmt.Println("Error in plannerV2 handlers")
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"msg":  "Error da macha",
			"data": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":  "intha vachuko",
		"data": data,
	})

}
