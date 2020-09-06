package shortUrls

import (
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestShortUrlService(t *testing.T) {
	service := shortUrlServiceImpl{
		user: &models.SignedInUser{UserId: 1},
	}

	mockUid := "testuid"
	mockNotFoundUid := "testnotfounduid"
	mockPath := "mock/path?test=true"
	mockShortUrl := models.ShortUrl{
		Uid:       mockUid,
		Path:      mockPath,
		CreatedBy: service.user.UserId,
		CreatedAt: time.Now(),
	}
	mockNotFoundShortUrl := models.ShortUrl{
		Uid:       mockNotFoundUid,
		Path:      "",
		CreatedBy: service.user.UserId,
		CreatedAt: time.Now(),
	}

	bus.AddHandler("test", func(query *models.CreateShortUrlCommand) error {
		query.Result = &mockShortUrl
		return nil
	})

	bus.AddHandler("test", func(query *models.GetFullUrlQuery) error {
		result := &mockShortUrl
		if query.Uid == mockNotFoundUid {
			result = &mockNotFoundShortUrl
		}
		query.Result = result
		return nil
	})

	t.Run("User can create and read short URLs", func(t *testing.T) {
		uid, err := service.CreateShortUrl(&dtos.CreateShortUrlForm{
			Path: "test/short/url",
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, uid)
		assert.Equal(t, uid, mockUid)
		path, err := service.GetFullUrlByUID(uid)
		assert.Nil(t, err)
		assert.NotEmpty(t, path)
		assert.Equal(t, path, mockPath)
	})

	t.Run("User cannot look up nonexistent short urls", func(t *testing.T) {
		service := shortUrlServiceImpl{
			user: &models.SignedInUser{UserId: 1},
		}

		path, err := service.GetFullUrlByUID(mockNotFoundUid)
		assert.NotNil(t, err)
		assert.Empty(t, path)
		assert.Equal(t, err, models.ErrShortUrlNotFound)
	})
}
