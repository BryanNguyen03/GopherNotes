package handlers

import (
	"net/http"

	"GopherNotes/ai"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func Chat(c *gin.Context) {

	// parse incmoing JSON body
	var req ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// handle empty message
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message cannot be empty"})
		return
	}

	// build prompt and send to llama
	prompt := buildPrompt(req.Message)
	reply, err := ai.AskLlama(prompt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "message cannot be empty"})
		return
	}

	//fmt.Println("successfully sent prompt")
	c.JSON(http.StatusOK, gin.H{"reply": reply})

}

func buildPrompt(userMessage string) string {
	noteSection := "No note has been loaded"
	if NoteText != "" {
		noteSection = NoteText
	}

	return "[INST] <<SYS>>\n" +
		"You are a helpful learning assistant tool. " +
		"Answer using ONLY the note provided. " +
		"Be concise and direct. " +
		"Do NOT explain what you are doing. " +
		"Do NOT narrate your process. " +
		"Just answer.\n" +
		"<</SYS>>\n\n" +
		"Note:\n" + noteSection + "\n\n" +
		"Question: " + userMessage + " [/INST]"

}
