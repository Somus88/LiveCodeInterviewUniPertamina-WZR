package main

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Major      string    `json:"major"`
	Year       int       `json:"year"`
	Status     string    `json:"status" binding:"required, oneof=active inactive"`
	Created_at time.Time `json:"created_at"`
}

var allStudents = []Student{
	{Id: "1", Name: "Wildan1", Email: "wzacky06@gmail.com", Major: "IT", Year: 2022, Status: "active", Created_at: time.Now()},
	{Id: "2", Name: "Zacky", Email: "wildanzr@gmail.com", Major: "IF", Year: 2023, Status: "inactive", Created_at: time.Now()},
	{Id: "3", Name: "Ramandhito", Email: "wzacky06@gmail.com", Major: "SE", Year: 2024, Status: "active", Created_at: time.Now()},
}

func main() {
	router := gin.Default()
	router.GET("/api/students", getAllMahasiswa)
	router.GET("/api/students?search={keyword}", searchMahasiswa)
	router.GET("/api/students?sort={field}&order={order}", sortMahasiswa)

	router.Run("localhost:8080")
}

func getAllMahasiswa(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, allStudents)
}

func searchMahasiswa(c *gin.Context) {
	keyword := c.Query("keyword")
	major := c.Query("major")
	status := c.Query("status")
	keep := true

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Success: False": "Message: Search keyword is empty"})
	}

	var results []Student

	for _, s := range allStudents {
		if major != "" && strings.ToLower(s.Major) != major {
			keep = false
		}

		if status != "" && strings.ToLower(s.Status) != status {
			keep = false
		}

		nameMatch := strings.Contains(strings.ToLower(s.Name), keyword)
		emailMatch := strings.Contains(strings.ToLower(s.Email), keyword)
		if nameMatch && emailMatch {
			keep = false
		}

		if keep {
			results = append(results, s)
		}
	}
	c.IndentedJSON(http.StatusOK, results)
}

func sortMahasiswa(c *gin.Context) {
	sortBy := strings.ToLower(c.DefaultQuery("sortBy", "created_at"))
	order := strings.ToLower(c.DefaultQuery("order", "asc"))

	allowedFields := map[string]bool{
		"name":       true,
		"created_at": true,
		"status":     true,
	}

	if !allowedFields[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"Success: False": "Message: Invalid sort field"})
	}

	results := make([]Student, len(allStudents))
	copy(results, allStudents)

	sort.Slice(results, func(i, j int) bool {
		var isLessThan bool

		switch sortBy {
		case "name":
			isLessThan = results[i].Name < results[j].Name
		case "status":
			isLessThan = results[i].Status < results[j].Status
		case "created_at":
			isLessThan = results[i].Created_at.Before(results[j].Created_at)
		}

		if order == "desc" {
			return !isLessThan
		}
		return isLessThan
	})
	c.IndentedJSON(http.StatusOK, results)
}
