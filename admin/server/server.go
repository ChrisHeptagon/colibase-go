package server

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"time"

	"encoding/json"

	"github.com/ChrisHeptagon/colibase/admin/models"
	"github.com/ChrisHeptagon/colibase/admin/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func MainServer(db *sql.DB) {
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.POST("/api/login", func(c *gin.Context) {
		handleUserLogin(c, db)
	})
	r.GET("/api/login-schema", func(c *gin.Context) {
		loginSchema(db, c)
	})
	r.POST("/api/init-login", func(c *gin.Context) {
		handleUserInitializaton(c, db)
	})
	r.POST("/api/logout", func(c *gin.Context) {
		handleUserLogout(c, db)
	})

	r.GET("/api/server-stats", func(c *gin.Context) {
		handleStats(c)
	})
	var handler gin.HandlerFunc
	if os.Getenv("MODE") == "DEV" {
		devServer, err := url.Parse(fmt.Sprintf("http://localhost:%s", os.Getenv("DEV_PORT")))
		if err != nil {
			fmt.Println("Error parsing dev server URL: ", err)
		}
		handler = func(c *gin.Context) {
			(*c).Request.Host = devServer.Host
			(*c).Request.URL.Host = devServer.Host
			(*c).Request.URL.Scheme = devServer.Scheme
			(*c).Request.RequestURI = ""

			if (*c).Request.URL.Path == "/" && (*c).Request.URL.RawQuery == "" {
				(*c).Writer.WriteHeader(http.StatusSwitchingProtocols)
				var ws websocket.Upgrader = websocket.Upgrader{
					HandshakeTimeout: 10 * time.Second,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				}
				conn, err := ws.Upgrade((*c).Writer, (*c).Request, nil)
				if err != nil {
					fmt.Println("Error upgrading websocket: ", err)
					return
				}
				defer conn.Close()
				for {
					msgT, msgB, err := conn.ReadMessage()
					if err != nil {
						fmt.Println("Error reading message: ", err)
					}
					fmt.Printf("Message Type: %d\n", msgT)
					fmt.Printf("Message: %s\n", msgB)
					err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
					if err != nil {
						fmt.Println("Error writing message: ", err)
						return
					}
				}
			}

			devServerResponse, err := http.DefaultClient.Do((*c).Request)
			if err != nil {
				fmt.Println("Error sending request to dev server: ", err)
				(*c).Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf((*c).Writer, "Error sending request to dev server: %v", err)
				return
			}
			(*c).Writer.WriteHeader(devServerResponse.StatusCode)
			(*c).Writer.Header().Set("Content-Type", devServerResponse.Header.Get("Content-Type"))
			io.Copy((*c).Writer, devServerResponse.Body)
		}
		r.GET("/src/*wildcard",
			handler,
		)
		r.GET("/@vite/client",
			handler,
		)
		r.GET("/@fs/*wildcard",
			handler,
		)
		r.GET("/node_modules/*wildcard",
			handler,
		)

	} else if os.Getenv("MODE") == "PROD" {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		r.Use(cors.New(corsConfig))
		nodeServer, err := url.Parse(fmt.Sprintf("http://localhost:%s", os.Getenv("PORT")))
		if err != nil {
			fmt.Println("Error parsing node server URL: ", err)
		}
		handler = gin.HandlerFunc(func(c *gin.Context) {
			c.Request.Host = nodeServer.Host
			c.Request.URL.Host = nodeServer.Host
			c.Request.URL.Scheme = nodeServer.Scheme
			c.Request.RequestURI = ""

			nodeServerResponse, err := http.DefaultClient.Do(c.Request)
			if err != nil {
				fmt.Println("Error sending request to node server: ", err)
				c.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(c.Writer, "Error sending request to node server: %v", err)
				return
			}
			c.Writer.WriteHeader(nodeServerResponse.StatusCode)
			c.Writer.Header().Set("Content-Type", nodeServerResponse.Header.Get("Content-Type"))
			io.Copy(c.Writer, nodeServerResponse.Body)

		})
		r.GET("/_astro/*wildcard",
			handler,
		)
	}
	entry := r.Group("/entry")
	entry.POST("/login",
		gin.HandlerFunc(func(c *gin.Context) {
			handleUserLogin(c, db)
		}),
	)
	entry.GET("/login",
		gin.HandlerFunc(func(c *gin.Context) {
			handler(c)
		}),
	)
	entry.GET("/logout",
		gin.HandlerFunc(func(c *gin.Context) {
			handleUserLogout(c, db)
		}),
	)
	r.GET("/ui/*wildcard",
		handler,
	)
	fmt.Println("Server running at http://localhost:6701")
	r.Run(":6701")
}

