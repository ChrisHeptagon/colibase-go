package server

import (
	"database/sql"
	"fmt"

	"github.com/ChrisHeptagon/colibase/models"
	"github.com/ChrisHeptagon/colibase/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func MainServer(db *sql.DB) {
	app := fiber.New()
	app.Use(compress.New())
	app.Static("/assets", "./admin-ui/dist/assets")
	loginSchema(app)
	handleUserInitializatonStatus(app, db, "users")
	handleUserInitializaton(app, db)
	handleUserLogin(app, db)
	app.Static("/admin-ui/*", "./admin-ui/dist/index.html")
	app.Listen(":6700")
}

func handleUserLogin(a *fiber.App, db *sql.DB) error {
	a.Post("/api/login", func(c *fiber.Ctx) error {
		var userSchema models.UserSchema
		var formData map[string]interface{}
		err := models.GenerateSchema("./config.json", &userSchema)
		if err := c.BodyParser(&formData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for _, field := range userSchema.User.Fields {
			value, exists := formData[field.Name]
			if exists && field.Name == "Password" || field.Name == "password" {
				hashedPassword, err := utils.HashPassword(value.(string))
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				fmt.Printf("%s: %s\n", field.Name, hashedPassword)
				formData[field.Name] = hashedPassword
			} else if exists && field.Name != "Password" {
				fmt.Printf("%s: %s\n", field.Name, value)
			} else {
				formData[field.Name] = "[Not Provided]"
				fmt.Printf("%s: [Not Provided]\n", field.Name)
			}
		}
		fmt.Printf("Modified formData: %v\n", formData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		structFormData, err := models.MapToStruct(formData)
		fmt.Println("converted formData to struct:", structFormData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		query1 := models.QueryAdminUserFromDB("users", structFormData)
		fmt.Println("query uno:", query1)
		result, err := db.Query(query1)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for result.Next() {
			column, err := result.Columns()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			for i := range column {
				columnInterface := make([]interface{}, len(column))
				columnInterface[i] = &column[i]
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}

				fmt.Println(column[i], ":", value)

			}
		}
		return c.SendString("ok")
	})
	return nil
}

func handleUserInitializaton(a *fiber.App, db *sql.DB) error {
	a.Post("/api/init-login", func(c *fiber.Ctx) error {
		var userSchema models.UserSchema
		var formData map[string]interface{}
		err := models.GenerateSchema("./config.json", &userSchema)
		if err := c.BodyParser(&formData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for _, field := range userSchema.User.Fields {
			value, exists := formData[field.Name]
			if exists && field.Name == "Password" || field.Name == "password" {
				hashedPassword, err := utils.HashPassword(value.(string))
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				fmt.Printf("%s: %s\n", field.Name, hashedPassword)
				formData[field.Name] = hashedPassword
			} else if exists && field.Name != "Password" {
				fmt.Printf("%s: %s\n", field.Name, value)
			} else {
				formData[field.Name] = "[Not Provided]"
				fmt.Printf("%s: [Not Provided]\n", field.Name)
			}
		}
		fmt.Printf("Modified formData: %v\n", formData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		structFormData, err := models.MapToStruct(formData)
		fmt.Println("converted formData to struct:", structFormData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		query1 := models.GeneratePostgreSQLTable("users", structFormData)
		result, err := db.Exec(query1)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		r1, err := result.RowsAffected()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println("rows affected:", r1)
		result, err = db.Exec(models.InsertDataFromStruct("users", structFormData))
		fmt.Println("query dos:", models.InsertDataFromStruct("users", structFormData))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		r2, err := result.RowsAffected()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println("rows affected:", r2)
		return c.SendString(fmt.Sprintf("rows affected: %d", r2+r1))
	})
	return nil
}

func handleUserInitializatonStatus(a *fiber.App, db *sql.DB, tn string) error {
	a.Get("/api/user-initialization-status", func(c *fiber.Ctx) error {
		if models.IsUserInitialized(db) {
			return c.SendStatus(fiber.StatusOK)
		} else {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	})
	return nil
}

func loginSchema(a *fiber.App) error {
	a.Get("/api/login-schema", func(c *fiber.Ctx) error {
		var userSchema models.UserSchema
		err := models.GenerateSchema("./config.json", &userSchema)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(userSchema)
	})
	return nil
}
