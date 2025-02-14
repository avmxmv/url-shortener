package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"url-shortener/api"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

type MockStorage struct {
	storage.Storage
	CreateLinkFunc func(string) (string, error)
	GetLinkFunc    func(string) (string, error)
}

func (m *MockStorage) CreateLink(originalURL string) (string, error) {
	return m.CreateLinkFunc(originalURL)
}

func (m *MockStorage) GetLink(shortURL string) (string, error) {
	return m.GetLinkFunc(shortURL)
}

func TestService_CreateLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful creation", func(t *testing.T) {
		mockStorage := &MockStorage{
			CreateLinkFunc: func(originalURL string) (string, error) {
				return "abc123", nil
			},
		}

		srv := service.NewService(mockStorage)
		resp, err := srv.CreateLink(context.Background(), &api.CreateLinkRequest{
			OriginalUrl: "https://example.com",
		})

		assert.NoError(t, err)
		assert.Equal(t, "abc123", resp.ShortUrl)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := &MockStorage{
			CreateLinkFunc: func(originalURL string) (string, error) {
				return "", assert.AnError
			},
		}

		srv := service.NewService(mockStorage)
		resp, err := srv.CreateLink(context.Background(), &api.CreateLinkRequest{
			OriginalUrl: "https://example.com",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestService_GetLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("successful get", func(t *testing.T) {
		mockStorage := &MockStorage{
			GetLinkFunc: func(shortURL string) (string, error) {
				return "https://example.com", nil
			},
		}

		srv := service.NewService(mockStorage)
		resp, err := srv.GetLink(context.Background(), &api.GetLinkRequest{
			ShortUrl: "abc123",
		})

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", resp.OriginalUrl)
	})

	t.Run("not found", func(t *testing.T) {
		mockStorage := &MockStorage{
			GetLinkFunc: func(shortURL string) (string, error) {
				return "", storage.ErrNotFound
			},
		}

		srv := service.NewService(mockStorage)
		resp, err := srv.GetLink(context.Background(), &api.GetLinkRequest{
			ShortUrl: "invalid",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := &MockStorage{
			GetLinkFunc: func(shortURL string) (string, error) {
				return "", assert.AnError
			},
		}

		srv := service.NewService(mockStorage)
		resp, err := srv.GetLink(context.Background(), &api.GetLinkRequest{
			ShortUrl: "abc123",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
