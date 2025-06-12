package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Search struct {
	f IndexerFactory
	j JobsManager
}

func NewSearch(f IndexerFactory, j JobsManager) *Search {
	return &Search{f: f, j: j}
}

func (h *Search) asyncProcessing(c *gin.Context, scanner *bufio.Scanner, words []string, done chan struct{}) {
	id := c.GetString("request-id")

	// Set active job
	h.j.CreateJob(id)

	// Make as background operation
	go func() {
		defer close(done) // Signal completion

		// Start counting
		start := time.Now()

		idx := h.f.CreateNewIndex(scanner)
		idx.Search(words)

		// Finish counting
		duration := time.Since(start)
		errs := h.j.SaveJobResult(c, id, duration.Milliseconds(), idx)
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Printf("error saving results: %s\n", err.Error())
			}
		}
		idx.Close()
		fmt.Printf("Duration searching words: %d ms. Your PC is shit, check your mother \n", duration.Milliseconds())
	}()

	// Return jobId operation
	c.JSON(http.StatusOK, gin.H{"jobId": id})
}

func (h *Search) ProcessInput(c *gin.Context) {
	// Parse body
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}
	// Get words
	rawWords := form.Value["words"]
	var words []string
	if err := json.Unmarshal([]byte(rawWords[0]), &words); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid words array"})
		return
	}

	// Parse attached file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open file stream
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}

	// Create a scanner to read line by line
	scanner := bufio.NewScanner(file)

	// Create completion channel
	done := make(chan struct{})

	// Start async search
	h.asyncProcessing(c, scanner, words, done)

	// Close file after processing is complete
	go func() {
		<-done // Wait for processing to complete
		file.Close()
	}()
}
