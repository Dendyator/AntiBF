package api

import (
	"context"
	"testing"

	pb "github.com/Dendyator/AntiBF/api/proto/pb" //nolint
	"github.com/Dendyator/AntiBF/internal/config" //nolint
	"github.com/Dendyator/AntiBF/internal/core"   //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint

	"github.com/stretchr/testify/assert"
)

func TestCheckAuthorization(t *testing.T) {
	mockLogger := logger.New("info")
	core.InitLogger(mockLogger)
	InitLogger(mockLogger)

	core.InitRateLimiter(config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	})

	originalWhitelistFunc := core.WhitelistFunc
	core.WhitelistFunc = func(ip string) bool { return false }
	defer func() { core.WhitelistFunc = originalWhitelistFunc }()

	originalBlacklistFunc := core.BlacklistFunc
	core.BlacklistFunc = func(ip string) bool { return false }
	defer func() { core.BlacklistFunc = originalBlacklistFunc }()

	srv := NewServer(mockLogger)

	req := &pb.AuthRequest{
		Login:    "testuser",
		Password: "testpass",
		Ip:       "192.168.1.30/25",
	}

	resp, err := srv.CheckAuthorization(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
}

func TestResetBucket(t *testing.T) {
	core.InitLogger(logger.New("info"))

	core.InitRateLimiter(config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	})

	originalResetBucketFunc := core.ResetBucketFunc
	core.ResetBucketFunc = func(login, ip string) bool { return true }
	defer func() { core.ResetBucketFunc = originalResetBucketFunc }()

	log := logger.New("info")
	srv := NewServer(log)

	req := &pb.ResetRequest{
		Login: "testuser",
		Ip:    "192.168.1.30/25",
	}

	resp, err := srv.ResetBucket(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAddToBlacklist(t *testing.T) {
	originalManageListFunc := core.ManageListFunc
	core.ManageListFunc = func(subnet, listType string, add bool) bool {
		return listType == core.Blacklist && add
	}
	defer func() { core.ManageListFunc = originalManageListFunc }()

	log := logger.New("info")
	srv := NewServer(log)

	req := &pb.ListRequest{
		Subnet: "192.168.1.0/24",
	}

	resp, err := srv.AddToBlacklist(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}
