package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Load environment variables from .env file if exists
	loadEnvFile()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	var err error
	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "blog.db"
	}

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize database schema and data
	InitDB(db)

	// Seed database with sample data
	SeedDatabase(db)

	// Initialize Gin router
	r := gin.Default()

	// Start rate limiter cleanup
	CleanupRateLimiters()

	// Security headers middleware
	r.Use(SecurityHeadersMiddleware())

	// General rate limiting for all endpoints
	r.Use(GeneralRateLimitMiddleware())

	// Input sanitization middleware
	r.Use(InputSanitizationMiddleware())

	// CORS configuration - more restrictive for production
	allowedOrigins := []string{"http://localhost:8080", "http://127.0.0.1:8080"}
	if os.Getenv("PRODUCTION_DOMAIN") != "" {
		allowedOrigins = append(allowedOrigins, os.Getenv("PRODUCTION_DOMAIN"))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	setupRoutes(r)

	// Serve static files with path validation
	baseDir := getBaseDir()
	log.Printf("DEBUG: baseDir = %s", baseDir)

	// Test file existence
	if _, err := os.Stat(baseDir + "/index.html"); os.IsNotExist(err) {
		log.Printf("ERROR: index.html not found at %s", baseDir+"/index.html")
	} else {
		log.Printf("DEBUG: index.html found at %s", baseDir+"/index.html")
	}

	r.GET("/sea.jpeg", func(c *gin.Context) {
		c.Header("Content-Type", "image/jpeg")
		c.File(baseDir + "/sea.jpeg")
	})

	// Serve SEO files
	r.GET("/sitemap.xml", func(c *gin.Context) {
		c.Header("Content-Type", "application/xml")
		c.File(baseDir + "/sitemap.xml")
	})

	r.GET("/robots.txt", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain")
		c.File(baseDir + "/robots.txt")
	})

	// Serve Google verification file
	r.GET("/googlea40d1f1aabd0d48d.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.File(baseDir + "/googlea40d1f1aabd0d48d.html")
	})

	// Serve favicon
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Header("Content-Type", "image/x-icon")
		c.File(baseDir + "/favicon.ico")
	})

	// Serve HTML files with security headers
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/index.html")
	})

	// Serve HTML files
	r.GET("/admin", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/admin.html")
	})
	r.GET("/iletisim", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/iletisim.html")
	})
	r.GET("/hakkimda", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/hakkimda.html")
	})
	r.GET("/yazilar", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/yazilar.html")
	})
	r.GET("/projeler", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/projeler.html")
	})
	r.GET("/post1", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/post1.html")
	})
	r.GET("/yakinda", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/yakinda.html")
	})
	r.GET("/posts/:slug", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(baseDir + "/post-detay.html")
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start server
	log.Printf("Server starting on :%s...", port)
	log.Fatal(r.Run(":" + port))
}

func setupRoutes(r *gin.Engine) {
	// Public routes
	api := r.Group("/api")
	{
		// Authentication - with stricter rate limiting
		authGroup := api.Group("/auth").Use(AuthRateLimitMiddleware())
		{
			authGroup.POST("/login", login)
			authGroup.POST("/refresh", refreshToken)
		}

		// Public content
		api.GET("/posts", getPosts)
		api.GET("/posts/:slug", getPost)
		api.GET("/projects", getProjects)
		api.GET("/projects/:slug", getProject)
		api.GET("/settings", getSettings)
		api.POST("/messages", createMessage)
		api.GET("/comingsoon", getComingSoon)
		api.POST("/newsletter", subscribeNewsletter)

		// Increment view count
		api.POST("/posts/:id/view", incrementViewCount)
	}

	// Protected routes (require authentication)
	auth := api.Group("/").Use(AuthMiddleware())
	{
		// Dashboard
		auth.GET("/dashboard/stats", getDashboardStats)

		// Posts management
		auth.GET("/admin/posts", getAdminPosts)
		auth.POST("/admin/posts", createPost)
		auth.GET("/admin/posts/:id", getAdminPost)
		auth.PUT("/admin/posts/:id", updatePost)
		auth.DELETE("/admin/posts/:id", deletePost)

		// Projects management
		auth.GET("/admin/projects", getAdminProjects)
		auth.POST("/admin/projects", createProject)
		auth.GET("/admin/projects/:id", getAdminProject)
		auth.PUT("/admin/projects/:id", updateProject)
		auth.DELETE("/admin/projects/:id", deleteProject)

		// Messages management
		auth.GET("/admin/messages", getMessages)
		auth.GET("/admin/messages/:id", getMessage)
		auth.PUT("/admin/messages/:id", updateMessage)
		auth.DELETE("/admin/messages/:id", deleteMessage)

		// Coming Soon management
		auth.GET("/admin/comingsoon", getComingSoon)
		auth.POST("/admin/comingsoon", createComingSoon)
		auth.GET("/admin/comingsoon/:id", getComingSoonByID)
		auth.PUT("/admin/comingsoon/:id", updateComingSoon)
		auth.DELETE("/admin/comingsoon/:id", deleteComingSoon)

		// Newsletter management
		auth.GET("/admin/newsletter", getNewsletterSubscriptions)
		auth.DELETE("/admin/newsletter/:id", deleteNewsletterSubscription)

		// Settings management
		auth.PUT("/admin/settings", updateSettings)

		// User management
		auth.GET("/user/profile", getUserProfile)
		auth.PUT("/user/profile", updateUserProfile)

		// SEO management
		auth.GET("/sitemap/generate", generateSitemap)
	}
}

