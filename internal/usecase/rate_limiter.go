package usecase

import (
	"github.com/Dendyator/AntiBF/internal/repository"
	"net"
	"sync"
	"time"

	"github.com/Dendyator/AntiBF/pkg/config"
	"github.com/Dendyator/AntiBF/pkg/logger"
)

const (
	Whitelist = "whitelist"
	Blacklist = "blacklist"
)

type RateLimiter struct {
	repo      repository.Repository
	userRepo  repository.UserRepositoryInterface // Добавляем UserRepository
	cfg       config.RateLimiterConfig
	appLogger *logger.Logger
	limiter   sync.Map
}

// Создание нового RateLimiter
func NewRateLimiter(repo repository.Repository, userRepo repository.UserRepositoryInterface, cfg config.RateLimiterConfig, appLogger *logger.Logger) *RateLimiter {
	return &RateLimiter{
		repo:      repo,
		userRepo:  userRepo,
		cfg:       cfg,
		appLogger: appLogger,
	}
}

func (r *RateLimiter) IsWhitelisted(ip string) bool {
	if !isValidCIDR(ip) {
		r.appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}
	inWhitelist, err := r.repo.CheckInList(Whitelist, ip)
	if err != nil {
		r.appLogger.Warnf("Error checking whitelist for %s: %v", ip, err)
		return false
	}
	r.appLogger.Debugf("IP %s whitelisted: %v", ip, inWhitelist)
	return inWhitelist
}

func (r *RateLimiter) IsBlacklisted(ip string) bool {
	if !isValidCIDR(ip) {
		r.appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}
	inBlacklist, err := r.repo.CheckInList(Blacklist, ip)
	if err != nil {
		r.appLogger.Warnf("Error checking blacklist for %s: %v", ip, err)
		return false
	}
	r.appLogger.Debugf("IP %s blacklisted: %v", ip, inBlacklist)
	return inBlacklist
}

// Сброс bucket'ов для login и IP
func ResetBucketFunc(r *RateLimiter, login, ip string) bool {
	if !isValidCIDR(ip) {
		r.appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}
	r.limiter.Delete(login)
	r.limiter.Delete(ip)
	return true
}

// Управление списками (white/black)
func ManageListFunc(r *RateLimiter, subnet string, listType string, add bool) bool {
	if !isValidCIDR(subnet) {
		r.appLogger.Warnf("Invalid subnet format: %s", subnet)
		return false
	}
	return r.repo.UpdateList(listType, subnet, add)
}

// Основная функция проверки авторизации
func (r *RateLimiter) CheckAuthorization(login, password, ip string) bool {
	if !isValidCIDR(ip) {
		r.appLogger.Warnf("Invalid IP format: %s", ip)
		return false
	}

	if r.IsWhitelisted(ip) {
		r.appLogger.Infof("IP %s is whitelisted", ip)
		return true
	}

	if r.IsBlacklisted(ip) {
		r.appLogger.Infof("IP %s is blacklisted", ip)
		return false
	}

	if r.checkLogin(login) && r.checkPassword(password) && r.checkIP(ip) {
		r.appLogger.Infof("Authorization allowed for IP: %s, login: %s", ip, login)
		return true
	}

	r.appLogger.Warnf("Authorization attempt blocked for IP: %s, login: %s", ip, login)
	return false
}

// Сброс bucket'ов
func (r *RateLimiter) ResetBucket(login, ip string) bool {
	return ResetBucketFunc(r, login, ip)
}

// Управление списками
func (r *RateLimiter) ManageList(subnet string, listType string, add bool) bool {
	return ManageListFunc(r, subnet, listType, add)
}

// Валидация CIDR
func isValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// Проверка логина
func (r *RateLimiter) checkLogin(login string) bool {
	allowed := r.performRateLimiting(login, r.cfg.LoginLimit)
	r.appLogger.Debugf("Login attempts for %s allowed: %v", login, allowed)
	return allowed
}

// Проверка пароля
func (r *RateLimiter) checkPassword(password string) bool {
	allowed := r.performRateLimiting(password, r.cfg.PasswordLimit)
	r.appLogger.Debugf("Password attempts for %s allowed: %v", password, allowed)
	return allowed
}

// Проверка IP
func (r *RateLimiter) checkIP(ip string) bool {
	allowed := r.performRateLimiting(ip, r.cfg.IPLimit)
	r.appLogger.Debugf("IP attempts for %s allowed: %v", ip, allowed)
	return allowed
}

// Алгоритм rate limiting
func (r *RateLimiter) performRateLimiting(key string, limit int) bool {
	value, loaded := r.limiter.LoadOrStore(key, &Bucket{0, time.Now()})
	bucket := value.(*Bucket)

	if !loaded || time.Since(bucket.lastAttempt) > time.Minute {
		r.appLogger.Debugf("Resetting attempts for key: %s due to timeout", key)
		bucket.attempts = 0
	}

	bucket.lastAttempt = time.Now()
	bucket.attempts++
	allowed := bucket.attempts <= limit

	if allowed {
		r.appLogger.Debugf("Attempt for key: %s, attempts: %d, within limit: %d", key, bucket.attempts, limit)
	} else {
		r.appLogger.Warnf("Attempt for key: %s, attempts: %d, exceeds limit: %d", key, bucket.attempts, limit)
	}

	return allowed
}

// Bucket для хранения счетчиков попыток
type Bucket struct {
	attempts    int
	lastAttempt time.Time
}
