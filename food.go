package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	SqlFood  = `
		SELECT
			id,
			name,
			brand,
			nutrition_grade,
			nutrition_score,
			url,
			image_food,
			image_nutrition,
			category_id
		FROM food_food
		WHERE id = $1;`
)

type Food struct {
	Id             int
	Name           string
	Brand          string
	NutritionGrade string
	NutritionScore int
	URL            string
	ImageFood      string
	ImageNutrition string
	CategoryId     int
}

var db *sql.DB

func ConnectDb() error {
	var err error
	config := fmt.Sprintf("host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode='disable'",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"))
	db, err = sql.Open("postgres", config)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func food(ctx *gin.Context) {
	id := ctx.Param("id")
	row := db.QueryRow(SqlFood, id)
	food := Food{}
	err := row.Scan(
		&food.Id,
		&food.Name,
		&food.Brand,
		&food.NutritionGrade,
		&food.NutritionScore,
		&food.URL,
		&food.ImageFood,
		&food.ImageNutrition,
		&food.CategoryId)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	out, err := json.Marshal(food)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, string(out))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/food/:id", food)
	return router
}

func main() {
	err := ConnectDb()
	if err != nil {
		panic(err)
	}
	log.Print("Connected to database")
	defer db.Close()
	router := setupRouter()
	router.Run()
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
