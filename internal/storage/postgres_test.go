package storage

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewPostgresStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	storage, err := NewPostgresStorage("mock_dsn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.NotNil(t, storage)
	assert.NotNil(t, storage.db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := &PostgresStorage{db: db}

	originalURL := "https://example.com"
	shortURL := "abc123"

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO links (original_url, short_url) 
		VALUES ($1, $2) 
		ON CONFLICT (original_url) DO NOTHING
		RETURNING short_url`,
	)).WithArgs(originalURL, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"short_url"}).AddRow(shortURL))

	result, err := storage.CreateLink(originalURL)
	assert.NoError(t, err)
	assert.Equal(t, shortURL, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := &PostgresStorage{db: db}

	shortURL := "abc123"
	originalURL := "https://example.com"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT original_url FROM links WHERE short_url = $1",
	)).WithArgs(shortURL).WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow(originalURL))

	result, err := storage.GetLink(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLink_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := &PostgresStorage{db: db}

	shortURL := "abc123"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT original_url FROM links WHERE short_url = $1",
	)).WithArgs(shortURL).WillReturnError(sql.ErrNoRows)

	result, err := storage.GetLink(shortURL)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Empty(t, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
