package core

import (
	"net"
	"sync"
	"time"

	"github.com/Dendyator/AntiBF/internal/config" //nolint
	"github.com/Dendyator/AntiBF/internal/db"     //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint
)

const (
	Whitelist = "whitelist"
	Blacklist = "blacklist"
)

var (
	limiter   sync.Map
	appLogger *logger.Logger
	cfg       config.RateLimiterConfig
)

type Bucket struct {
	attempts    int
	lastAttempt time.Time
}

func InitLogger(log *logger.Logger) {
	appLogger = log
}

func InitRateLimiter(config config.RateLimiterConfig) {
	cfg = config
	appLogger.Infof("RateLimiter initialized with LoginLimit: %d, PasswordLimit: %d, IPLimit: %d",
		cfg.LoginLimit, cfg.PasswordLimit, cfg.IPLimit)
}

var WhitelistFunc = func(ip string) bool {
	appLogger.Debug("Entering WhitelistFunc")
	if !isValidCIDR(ip) {
		appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}

	inWhitelist, err := db.CheckInRedis(Whitelist, ip)
	if err != nil {
		appLogger.Warnf("Error checking whitelist for %s: %v", ip, err)
		return false
	}
	appLogger.Debugf("IP %s whitelisted: %v", ip, inWhitelist)
	return inWhitelist
}

var BlacklistFunc = func(ip string) bool {
	appLogger.Debug("Entering BlacklistFunc")
	if !isValidCIDR(ip) {
		appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}

	inBlacklist, err := db.CheckInRedis(Blacklist, ip)
	if err != nil {
		appLogger.Warnf("Error checking blacklist for %s: %v", ip, err)
		return false
	}
	appLogger.Debugf("IP %s blacklisted: %v", ip, inBlacklist)
	return inBlacklist
}

var ResetBucketFunc = func(login, ip string) bool {
	if !isValidCIDR(ip) {
		appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}
	limiter.Delete(login)
	limiter.Delete(ip)
	return true
}

var ManageListFunc = func(subnet string, listType string, add bool) bool {
	if !isValidCIDR(subnet) {
		appLogger.Warnf("Invalid subnet format: %s", subnet)
		return false
	}
	return db.UpdateRedis(listType, subnet, add)
}

func CheckAuthorization(login, password, ip string) bool {
	if !isValidCIDR(ip) {
		appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}

	if WhitelistFunc(ip) {
		appLogger.Infof("IP %s is whitelisted", ip)
		return true
	}
	if BlacklistFunc(ip) {
		appLogger.Infof("IP %s is blacklisted", ip)
		return false
	}
	if checkLogin(login) && checkIP(ip) && checkPassword(password) {
		appLogger.Infof("Authorization allowed for IP: %s, login: %s", ip, login)
		return true
	}
	appLogger.Warnf("Authorization attempt blocked for IP: %s, login: %s", ip, login)
	return false
}

func ResetBucket(login, ip string) bool {
	return ResetBucketFunc(login, ip)
}

func ManageList(subnet string, listType string, add bool) bool {
	return ManageListFunc(subnet, listType, add)
}

func isValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

func checkLogin(login string) bool {
	allowed := performRateLimiting(login, cfg.LoginLimit)
	appLogger.Debugf("Login attempts for %s allowed: %v", login, allowed)
	return allowed
}

func checkPassword(password string) bool {
	allowed := performRateLimiting(password, cfg.PasswordLimit)
	appLogger.Debugf("Password attempts for %s allowed: %v", password, allowed)
	return allowed
}

func checkIP(ip string) bool {
	allowed := performRateLimiting(ip, cfg.IPLimit)
	appLogger.Debugf("IP attempts for %s allowed: %v", ip, allowed)
	return allowed
}

func performRateLimiting(key string, limit int) bool {
	value, loaded := limiter.LoadOrStore(key, &Bucket{0, time.Now()})
	bucket := value.(*Bucket)

	if !loaded {
		appLogger.Debugf("New rate limit bucket created for key: %s", key)
	}

	if time.Since(bucket.lastAttempt) > time.Minute {
		appLogger.Debugf("Resetting attempts for key: %s due to timeout", key)
		bucket.attempts = 0
	}

	bucket.lastAttempt = time.Now()
	bucket.attempts++

	allowed := bucket.attempts <= limit
	if allowed {
		appLogger.Debugf("Attempt for key: %s, attempts: %d, within limit: %d", key, bucket.attempts, limit)
	} else {
		appLogger.Warnf("Attempt for key: %s, attempts: %d, exceeds limit: %d", key, bucket.attempts, limit)
	}

	return allowed
}
