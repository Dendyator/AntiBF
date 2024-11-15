package core

import (
	"github.com/Dendyator/AntiBF/internal/config"
	"github.com/Dendyator/AntiBF/internal/db"
	"github.com/Dendyator/AntiBF/internal/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckAuthorization(t *testing.T) {
	InitLogger(logger.New("info"))
	InitRateLimiter(config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	})

	WhitelistFunc = func(ip string) bool {
		return ip == "192.168.1.10/25"
	}
	BlacklistFunc = func(ip string) bool {
		return ip == "192.168.1.20/25"
	}

	ok := CheckAuthorization("testuser", "testpass", "192.168.1.10/25")
	if !ok {
		t.Errorf("Expected true for whitelisted IP, got false")
	}

	ok = CheckAuthorization("testuser", "testpass", "192.168.1.20/25")
	if ok {
		t.Errorf("Expected false for blacklisted IP, got true")
	}

	for i := 0; i < 10; i++ {
		if !CheckAuthorization("testuser", "testpass", "192.168.1.30/25") {
			t.Errorf("Expected true for valid login attempt, got false at iteration %d", i)
		}
	}

	if CheckAuthorization("testuser", "testpass", "192.168.1.30/25") {
		t.Errorf("Expected false after exceeding rate limit, got true")
	}
}

func TestResetBucket(t *testing.T) {
	InitLogger(logger.New("info"))
	performRateLimiting("testuser", 10)
	performRateLimiting("192.168.1.30/25", 1000)

	ResetBucket("testuser", "192.168.1.30/25")

	if _, exists := limiter.Load("testuser"); exists {
		t.Errorf("Expected bucket to be reset for login-key, but entry still exists")
	}

	if _, exists := limiter.Load("192.168.1.30/25"); exists {
		t.Errorf("Expected bucket to be reset for IP, but entry still exists")
	}
}

func TestAuthorizationLimits(t *testing.T) {
	InitLogger(logger.New("debug"))

	db.InitRedis("localhost:6379", logger.New("debug"))

	InitRateLimiter(config.RateLimiterConfig{
		LoginLimit:    5,
		PasswordLimit: 5,
		IPLimit:       5,
	})

	ResetBucket("testuser", "192.168.1.1/25")
	ResetBucket("testpass", "192.168.1.1/25")

	for i := 0; i < 5; i++ {
		ok := CheckAuthorization("testuser", "testpass", "192.168.1.1/25")
		assert.True(t, ok, "Authorization should be allowed")
	}

	ok := CheckAuthorization("testuser", "testpass", "192.168.1.1/25")
	assert.False(t, ok, "Authorization should be blocked after exceeding limit")
}