func handleStats(c *gin.Context) {
	stats, err := utils.GetStats()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error getting stats: %v", err)
	}
	c.JSON(http.StatusOK, stats)
}

func handleUserLogout(c *gin.Context, db *sql.DB) {
	err := c.Request.ParseForm()
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error parsing cookie: %v", err)
		return
	}
	cookieValue := c.Request.Form.Get("cookie")
	err = models.DeleteCookie(db, "sessions", cookieValue)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error deleting cookie: %v", err)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(c.Writer, "Cookie deleted")
}

func handleUserLogin(c *gin.Context, db *sql.DB) {
	c.Request.ParseMultipartForm(2048)
	formData := make(map[string]string)
	var cookieValue string
	for key, value := range c.Request.Form {
		formData[key] = value[0]
	}
	var formErrors []string
	schema, err := models.GenAdminSchema(db, "admin_schema")
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error generating schema: %v", err)
		return
	}
	for key, value := range formData {
		switch value {
		case "":
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(c.Writer, "empty field(s)")
			return
		}
		if schema[key] == nil {
			defer delete(formData, key)
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(c.Writer, "Invalid Field: %s", key)
			return
		}
		if schema[key] != nil {
			if schema[key]["required"] == "true" {
				if value == "" {
					formErrors = append(formErrors, fmt.Sprintf("Empty Field: %s", key))
				}
			}
			if schema[key]["pattern"] != "" {
				if !regexp.MustCompile(schema[key]["pattern"]).MatchString(value) {
					formErrors = append(formErrors, fmt.Sprintf("Invalid %s", key))
				}
			}
		}

	}
	if len(formErrors) > 0 {
		for key, value := range formErrors {
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(c.Writer, "Error #%d: %s\n", key+1, value)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(c.Writer, "Error(s) in form, \nplease check your form and try again")
		return
	}

	structFormData, err := models.MapToStruct(formData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error mapping form to struct: %v", err))
		return
	}
	rows, err := models.QueryAdminUserDB(db, "users", structFormData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error querying database: %v", err))
		return
	}
	userDeets := make(map[string]interface{})

	defer rows.Close()
	column, err := rows.Columns()
	columnsInterface := make([]interface{}, len(column))
	for rows.Next() {
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.WriteString(fmt.Sprintf("Error getting columns: %v", err))
		}

		for i := range column {
			columnsInterface[i] = &columnsInterface[i]
		}
		err = rows.Scan(columnsInterface...)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.WriteString(fmt.Sprintf("Error scanning columns: %v", err))
			return
		}
		for i := range column {
			userDeets[column[i]] = columnsInterface[i]
		}
	}
	if len(userDeets) == 0 {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.WriteString("Invalid Credentials")
		return
	}
	fmt.Println(userDeets)
	fmt.Println(formData)
	for key, value := range userDeets {
		if regexp.MustCompile(`(?i)password`).FindString(key) != "" {
			if utils.CheckPassword(formData[key], value.(string)) == nil {
				continue
			} else if utils.CheckPassword(formData[key], value.(string)) != nil {
				fmt.Println("Password:", formData[key])
				fmt.Println("Hashed Password:", value.(string))
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.WriteString("Invalid Credentials")
				return
			}
		} else if regexp.MustCompile(`(?i)id`).FindString(key) != "" {
			continue
		} else {
			if value != formData[key] {
				c.Writer.WriteHeader(http.StatusUnauthorized)
				c.Writer.WriteString("Invalid Credentials")
				return
			}
		}
	}

	// Convert Form Data to String
	cookieValue = string(uuid.New().String())
	resMap := map[string]interface{}{
		"message": "User Logged In",
		"value":   cookieValue,
	}
	rawJson, err := json.Marshal(resMap)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error marshaling response: %v", err))
		return
	}
	c.SetCookie("session", cookieValue, 3600, "/", c.Request.URL.Hostname(), false, true)
	c.Writer.WriteHeader(http.StatusOK)
	c.Redirect(http.StatusFound, "/ui/dashboard")
	c.Writer.Write(rawJson)
}

