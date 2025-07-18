package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Rate limiter map to store limiters for each IP
var rateLimiters = make(map[string]*rate.Limiter)

// RateLimitMiddleware creates a rate limiter for each IP
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Get or create rate limiter for this IP
		limiter, exists := rateLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(requestsPerMinute), requestsPerMinute)
			rateLimiters[ip] = limiter
		}

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, APIResponse{
				Success: false,
				Message: "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Strict rate limiter for authentication endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(5) // 5 requests per minute for auth
}

// General rate limiter for API endpoints
func GeneralRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(100) // 100 requests per minute for general API
}

// Security headers middleware
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict transport security (HTTPS only)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy - Allow CDN resources
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https:; font-src 'self' data: https:; connect-src 'self'")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// Input sanitization middleware
func InputSanitizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request body if exists
		if c.Request.Body != nil {
			// Add custom validation here if needed
			// For now, we rely on Gin's binding validation
		}

		c.Next()
	}
}

// Clean up old rate limiter entries periodically
func CleanupRateLimiters() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Clear all rate limiters to prevent memory leaks
				// In production, you might want to be more selective
				rateLimiters = make(map[string]*rate.Limiter)
			}
		}
	}()
}
