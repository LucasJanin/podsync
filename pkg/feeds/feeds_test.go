//go:generate mockgen -source=feeds.go -destination=feeds_mock_test.go -package=feeds

package feeds

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	itunes "github.com/mxpv/podcast"

	"github.com/mxpv/podsync/pkg/api"
	"github.com/mxpv/podsync/pkg/model"
)

var feed = &model.Feed{
	HashID:   "123",
	ItemID:   "xyz",
	Provider: api.ProviderVimeo,
	LinkType: api.LinkTypeChannel,
	PageSize: 50,
	Quality:  api.QualityHigh,
	Format:   api.FormatVideo,
}

func TestService_CreateFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockstorage(ctrl)
	db.EXPECT().SaveFeed(gomock.Any()).Times(1).Return(nil)

	gen, _ := NewIDGen()

	s := Service{
		generator: gen,
		db:        db,
		builders:  map[api.Provider]Builder{api.ProviderYoutube: nil},
	}

	req := &api.CreateFeedRequest{
		URL:      "youtube.com/channel/123",
		PageSize: 50,
		Quality:  api.QualityHigh,
		Format:   api.FormatVideo,
	}

	hashID, err := s.CreateFeed(req, &api.Identity{})
	require.NoError(t, err)
	require.NotEmpty(t, hashID)
}

func TestService_makeFeed(t *testing.T) {
	req := &api.CreateFeedRequest{
		URL:      "youtube.com/channel/123",
		PageSize: 1000,
		Quality:  api.QualityLow,
		Format:   api.FormatAudio,
	}

	gen, _ := NewIDGen()

	s := Service{
		generator: gen,
	}

	feed, err := s.makeFeed(req, &api.Identity{})
	require.NoError(t, err)
	require.Equal(t, 50, feed.PageSize)
	require.Equal(t, api.QualityHigh, feed.Quality)
	require.Equal(t, api.FormatVideo, feed.Format)

	feed, err = s.makeFeed(req, &api.Identity{FeatureLevel: api.ExtendedFeatures})
	require.NoError(t, err)
	require.Equal(t, 150, feed.PageSize)
	require.Equal(t, api.QualityLow, feed.Quality)
	require.Equal(t, api.FormatAudio, feed.Format)

	feed, err = s.makeFeed(req, &api.Identity{FeatureLevel: api.ExtendedPagination})
	require.NoError(t, err)
	require.Equal(t, 600, feed.PageSize)
	require.Equal(t, api.QualityLow, feed.Quality)
	require.Equal(t, api.FormatAudio, feed.Format)
}

func TestService_QueryFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockstorage(ctrl)
	db.EXPECT().GetFeed("123").Times(1).Return(nil, nil)

	s := Service{db: db}
	_, err := s.QueryFeed("123")
	require.NoError(t, err)
}

func TestService_GetFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	item := CacheItem{
		UpdatedAt: time.Now().UTC(),
		Feed:      []byte("test"),
	}

	cache := NewMockcacheService(ctrl)
	cache.EXPECT().GetItem("123", gomock.Any()).DoAndReturn(func(_ string, ret *CacheItem) error {
		*ret = item
		return nil
	})

	s := Service{cache: cache}

	data, err := s.BuildFeed("123")
	assert.NoError(t, err)
	assert.Equal(t, item.Feed, data)
}

func TestService_VerifyCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := NewMockcacheService(ctrl)
	cache.EXPECT().GetItem("123", gomock.Any()).DoAndReturn(func(_ string, ret *CacheItem) error {
		ret.Feed = []byte("test")
		ret.UpdatedAt = time.Now().UTC().Add(-20 * time.Minute)
		ret.ItemCount = 30
		return nil
	})

	cache.EXPECT().SaveItem("123", gomock.Any(), 15*24*time.Hour).Times(1).Return(nil)

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(feed.HashID).Times(1).Return(feed, nil)

	builder := NewMockBuilder(ctrl)
	builder.EXPECT().GetVideoCount(feed).Return(uint64(30), nil)

	s := Service{db: stor, cache: cache, builders: map[api.Provider]Builder{
		api.ProviderVimeo: builder,
	}}

	data, err := s.BuildFeed(feed.HashID)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), data)
}

func TestService_BuildFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(feed.HashID).Times(1).Return(feed, nil)

	cache := NewMockcacheService(ctrl)
	cache.EXPECT().GetItem(feed.HashID, gomock.Any()).Return(errors.New("not found"))
	cache.EXPECT().SaveItem(feed.HashID, gomock.Any(), gomock.Any()).Return(nil)

	podcast := itunes.New("", "", "", nil, nil)

	builder := NewMockBuilder(ctrl)
	builder.EXPECT().Build(feed).Return(&podcast, nil)
	builder.EXPECT().GetVideoCount(feed).Return(uint64(25), nil)

	s := Service{db: stor, cache: cache, builders: map[api.Provider]Builder{
		api.ProviderVimeo: builder,
	}}

	_, err := s.BuildFeed(feed.HashID)
	require.NoError(t, err)
}

func TestService_RebuildCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(feed.HashID).Times(1).Return(feed, nil)

	cache := NewMockcacheService(ctrl)
	cache.EXPECT().GetItem("123", gomock.Any()).DoAndReturn(func(_ string, ret *CacheItem) error {
		ret.Feed = []byte("test")
		ret.UpdatedAt = time.Now().UTC().Add(-20 * time.Minute)
		ret.ItemCount = 30
		return nil
	})

	cache.EXPECT().SaveItem(feed.HashID, gomock.Any(), gomock.Any()).Return(nil)

	podcast := itunes.New("", "", "", nil, nil)

	builder := NewMockBuilder(ctrl)
	builder.EXPECT().Build(feed).Return(&podcast, nil)
	builder.EXPECT().GetVideoCount(feed).Return(uint64(25), nil)

	s := Service{db: stor, cache: cache, builders: map[api.Provider]Builder{
		api.ProviderVimeo: builder,
	}}

	_, err := s.BuildFeed(feed.HashID)
	require.NoError(t, err)
}

func TestService_WrongID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := NewMockcacheService(ctrl)
	cache.EXPECT().GetItem(gomock.Any(), gomock.Any()).Return(errors.New("not found"))

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetFeed(gomock.Any()).Times(1).Return(nil, errors.New("not found"))

	s := Service{db: stor, cache: cache}

	_, err := s.BuildFeed("invalid_feed_id")
	require.Error(t, err)
}

func TestService_GetMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stor := NewMockstorage(ctrl)
	stor.EXPECT().GetMetadata(feed.HashID).Times(1).Return(feed, nil)

	s := Service{db: stor}

	m, err := s.GetMetadata(feed.HashID)
	require.NoError(t, err)
	require.EqualValues(t, 0, m.Downloads)
}