func handleUserInitializaton(c *gin.Context, db *sql.DB) {
	var formData map[string]string
	var cookieValue string
	if err := c.Request.ParseForm(); err != nil {
		fmt.Println("Error parsing form: ", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error parsing form: %v", err)
	}
	var formErrors []string
	schema, err := models.GenAdminSchema(db, "admin_schema")
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error generating schema: %v", err)
	}
	for key, value := range formData {
		if value == "" {
			c.Writer.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(c.Writer, "empty field(s)")
			return
		}
		if schema[key] == nil {
			delete(formData, key)

		}
		if schema[key] != nil {
			if schema[key]["required"] == "true" {
				if value == "" {
					formErrors = append(formErrors, fmt.Sprintf("Empty Field: %s", key))
				}
			}
			if schema[key]["pattern"] != "" {
				if !regexp.MustCompile(schema[key]["pattern"]).MatchString(value) {
					formData = map[string]string{
						"failure": "true",
					}
					formErrors = append(formErrors, fmt.Sprintf("Invalid %s\n", key))
				}
			}
			if regexp.MustCompile(`(?i)password`).FindString(key) != "" {
				hashedPass, err := utils.HashPassword(value)
				if err != nil {
					c.Writer.WriteHeader(http.StatusInternalServerError)
					c.Writer.WriteString(fmt.Sprintf("Error hashing password: %v", err))
				}
				formData[key] = hashedPass
			}
		}
	}
	if len(formErrors) > 0 || formData["failure"] == "true" {
		for key, value := range formErrors {
			fmt.Printf("Error #%d: %s\n", key+1, value)
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(c.Writer, "Error(s) in form, \nplease check your form and try again")
		return
	}
	structFormData, err := models.MapToStruct(formData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error mapping form to struct: %v", err))
	}
	err = models.GenerateAdminTable(db, "users", structFormData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error generating admin table: %v", err))
	}
	err = models.InsertDataFromStruct(db, "users", structFormData)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteString(fmt.Sprintf("Error inserting data into table: %v", err))
	}

	c.PureJSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "User Initialized",
			"value":   cookieValue,
		},
	)
}

func loginSchema(db *sql.DB, c *gin.Context) {
	schema, err := models.GenAdminSchema(db, "admin_schema")
	var jsonSchema []map[string]string
	for key, value := range schema {
		modVal := make(map[string]string)
		modVal["name"] = key
		for k, v := range value {
			modVal[k] = v
		}
		jsonSchema = append(jsonSchema, modVal)
	}
	sort.Slice(jsonSchema, func(i, j int) bool {
		return jsonSchema[i]["order"] < jsonSchema[j]["order"]
	})
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error generating schema: %v", err)
	}
	if schema == nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error generating schema: %v", err)
	} else if len(schema) == 0 {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Error generating schema: %v", err)
	} else {
		jsonTest, err := json.Marshal(jsonSchema)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(c.Writer, "Error marshaling schema: %v", err)
		}
		fmt.Println(string(jsonTest))

		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(jsonTest)
	}
}