// Authentication handlers
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// Log the actual error for debugging but don't expose it
		log.Printf("Login error for user %s: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		log.Printf("Invalid password attempt for user: %s", req.Username)
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}

// Refresh token handler
func refreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	// Verify refresh token
	claims, err := VerifyTokenWithType(req.RefreshToken, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Invalid or expired refresh token",
		})
		return
	}

	// Get user from database
	var user User
	if err := db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Generate new access token
	newAccessToken, err := GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to generate new access token",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: RefreshTokenResponse{
			AccessToken: newAccessToken,
		},
	})
}

// Dashboard handlers
func getDashboardStats(c *gin.Context) {
	var stats DashboardStats

	db.Model(&Post{}).Count(&stats.TotalPosts)
	db.Model(&Post{}).Where("status = ?", "published").Count(&stats.PublishedPosts)
	db.Model(&Post{}).Where("status = ?", "draft").Count(&stats.DraftPosts)
	db.Model(&Project{}).Count(&stats.TotalProjects)
	db.Model(&Project{}).Where("status = ?", "active").Count(&stats.ActiveProjects)
	db.Model(&Project{}).Where("status = ?", "completed").Count(&stats.CompletedProjects)
	db.Model(&Message{}).Where("is_read = ?", false).Count(&stats.UnreadMessages)
	db.Model(&Post{}).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews)
	db.Model(&Newsletter{}).Where("is_active = ?", true).Count(&stats.NewsletterSubscribers)

	log.Printf("Dashboard stats: Newsletter subscribers = %d", stats.NewsletterSubscribers)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

