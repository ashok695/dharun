package handlers

import (
	"fmt"

	"github.com/dharun/poc/internals/controllers"
	"github.com/gofiber/fiber/v2"
)

func GetWorkbookData(c *fiber.Ctx) error {
	data, err := controllers.GetWorkbook(c)
	if err != nil {
		fmt.Println("Error in workbook handlers", err.Error())
		return c.JSON("error")
	}
	return c.JSON(fiber.Map{
		"msg":  "hi",
		"data": data,
	})
}
