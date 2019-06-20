package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Names struct {
	Id           int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Projectname  string `gorm:"not null" form:"PROJECTNAME" json:"PROJECTNAME"`
	Universe     string `gorm:"not null" form:"UNIVERSE" json:"UNIVERSE"`
	Creationdate string `gorm:"not null" form:"CREATIONDATE" json:"CREATIONDATE"`
	Image        string `form:"IMAGE" json:"IMAGE"`
}

func main() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	// Cors setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "HEAD", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	api := router.Group("api/")
	{
		api.POST("/names", PostNames)
		api.GET("/names", GetNames)
		api.GET("/names/:id", GetName)
		api.PUT("/names/:id", UpdateName)
		api.DELETE("/names/:id", DeleteName)
	}

	router.Run(":8080")
}

func PostNames(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var names Names
	c.Bind(&names)

	if names.Projectname != "" && names.Universe != "" {
		db.Create(&names)
		// Display error
		c.JSON(201, gin.H{"success": names})
	} else {
		// Display error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"Projectname\": \"Adrasteia\", \"universe\": \"Battlestar Galactica\" , \"Creationdate\": \"2019-06-19\", \"Image\": \"shipimage.jpg\" }" http://localhost:8080/api/names
}

func GetNames(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	var names []Names
	// SELECT * FROM names
	db.Find(&names)

	// Display JSON result
	c.JSON(200, names)

	// curl -i http://localhost:8080/api/names
}

func GetName(c *gin.Context) {
	db := InitDb()
	// Close connection database
	defer db.Close()

	id := c.Params.ByName("id")
	var name Names
	// SELECT * FROM names WHERE id = 1;
	db.First(&name, id)

	if name.Id != 0 {
		// Display JSON result
		c.JSON(200, name)
	} else {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Not found"})
	}
	// curl -i http://localhost:8080/api/names/1
	// Change one to the ID you are wanting returned
}

func UpdateName(c *gin.Context) {
	// Connection to the database
	db := InitDb()
	// Close connection database
	defer db.Close()

	// Get id name
	id := c.Params.ByName("id")
	var name Names
	// SELECT * FROM names WHERE id = Entry;
	db.First(&name, id)

	if name.Projectname != "" && name.Universe != "" {

		if name.Id != 0 {
			var newName Names
			c.Bind(&newName)

			updates := Names{
				Id:           name.Id,
				Projectname:  newName.Projectname,
				Universe:     newName.Universe,
				Creationdate: newName.Creationdate,
				Image:        newName.Image,
			}

			// UPDATE names SET projectname='newName.Projectname', universe='newName.Universe' WHERE id = names.Id;
			db.Save(&updates)
			// Display modified data in JSON message "success"
			c.JSON(200, gin.H{"success": updates})
		} else {
			// Display JSON error
			c.JSON(404, gin.H{"error": "Name not found"})
		}

	} else {
		// Display JSON error
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"Projectname\": \"Cylons\", \"Universe\": \"BattleStar Galactica\", \"Creationdate\": \"2019-06-19\", \"Image\": \"\" }" http://localhost:8080/api/names/2
}

func DeleteName(c *gin.Context) {
	db := InitDb()
	id := c.Params.ByName("id")
	var name Names
	d := db.Where("id = ?", id).Delete(&name)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})

	// curl -i -X  DELETE http://localhost:8080/api/names/1
}

func InitDb() *gorm.DB {
	// Openning file
	db, err := gorm.Open("sqlite3", "./gorm.db")
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&Names{}) {
		db.CreateTable(&Names{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Names{})
	}

	return db
}