// Public post handlers
func getPosts(c *gin.Context) {
	page, limit, offset := GetPaginationParams(c)
	search := GetSearchQuery(c)
	category := GetCategoryFilter(c)

	var posts []Post
	query := db.Model(&Post{}).Where("status = ?", "published").Order("created_at DESC")

	if search != "" {
		// Sanitize search input to prevent SQL injection
		search = strings.ReplaceAll(search, "%", "\\%")
		search = strings.ReplaceAll(search, "_", "\\_")
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if category != "" {
		// Sanitize category input
		category = strings.ReplaceAll(category, "%", "\\%")
		category = strings.ReplaceAll(category, "_", "\\_")
		query = query.Where("categories LIKE ?", "%"+category+"%")
	}

	var total int64
	countResult := query.Count(&total)
	if countResult.Error != nil {
		log.Printf("Error counting posts: %v", countResult.Error)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	findResult := query.Offset(offset).Limit(limit).Find(&posts)
	if findResult.Error != nil {
		log.Printf("Error finding posts: %v", findResult.Error)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	log.Printf("Found %d posts, total: %d", len(posts), total)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"posts":       posts,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func getPost(c *gin.Context) {
	slug := c.Param("slug")
	var post Post

	if err := db.Where("slug = ? AND status = ?", slug, "published").First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    post,
	})
}

func incrementViewCount(c *gin.Context) {
	id := c.Param("id")
	clientIP := c.ClientIP()

	log.Printf("View count request - Post ID: %s, Client IP: %s", id, clientIP)

	// Convert id to uint
	var postID uint
	if _, err := fmt.Sscanf(id, "%d", &postID); err != nil {
		log.Printf("Invalid post ID: %s", id)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid post ID",
		})
		return
	}

	// Check if this IP has viewed this post in the last 30 minutes
	var existingView ViewRecord
	result := db.Where("post_id = ? AND client_ip = ? AND created_at > ?", postID, clientIP, time.Now().Add(-30*time.Minute)).
		First(&existingView)

	// If no recent view found, increment counter and record view
	if result.Error != nil {
		log.Printf("No recent view found for post %d from IP %s, incrementing view count", postID, clientIP)

		// Record this view
		viewRecord := ViewRecord{
			PostID:    postID,
			ClientIP:  clientIP,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&viewRecord).Error; err != nil {
			log.Printf("Error creating view record: %v", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to record view",
			})
			return
		}

		// Increment view count
		if err := db.Model(&Post{}).Where("id = ?", postID).Update("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
			log.Printf("Error updating view count: %v", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to update view count",
			})
			return
		}

		// Get the updated post to return the new view count
		var updatedPost Post
		if err := db.Where("id = ?", postID).First(&updatedPost).Error; err != nil {
			log.Printf("Error fetching updated post: %v", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch updated post",
			})
			return
		}

		log.Printf("View count incremented for post %d, new count: %d", postID, updatedPost.ViewCount)

		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "View count updated",
			Data: gin.H{
				"post_id":    postID,
				"view_count": updatedPost.ViewCount,
			},
		})
	} else {
		log.Printf("Recent view found for post %d from IP %s (viewed at %s), not incrementing", postID, clientIP, existingView.CreatedAt)

		// Even if view wasn't incremented, return current view count
		var currentPost Post
		if err := db.Where("id = ?", postID).First(&currentPost).Error; err != nil {
			log.Printf("Error fetching current post: %v", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to fetch current post",
			})
			return
		}

		log.Printf("Returning current view count for post %d: %d", postID, currentPost.ViewCount)

		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "View already recorded",
			Data: gin.H{
				"post_id":    postID,
				"view_count": currentPost.ViewCount,
			},
		})
	}
}

// Admin post handlers
func getAdminPosts(c *gin.Context) {
	page, limit, offset := GetPaginationParams(c)
	search := GetSearchQuery(c)
	status := GetStatusFilter(c)

	var posts []Post
	query := db.Model(&Post{}).Preload("Author").Order("created_at DESC")

	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	query.Offset(offset).Limit(limit).Find(&posts)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"posts":       posts,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func createPost(c *gin.Context) {
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	userID := GetCurrentUserID(c)
	slug := GenerateSlug(req.Title)

	post := Post{
		Title:      req.Title,
		Slug:       slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		Status:     req.Status,
		Categories: req.Categories,
		Tags:       req.Tags,
		ImageURL:   req.ImageURL,
		AuthorID:   userID,
	}

	if req.Status == "published" {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create post",
		})
		return
	}

	// Eğer yazı yayınlandıysa abonelere e-posta gönder
	if req.Status == "published" {
		go func() {
			subject := fmt.Sprintf("🚀 Yeni Blog Yazısı: %s", post.Title)
			content := generateNewPostEmailContent(post)
			if err := sendEmailToSubscribers(subject, content); err != nil {
				log.Printf("Failed to send newsletter emails: %v", err)
			}
		}()
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Post created successfully",
		Data:    post,
	})
}

func getAdminPost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	if err := db.Preload("Author").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    post,
	})
}

func updatePost(c *gin.Context) {
	id := c.Param("id")
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var post Post
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Post not found",
		})
		return
	}

	// Update fields
	post.Title = req.Title
	post.Slug = GenerateSlug(req.Title)
	post.Content = req.Content
	post.Excerpt = req.Excerpt
	post.Status = req.Status
	post.Categories = req.Categories
	post.Tags = req.Tags
	post.ImageURL = req.ImageURL

	if req.Status == "published" && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Post updated successfully",
		Data:    post,
	})
}

func deletePost(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Post{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

// Public project handlers
func getProjects(c *gin.Context) {
	var projects []Project
	db.Order("created_at DESC").Find(&projects)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    projects,
	})
}

func getProject(c *gin.Context) {
	slug := c.Param("slug")
	var project Project

	if err := db.Where("slug = ?", slug).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    project,
	})
}

