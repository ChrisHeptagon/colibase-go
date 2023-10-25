package server

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/ChrisHeptagon/colibase/models"
	"github.com/ChrisHeptagon/colibase/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func MainServer(db *sql.DB) {
	app := fiber.New()
	store := session.New()
	app.Use(compress.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: "BoF3aT8EEFv7NB+1eVpkXVNpZQ1nbqi9gP6xABEOIew=",
	}))
	publicAdminUIPaths := []string{
		"login",
		"logout",
		"init",
	}
	privateAdminUIPaths := []string{
		"dashboard",
		"settings",
		"users",
		"tables",
	}
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return handleUserLogin(c, db)
	})
	app.Get("/api/login-schema", func(c *fiber.Ctx) error {
		return loginSchema(c)
	})
	app.Get("/api/user-initialization-status", func(c *fiber.Ctx) error {
		return handleUserInitializatonStatus(c, db, "users")
	})
	app.Post("/api/init-login", func(c *fiber.Ctx) error {
		return handleUserInitializaton(c, db)
	})

	publicGroup := app.Group("/admin-entry")
	privateGroup := app.Group("/admin-ui", func(c *fiber.Ctx) error {
		return AuthMiddleware(c, store)
	})
	for _, path := range publicAdminUIPaths {
		publicGroup.Static(path, "./admin-ui/dist/index.html")
	}
	for _, path := range privateAdminUIPaths {
		privateGroup.Static(path, "./admin-ui/dist/index.html")
	}
	app.Static("/assets", "./admin-ui/dist/assets")

	app.Listen(":6700")
}

func AuthMiddleware(c *fiber.Ctx, st *session.Store) error {
	session, err := st.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	name := session.Get("name")
	cookie := c.Cookies("colibase")
	if cookie == "" {
		return c.Redirect("/admin-entry/login")
	}
	return c.Next()
}

func handleUserLogin(c *fiber.Ctx, db *sql.DB) error {
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
			formData[field.Name] = value.(string)
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
	query1 := models.QueryAdminUserDB("users", structFormData)
	fmt.Println("query uno:", query1)
	rows, err := db.Query(query1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userDeets := make(map[string]interface{})

	defer rows.Close()
	column, err := rows.Columns()
	columnsInterface := make([]interface{}, len(column))
	for rows.Next() {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for i := range column {
			columnsInterface[i] = &columnsInterface[i]
		}
		err = rows.Scan(columnsInterface...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for i := range column {
			userDeets[column[i]] = columnsInterface[i]
		}
	}
	if len(userDeets) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}
	for key, value := range userDeets {
		fmt.Printf("%s: %v\n", key, value)
		switch key {
		case regexp.MustCompile(`(?i)password`).FindString(key):
			if utils.CheckPassword(formData[key].(string), value.(string)) != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			} else if utils.CheckPassword(formData[key].(string), value.(string)) == nil {
				continue
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		case "ID", "id":
			continue
		default:
			if value != formData[key] {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		}

	}
	fmt.Println("userDeets:", userDeets)
	for key, value := range userDeets {
		fmt.Printf("%s: %v\n", key, value)
		switch key {
		case regexp.MustCompile(`(?i)email`).FindString(key):
			c.Cookie(&fiber.Cookie{
				Name:     "colibase",
				Value:    value.(string),
				Expires:  time.Now().Add(24 * 7 * time.Hour),
				SameSite: fiber.CookieSameSiteStrictMode,
				MaxAge:   24 * 7 * 60 * 60,
				HTTPOnly: true,
			})
		default:
			continue
		}
	}
	c.Redirect("/admin-ui/dashboard")
	return c.JSON(
		fiber.Map{
			"message": "Login Successful",
			"status":  fiber.StatusOK,
		},
	)
}

func handleUserInitializaton(c *fiber.Ctx, db *sql.DB) error {
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
}

func handleUserInitializatonStatus(c *fiber.Ctx, db *sql.DB, tn string) error {
	if models.IsUserInitialized(db) {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func loginSchema(c *fiber.Ctx) error {
	var userSchema models.UserSchema
	err := models.GenerateSchema("./config.json", &userSchema)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(userSchema)
}
