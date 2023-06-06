package main

import (
	"fmt"

	"os"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

type MovieRequest struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Genre string `json:"genre"`
}

var db *sql.DB

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}
	var err error
	db, err = sql.Open("sqlite3", "movies.db")

	if err != nil {
		fmt.Println("error connecting to database")
	} else {
		fmt.Print("Connection to database made successfully!\n\n")
	}

	defer db.Close()

	query := "CREATE TABLE IF NOT EXISTS movies (id TEXT PRIMARY KEY, title TEXT, year INTEGER, genre TEXT);"
	result, err := db.Query(query)

	if err != nil {
		fmt.Println("Something wrong with table creation query", err.Error())
	} else {
		fmt.Print("Table creation query successful!\n\n")
	}

	defer result.Close()

	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieById)
	router.POST("/movies", createMovie)

	router.Run(":" + PORT)
}

func getMovies(c *gin.Context) {
	var movies []Movie

	rows, err := db.Query("SELECT * FROM movies;")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"detail": "Error occurred while retrieving records!"})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var movie Movie

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"detail": "Error occurred while retrieving records!"})
			return
		}

		movies = append(movies, movie)
	}

	c.IndentedJSON(http.StatusOK, movies)
}

func getMovieById(c *gin.Context) {
	id := c.Param("id")

	rows, err := db.Query("SELECT * FROM movies;")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"detail": "Error occurred while retrieving records!"})
		return
	}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"detail": "Error occurred while retrieving records!"})
			return
		}

		if movie.ID == id {
			c.IndentedJSON(http.StatusOK, movie)
			return
		}

	}

	defer rows.Close()

	c.IndentedJSON(http.StatusNotFound, gin.H{"detail": "Requested resource not found!"})
}

func createMovie(c *gin.Context) {
	var newMovieRequest MovieRequest

	if err := c.BindJSON(&newMovieRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"detail": "Invalid format"})
		return
	}

	var newMovie = Movie{ID: uuid.NewString(), Title: newMovieRequest.Title, Year: newMovieRequest.Year, Genre: newMovieRequest.Genre}

	query := fmt.Sprintf("INSERT INTO movies VALUES('%s', '%s', %d, '%s');", newMovie.ID, newMovie.Title, newMovie.Year, newMovie.Genre)
	_, err := db.Exec(query)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newMovie)
}