// Admin project handlers
func getAdminProjects(c *gin.Context) {
	var projects []Project
	result := db.Order("created_at DESC").Find(&projects)

	log.Printf("getAdminProjects - Found %d projects, error: %v", len(projects), result.Error)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"projects":    projects,
			"total":       int64(len(projects)),
			"page":        1,
			"limit":       10,
			"total_pages": 1,
		},
	})
}

func createProject(c *gin.Context) {
	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	slug := GenerateSlug(req.Title)

	project := Project{
		Title:        req.Title,
		Slug:         slug,
		Description:  req.Description,
		Content:      req.Content,
		Status:       req.Status,
		Technologies: req.Technologies,
		GithubURL:    req.GithubURL,
		DemoURL:      req.DemoURL,
		ImageURL:     req.ImageURL,
		StartDate:    time.Now(),
	}

	if err := db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create project",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Project created successfully",
		Data:    project,
	})
}

func getAdminProject(c *gin.Context) {
	id := c.Param("id")
	var project Project

	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Project not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    project,
	})
}

func updateProject(c *gin.Context) {
	id := c.Param("id")
	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var project Project
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Project not found",
		})
		return
	}

	// Update fields
	project.Title = req.Title
	project.Slug = GenerateSlug(req.Title)
	project.Description = req.Description
	project.Content = req.Content
	project.Status = req.Status
	project.Technologies = req.Technologies
	project.GithubURL = req.GithubURL
	project.DemoURL = req.DemoURL
	project.ImageURL = req.ImageURL

	if req.Status == "completed" && project.EndDate == nil {
		now := time.Now()
		project.EndDate = &now
	}

	if err := db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update project",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Project updated successfully",
		Data:    project,
	})
}

func deleteProject(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete project",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Project deleted successfully",
	})
}

// Message handlers
func createMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	message := Message{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to send message",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    message,
	})
}

func getMessages(c *gin.Context) {
	page, limit, offset := GetPaginationParams(c)
	search := GetSearchQuery(c)

	var messages []Message
	query := db.Model(&Message{}).Order("created_at DESC")

	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR subject LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&messages)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"messages":    messages,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func getMessage(c *gin.Context) {
	id := c.Param("id")
	var message Message

	if err := db.First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Message not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    message,
	})
}

func updateMessage(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		IsRead    bool `json:"is_read"`
		IsReplied bool `json:"is_replied"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var message Message
	if err := db.First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Message not found",
		})
		return
	}

	message.IsRead = req.IsRead
	message.IsReplied = req.IsReplied

	if err := db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update message",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Message updated successfully",
		Data:    message,
	})
}

func deleteMessage(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Message{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete message",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Message deleted successfully",
	})
}

// Settings handlers
func getSettings(c *gin.Context) {
	var settings Settings
	if err := db.First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Settings not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    settings,
	})
}

func updateSettings(c *gin.Context) {
	var req Settings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var settings Settings
	if err := db.First(&settings).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Settings not found",
		})
		return
	}

	// Update fields
	settings.SiteTitle = req.SiteTitle
	settings.SiteTagline = req.SiteTagline
	settings.SiteDescription = req.SiteDescription
	settings.Email = req.Email
	settings.LinkedIn = req.LinkedIn
	settings.Github = req.Github
	settings.Location = req.Location
	settings.Company = req.Company
	settings.Position = req.Position
	settings.AboutContent = req.AboutContent

	if err := db.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update settings",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Settings updated successfully",
		Data:    settings,
	})
}

// User profile handlers
func getUserProfile(c *gin.Context) {
	userID := GetCurrentUserID(c)
	var user User

	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    user,
	})
}

func updateUserProfile(c *gin.Context) {
	userID := GetCurrentUserID(c)
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request data",
		})
		return
	}

	var user User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		// Validate password strength
		if err := ValidatePasswordStrength(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		hashedPassword, err := HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to hash password",
			})
			return
		}
		user.Password = hashedPassword
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    user,
	})
}

// SEO handlers
func generateSitemap(c *gin.Context) {
	// Get all published posts
	var posts []Post
	db.Where("status = ?", "published").Find(&posts)

	// Get all active projects
	var projects []Project
	db.Where("status IN ?", []string{"active", "completed"}).Find(&projects)

	// Generate sitemap XML
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>https://sametegeasik.com/</loc>
        <lastmod>` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>daily</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>https://sametegeasik.com/yazilar</loc>
        <lastmod>` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>daily</changefreq>
        <priority>0.9</priority>
    </url>
    <url>
        <loc>https://sametegeasik.com/projeler</loc>
        <lastmod>` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>weekly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>https://sametegeasik.com/hakkimda</loc>
        <lastmod>` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.7</priority>
    </url>
    <url>
        <loc>https://sametegeasik.com/iletisim</loc>
        <lastmod>` + time.Now().Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.6</priority>
    </url>`

	// Add blog posts to sitemap
	for _, post := range posts {
		sitemap += `
    <url>
        <loc>https://sametegeasik.com/posts/` + post.Slug + `</loc>
        <lastmod>` + post.UpdatedAt.Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.8</priority>
    </url>`
	}

	// Add projects to sitemap
	for _, project := range projects {
		sitemap += `
    <url>
        <loc>https://sametegeasik.com/projects/` + project.Slug + `</loc>
        <lastmod>` + project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00") + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.7</priority>
    </url>`
	}

	sitemap += `
</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

// Coming Soon handlers
func getComingSoon(c *gin.Context) {
	var comingSoonPosts []ComingSoon
	db.Order("priority DESC, estimated_date ASC").Find(&comingSoonPosts)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    comingSoonPosts,
	})
}

