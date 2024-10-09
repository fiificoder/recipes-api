package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe
func init() {
	recipes = make([]Recipe, 0)

	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe


	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
			return
	}

	recipe.ID = xid.New().String()
recipe.PublishedAt = time.Now()
recipes = append(recipes, recipe)
c.JSON(http.StatusOK, recipe)
}


func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"Error": err.Error()})
			return
	}

	index :=  -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H {
			"Error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler (c *gin.Context) {
	id  := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H {
			"error": "recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H {
		"message": "recipe has been deleted",
	})
}

func SearchRecipeHandler (c *gin.Context) {
	tag := c.Query("tag")
	ListOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false

		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			ListOfRecipes = append(ListOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, ListOfRecipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)


	router.Run()
}