package cache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisService(t *testing.T) {
	assertWithTest := assert.New(t)
	// Create a payload
	_, err := createTestClient()
	assertWithTest.Nil(err)

}

var newBooks = []*books.Book{
	{
		ID:        1,
		ISBN:      "978-1234567890",
		Title:     "The Great Gatsby",
		Author:    "F. Scott Fitzgerald",
		Publisher: "Scribner",
		Published: utils.CustomDate{Time: time.Date(1925, 4, 10, 0, 0, 0, 0, time.UTC)},
		Genre:     "Fiction",
		Language:  "English",
		Pages:     180,
		Available: books.Available,
		UpdatedAt: utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
		CreatedAt: utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
		DeletedAt: utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
	},
	// Add more books as needed...
	{
		ID:        2,
		ISBN:      "978-0451524935",
		Title:     "1984",
		Author:    "George Orwell",
		Publisher: "Signet Classic",
		Published: utils.CustomDate{Time: time.Date(1949, 6, 8, 0, 0, 0, 0, time.UTC)},
		Genre:     "Dystopian",
		Language:  "English",
		Pages:     328,
		Available: books.Available,
		UpdatedAt: utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
		CreatedAt: utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
		DeletedAt: utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
	},
}
var testClient, _ = createTestClient()

func TestStoreValue(t *testing.T) {
	assertWithTest := assert.New(t)
	// Create a payload
	testService := NewRedisCache(testClient)
	// Create a slice of cache payloads
	data := []*CachePayload{}
	for _, book := range newBooks {
		jsonData, err := json.Marshal(book)
		assertWithTest.Nil(err, "Data has been serialized")
		asset := &CachePayload{
			Key:        book.Author,
			Value:      jsonData,
			Expiration: 1 * time.Minute,
		}
		data = append(data, asset)
	}

	err := testService.Store(context.Background(), data)
	assertWithTest.Nil(err, "Store the value with no error")

	// Do an existence check
	exists, err := testService.ExistenceCheck(context.Background(), data[0].Key)
	assertWithTest.Nil(err)
	assertWithTest.True(exists, "The policy has just been added so must return true")
	// Do a false check
	exists, err = testService.ExistenceCheck(context.Background(), "fiuvosnviosjn")
	assertWithTest.Nil(err)
	assertWithTest.False(exists, "Ficticious key so existence is false")
}

func TestRetrieveValue(t *testing.T) {
	assertWithTest := assert.New(t)
	// Create a payload
	testService := NewRedisCache(testClient)
	// Create a slice of cache payloads
	data := []*CachePayload{}
	for _, book := range newBooks {
		jsonData, err := json.Marshal(book)
		assertWithTest.Nil(err, "Data has been serialized")
		asset := &CachePayload{
			Key:        book.Author,
			Value:      jsonData,
			Expiration: 1 * time.Minute,
		}
		data = append(data, asset)
	}

	// Store data
	err := testService.Store(context.Background(), data)
	assertWithTest.Nil(err, "Store the value with no error")
	// Retrieve the stored data
	retrievedData, err := testService.Retrieve(context.Background(),
		[]string{newBooks[0].Author})
	assertWithTest.Nil(err, "Data should be retrieved successfully")
	assertWithTest.Equal(data[0].Value, retrievedData[0].Value)

}

func TestClearCacheByKeys(t *testing.T) {
	assertWithTest := assert.New(t)
	ctx := context.Background()
	// Create a payload
	testService := NewRedisCache(testClient)
	// Create a slice of cache payloads
	data := []*CachePayload{}
	keys := []string{}
	for _, book := range newBooks {
		jsonData, err := json.Marshal(book)
		assertWithTest.Nil(err, "Data has been serialized")
		asset := &CachePayload{
			Key:        book.Author,
			Value:      jsonData,
			Expiration: 1 * time.Minute,
		}
		data = append(data, asset)
		// Create a slice of keys to clear
		keys = append(keys, asset.Key)
	}

	// Store data
	err := testService.Store(ctx, data)
	assertWithTest.Nil(err, "Store the value with no error")
	// Delete data
	err = testService.ClearCacheByKeys(ctx, keys)
	assertWithTest.Nil(err, "Delete values from cache successfully")
	// Retrieve the stored data
	values, err := testService.Retrieve(context.Background(),
		[]string{newBooks[0].Author})
	assertWithTest.Nil(err, "Data should be retrieved successfully")
	assertWithTest.Nil(values, "Nothing in cache to return")

}

func createTestClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6389",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