func getComingSoonByID(c *gin.Context) {
	id := c.Param("id")
	var comingSoonPost ComingSoon

	if err := db.First(&comingSoonPost, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Coming soon post not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    comingSoonPost,
	})
}

func createComingSoon(c *gin.Context) {
	var request ComingSoonRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request format: %v", err),
		})
		return
	}

	// Debug: Log the parsed request
	log.Printf("Parsed request: %+v", request)

	// Parse estimated date
	var estimatedDate *time.Time
	if request.EstimatedDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", request.EstimatedDate); err == nil {
			estimatedDate = &parsedDate
		}
	}

	comingSoonPost := ComingSoon{
		Title:         request.Title,
		Description:   request.Description,
		Category:      request.Category,
		Tags:          request.Tags,
		Status:        request.Status,
		Priority:      request.Priority,
		EstimatedDate: estimatedDate,
	}

	if err := db.Create(&comingSoonPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to create coming soon post",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Coming soon post created successfully",
		Data:    comingSoonPost,
	})
}

func updateComingSoon(c *gin.Context) {
	id := c.Param("id")
	var comingSoonPost ComingSoon

	if err := db.First(&comingSoonPost, id).Error; err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "Coming soon post not found",
		})
		return
	}

	var request ComingSoonRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Parse estimated date
	var estimatedDate *time.Time
	if request.EstimatedDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000Z", request.EstimatedDate); err == nil {
			estimatedDate = &parsedDate
		}
	}

	comingSoonPost.Title = request.Title
	comingSoonPost.Description = request.Description
	comingSoonPost.Category = request.Category
	comingSoonPost.Tags = request.Tags
	comingSoonPost.Status = request.Status
	comingSoonPost.Priority = request.Priority
	comingSoonPost.EstimatedDate = estimatedDate

	if err := db.Save(&comingSoonPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update coming soon post",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Coming soon post updated successfully",
		Data:    comingSoonPost,
	})
}

func deleteComingSoon(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&ComingSoon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete coming soon post",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Coming soon post deleted successfully",
	})
}

