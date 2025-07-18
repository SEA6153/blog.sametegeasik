package main

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// User model for admin authentication
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" hides password from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Post model for blog posts
type Post struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Slug        string     `json:"slug" gorm:"unique;not null"`
	Content     string     `json:"content" gorm:"type:text"`
	Excerpt     string     `json:"excerpt" gorm:"type:text"`
	Status      string     `json:"status" gorm:"default:draft"` // draft, published
	Categories  string     `json:"categories"`
	Tags        string     `json:"tags"`
	ImageURL    string     `json:"image_url"`
	ViewCount   int        `json:"view_count" gorm:"default:0"`
	AuthorID    uint       `json:"author_id"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`
}

// Project model for portfolio projects
type Project struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Title        string     `json:"title" gorm:"not null"`
	Slug         string     `json:"slug" gorm:"unique;not null"`
	Description  string     `json:"description" gorm:"type:text"`
	Content      string     `json:"content" gorm:"type:text"`
	Status       string     `json:"status" gorm:"default:active"` // active, completed, archived
	Technologies string     `json:"technologies"`                 // JSON string of tech stack
	GithubURL    string     `json:"github_url"`
	DemoURL      string     `json:"demo_url"`
	ImageURL     string     `json:"image_url"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Message model for contact form messages
type Message struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Subject   string    `json:"subject" gorm:"not null"`
	Message   string    `json:"message" gorm:"type:text;not null"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	IsReplied bool      `json:"is_replied" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ComingSoon model for upcoming blog posts
type ComingSoon struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description" gorm:"type:text"`
	Category      string     `json:"category" gorm:"not null"`
	Tags          string     `json:"tags"`
	Status        string     `json:"status" gorm:"default:planned"` // planned, in_progress, completed
	Priority      int        `json:"priority" gorm:"default:1"`     // 1-5, 1 being highest priority
	EstimatedDate *time.Time `json:"estimated_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Settings model for site configuration
type Settings struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	SiteTitle       string    `json:"site_title" gorm:"default:Samet Ege Aşık"`
	SiteTagline     string    `json:"site_tagline" gorm:"default:Java Backend Developer"`
	SiteDescription string    `json:"site_description" gorm:"type:text"`
	Email           string    `json:"email"`
	LinkedIn        string    `json:"linkedin"`
	Github          string    `json:"github"`
	Location        string    `json:"location"`
	Company         string    `json:"company"`
	Position        string    `json:"position"`
	AboutContent    string    `json:"about_content" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Newsletter subscription model
type Newsletter struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null;unique"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ViewRecord model for tracking post views
type ViewRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"index"`
	ClientIP  string    `json:"client_ip" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}

// Response structures for API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type PostRequest struct {
	Title      string `json:"title" binding:"required,min=3,max=255"`
	Content    string `json:"content" binding:"required,min=10,max=50000"`
	Excerpt    string `json:"excerpt" binding:"max=500"`
	Status     string `json:"status" binding:"oneof=draft published"`
	Categories string `json:"categories" binding:"max=255"`
	Tags       string `json:"tags" binding:"max=255"`
	ImageURL   string `json:"image_url" binding:"omitempty,url,max=500"`
}

type ProjectRequest struct {
	Title        string `json:"title" binding:"required,min=3,max=255"`
	Description  string `json:"description" binding:"required,min=10,max=1000"`
	Content      string `json:"content" binding:"max=50000"`
	Status       string `json:"status" binding:"oneof=active completed archived"`
	Technologies string `json:"technologies" binding:"max=500"`
	GithubURL    string `json:"github_url" binding:"omitempty,url,max=500"`
	DemoURL      string `json:"demo_url" binding:"omitempty,url,max=500"`
	ImageURL     string `json:"image_url" binding:"omitempty,url,max=500"`
}

type MessageRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Email   string `json:"email" binding:"required,email,max=255"`
	Subject string `json:"subject" binding:"required,min=3,max=255"`
	Message string `json:"message" binding:"required,min=10,max=10000"`
}

type ComingSoonRequest struct {
	Title         string `json:"title" binding:"required,min=3,max=255"`
	Description   string `json:"description" binding:"required,min=10,max=1000"`
	Category      string `json:"category" binding:"required,min=2,max=100"`
	Tags          string `json:"tags" binding:"max=255"`
	Status        string `json:"status" binding:"oneof=planned in_progress completed"`
	Priority      int    `json:"priority" binding:"min=1,max=5"`
	EstimatedDate string `json:"estimated_date"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type NewsletterRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Dashboard stats response
type DashboardStats struct {
	TotalPosts            int64 `json:"total_posts"`
	TotalProjects         int64 `json:"total_projects"`
	UnreadMessages        int64 `json:"unread_messages"`
	TotalViews            int64 `json:"total_views"`
	PublishedPosts        int64 `json:"published_posts"`
	DraftPosts            int64 `json:"draft_posts"`
	ActiveProjects        int64 `json:"active_projects"`
	CompletedProjects     int64 `json:"completed_projects"`
	NewsletterSubscribers int64 `json:"newsletter_subscribers"`
}

// Initialize database with sample data
func InitDB(db *gorm.DB) {
	// Auto-migrate schemas
	db.AutoMigrate(&User{}, &Post{}, &Project{}, &Message{}, &Settings{}, &ComingSoon{}, &Newsletter{}, &ViewRecord{})

	// Create default admin user if not exists
	var adminUser User
	if err := db.Where("username = ?", "sea6153").First(&adminUser).Error; err != nil {
		// Hash password
		hashedPassword, _ := HashPassword("Trabzonspor1967*")
		defaultAdmin := User{
			Username: "sea6153",
			Email:    "blog.sametegeasik@gmail.com",
			Password: hashedPassword,
		}
		result := db.Create(&defaultAdmin)
		if result.Error != nil {
			log.Printf("Error creating admin user: %v", result.Error)
		} else {
			log.Printf("Admin user created successfully with ID: %d", defaultAdmin.ID)
		}
	} else {
		log.Printf("Admin user already exists with ID: %d", adminUser.ID)
	}

	// Create default settings if not exists
	var settings Settings
	if err := db.First(&settings).Error; err != nil {
		defaultSettings := Settings{
			SiteTitle:       "Samet Ege Aşık",
			SiteTagline:     "Java Backend Developer",
			SiteDescription: "Jeoloji mühendisliğinden yazılım geliştirmeye geçiş yapan, OAuth2, mikroservis mimarisi ve modern backend teknolojileri konusunda uzman Java Backend Developer.",
			Email:           "blog.sametegeasik@gmail.com",
			LinkedIn:        "https://linkedin.com/in/samet-ege-asik",
			Github:          "https://github.com/sametegeasik",
			Location:        "Ankara / Çankaya",
			Company:         "T.C. Millî Eğitim Bakanlığı YEĞİTEK",
			Position:        "Java Backend Developer",
			AboutContent:    "Hacettepe Üniversitesi Jeoloji Mühendisliği mezunu olarak kariyerime başladım. Alan dışı bir bölümden yazılım dünyasına geçiş yaparak, BilgeAdam Academy'de Java eğitimi aldım. Şu anda T.C. Millî Eğitim Bakanlığı YEĞİTEK'te Java Backend Developer olarak çalışıyorum.",
		}
		db.Create(&defaultSettings)
	}
}
