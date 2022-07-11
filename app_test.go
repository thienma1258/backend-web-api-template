package api

import (
	"context"
	"os"
	"testing"
)

var testApp *TestWrapper

type TestWrapper struct {
	*App
}

func NewTestApp() *TestWrapper {
	return &TestWrapper{}
}

// TestMain is the main entry for all tests
func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	var err error
	ctx := context.Background()
	testApp = NewTestApp()
	testApp.App, err = NewApp(ctx)
	if err != nil {
		panic(err)
	}
	testApp.App.redis, _ = NewRedisWrapper(&RedisConfig{
		DB:   0,
		Addr: "localhost:6379",
	})

	//testApp.App.redis.Client.FlushDB(ctx)
	return m.Run()
}
