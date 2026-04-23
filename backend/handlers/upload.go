package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ledongthuc/pdf"
)

var NoteText string

func UploadNote(c *gin.Context) {
	file, err := c.FormFile("note")

	// handle no file
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// handle incorrect file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".md" && ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .md and .pdf files supported"})
		return
	}

	if err := os.MkdirAll("../uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create upload folder"})
		return
	}

	dst := filepath.Join("../uploads", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Extract text
	var text string
	if ext == ".md" {
		text, err = extractMarkdown(dst)
	} else if ext == ".pdf" {
		text, err = extractPDF(dst)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to extract text: %v", err)})
		return
	}

	NoteText = text
	//fmt.Println("successfully uploaded")
	c.JSON(http.StatusOK, gin.H{
		"message":  "note uploaded successfully",
		"filename": file.Filename,
		"chars":    len(NoteText),
	})
}

// Extract text from .md file
func extractMarkdown(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Extract text from PDF
func extractPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("could not open PDF: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	totalPages := r.NumPage()

	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("Failed to read page %d: %w", i, err)
		}
		buf.WriteString(text)
	}

	return buf.String(), nil
}
