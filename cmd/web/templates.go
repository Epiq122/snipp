package main

import "snippet.robertgleason.ca/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