// Newsletter handlers
func getNewsletterSubscriptions(c *gin.Context) {
	page, limit, offset := GetPaginationParams(c)
	search := GetSearchQuery(c)

	var subscriptions []Newsletter
	query := db.Model(&Newsletter{}).Order("created_at DESC")

	if search != "" {
		query = query.Where("email LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&subscriptions)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"subscriptions": subscriptions,
			"total":         total,
			"page":          page,
			"limit":         limit,
			"total_pages":   (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func deleteNewsletterSubscription(c *gin.Context) {
	id := c.Param("id")

	if err := db.Delete(&Newsletter{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to delete newsletter subscription",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Newsletter subscription deleted successfully",
	})
}

// Email service functions
func sendEmailToSubscribers(subject, content string) error {
	var subscribers []Newsletter
	if err := db.Where("is_active = ?", true).Find(&subscribers).Error; err != nil {
		return err
	}

	if len(subscribers) == 0 {
		log.Println("No active subscribers found")
		return nil
	}

	// E-posta yapılandırması (çevre değişkenlerinden)
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")

	// Varsayılan değerler (test için)
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	if fromEmail == "" {
		fromEmail = "noreply@sametegeasik.com"
	}

	// SMTP yapılandırması eksikse e-posta gönderme
	log.Printf("SMTP Config Debug: HOST=%s, PORT=%s, USER=%s, PASS=%s, FROM=%s", smtpHost, smtpPort, smtpUser, smtpPass, fromEmail)
	if smtpUser == "" || smtpPass == "" {
		log.Println("SMTP credentials not configured, skipping email sending")
		return nil
	}

	// Her aboneye e-posta gönder
	for _, subscriber := range subscribers {
		if err := sendEmail(subscriber.Email, subject, content, smtpHost, smtpPort, smtpUser, smtpPass, fromEmail); err != nil {
			log.Printf("Failed to send email to %s: %v", subscriber.Email, err)
		} else {
			log.Printf("Email sent successfully to %s", subscriber.Email)
		}
	}

	return nil
}

func sendEmail(to, subject, content, smtpHost, smtpPort, smtpUser, smtpPass, fromEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	port := 587
	if smtpPort == "465" {
		port = 465
	}

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)

	return d.DialAndSend(m)
}

func generateNewPostEmailContent(post Post) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Yeni Blog Yazısı</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; }
			.content { background: #f9f9f9; padding: 20px; }
			.post-title { color: #333; font-size: 24px; margin-bottom: 10px; }
			.post-excerpt { color: #666; font-size: 16px; margin-bottom: 20px; }
			.btn { display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; }
			.footer { text-align: center; padding: 20px; color: #666; font-size: 14px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>🚀 Yeni Blog Yazısı Yayınlandı!</h1>
				<p>Samet Ege Aşık - Java Backend Developer</p>
			</div>
			<div class="content">
				<h2 class="post-title">%s</h2>
				<p class="post-excerpt">%s</p>
				<p>Yazının tamamını okumak için aşağıdaki bağlantıya tıklayın:</p>
				<a href="https://sametegeasik.com/posts/%s" class="btn">Yazıyı Oku</a>
			</div>
			<div class="footer">
				<p>Bu e-postayı aldınız çünkü blog yazılarımıza abone oldunuz.</p>
				<p>Aboneliğinizi iptal etmek için <a href="https://sametegeasik.com/unsubscribe">buraya tıklayın</a>.</p>
			</div>
		</div>
	</body>
	</html>
	`, post.Title, post.Excerpt, post.Slug)
}

// Load environment variables from .env file
func loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		// Read .env file
		content, err := os.ReadFile(envFile)
		if err != nil {
			log.Printf("Warning: Could not read .env file: %v", err)
			return
		}

		// Parse and set environment variables
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					os.Setenv(key, value)
				}
			}
		}
	}
}

// Get base directory from environment or use default
func getBaseDir() string {
	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		baseDir = "/Users/sametegeasik/blog-sitesi"
	}
	return baseDir
}

// Newsletter subscription handler
func subscribeNewsletter(c *gin.Context) {
	var req NewsletterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Geçersiz e-posta adresi",
		})
		return
	}

	// Check if email already exists
	var existingSubscription Newsletter
	if err := db.Where("email = ?", req.Email).First(&existingSubscription).Error; err == nil {
		if existingSubscription.IsActive {
			c.JSON(http.StatusConflict, APIResponse{
				Success: false,
				Message: "Bu e-posta adresi zaten abone listesinde",
			})
			return
		} else {
			// Reactivate existing subscription
			existingSubscription.IsActive = true
			if err := db.Save(&existingSubscription).Error; err != nil {
				c.JSON(http.StatusInternalServerError, APIResponse{
					Success: false,
					Message: "Abonelik güncellenirken hata oluştu",
				})
				return
			}
			c.JSON(http.StatusOK, APIResponse{
				Success: true,
				Message: "Aboneliğiniz yeniden aktif edildi! 🎉",
			})
			return
		}
	}

	// Create new subscription
	newsletter := Newsletter{
		Email:    req.Email,
		IsActive: true,
	}

	if err := db.Create(&newsletter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Abonelik oluşturulurken hata oluştu",
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Başarıyla abone oldunuz! 🚀",
		Data:    newsletter,
	})
}
