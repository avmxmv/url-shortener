package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemStorage_CreateLink(t *testing.T) {
	store := NewInMemStorage()

	originalURL := "https://example.com"

	shortURL, err := store.CreateLink(originalURL)
	assert.NoError(t, err, "CreateLink should not return an error")
	assert.NotEmpty(t, shortURL, "Short URL should not be empty")

	shortURL2, err := store.CreateLink(originalURL)
	assert.NoError(t, err, "CreateLink should not return an error for the same URL")
	assert.Equal(t, shortURL, shortURL2, "Short URL should be the same for the same original URL")
}

func TestInMemStorage_GetLink(t *testing.T) {
	store := NewInMemStorage()

	originalURL := "https://example.com"

	shortURL, err := store.CreateLink(originalURL)
	assert.NoError(t, err, "CreateLink should not return an error")

	retrievedURL, err := store.GetLink(shortURL)
	assert.NoError(t, err, "GetLink should not return an error")
	assert.Equal(t, originalURL, retrievedURL, "Retrieved URL should match the original URL")

	_, err = store.GetLink("nonexistent")
	assert.ErrorIs(t, err, ErrNotFound, "GetLink should return ErrNotFound for nonexistent short URL")
}

func TestInMemStorage_Concurrency(t *testing.T) {
	store := NewInMemStorage()

	numGoroutines := 100

	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			originalURL := "https://example.com/" + string(rune(i))
			_, err := store.CreateLink(originalURL)
			assert.NoError(t, err, "CreateLink should not return an error")
			done <- true
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	links := make(map[string]bool)
	for i := 0; i < numGoroutines; i++ {
		originalURL := "https://example.com/" + string(rune(i))
		shortURL, err := store.CreateLink(originalURL)
		assert.NoError(t, err, "CreateLink should not return an error")
		assert.False(t, links[shortURL], "Short URL should be unique")
		links[shortURL] = true
	}
}
