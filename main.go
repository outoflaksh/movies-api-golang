package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
)

type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

var movie_db = []Movie{
	{ID: 1, Title: "The Shawshank Redemption", Year: 1994, Genre: "Drama"},
	{ID: 2, Title: "The Godfather", Year: 1972, Genre: "Crime"},
	{ID: 3, Title: "Pulp Fiction", Year: 1994, Genre: "Crime"},
	{ID: 4, Title: "The Dark Knight", Year: 2008, Genre: "Action"},
}

func main() {
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieById)
	router.POST("/movies", createMovie)

	router.Run(":8000")
}

func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, movie_db)
}

func getMovieById(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	for _, val := range movie_db {
		if val.ID == id {
			c.IndentedJSON(http.StatusOK, val)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"detail": "Requested resource not found!"})

}

func createMovie(c *gin.Context) {
	var newMovie Movie

	if err := c.BindJSON(&newMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"detail": "Invalid format"})
		return
	}

	movie_db = append(movie_db, newMovie)

	c.IndentedJSON(http.StatusCreated, newMovie)
}
