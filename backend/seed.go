package main

import (
	"time"

	"gorm.io/gorm"
)

// Seed database with sample data
func SeedDatabase(db *gorm.DB) {
	// Create sample posts
	seedPosts(db)

	// Create sample projects
	seedProjects(db)

	// Create sample messages
	seedMessages(db)

	// Create sample coming soon posts
	seedComingSoonPosts(db)
}

func seedPosts(db *gorm.DB) {
	// Check if posts already exist
	var count int64
	db.Model(&Post{}).Count(&count)
	if count > 0 {
		return
	}

	// Get admin user
	var admin User
	db.Where("username = ?", "sea6153").First(&admin)

	posts := []Post{
		{
			Title:       "Spring Authorization Server ile SSO Entegrasyonu",
			Slug:        "spring-authorization-server-ile-sso-entegrasyonu",
			Content:     "OAuth2 ve OpenID Connect protokollerini kullanarak Single Sign-On (SSO) entegrasyonu nasıl kurulur? Bu yazıda Spring Authorization Server kullanarak modern kimlik doğrulama sistemlerini inceleyeceğiz.",
			Excerpt:     "OAuth2 tabanlı kimlik doğrulama sistemi kurulumu ve konfigürasyonu",
			Status:      "published",
			Categories:  "OAuth2,Security,Spring",
			Tags:        "spring,oauth2,security,sso",
			ViewCount:   245,
			AuthorID:    admin.ID,
			PublishedAt: &time.Time{},
		},
		{
			Title:       "Redis ile API Performansı Nasıl %30 İyileştirilir?",
			Slug:        "redis-ile-api-performansi-nasil-iyilestirilir",
			Content:     "Redis kullanarak API performansınızı önemli ölçüde artırabilirsiniz. Caching stratejileri, best practices ve gerçek dünya örnekleri ile performans optimizasyonu.",
			Excerpt:     "Redis caching stratejileri ve performans optimizasyonu teknikleri",
			Status:      "published",
			Categories:  "Performance,Caching,Redis",
			Tags:        "redis,performance,caching,optimization",
			ViewCount:   189,
			AuthorID:    admin.ID,
			PublishedAt: &time.Time{},
		},
		{
			Title:      "Kafka ile Mikroservisler Arası Gerçek Zamanlı İletişim",
			Slug:       "kafka-ile-mikroservisler-arasi-gercek-zamanli-iletisim",
			Content:    "Apache Kafka kullanarak mikroservisler arasında güvenilir, ölçeklenebilir ve gerçek zamanlı iletişim kurmanın yolları.",
			Excerpt:    "Kafka ile event-driven architecture ve mikroservis iletişimi",
			Status:     "draft",
			Categories: "Microservices,Kafka,Architecture",
			Tags:       "kafka,microservices,messaging,architecture",
			ViewCount:  0,
			AuthorID:   admin.ID,
		},
		{
			Title:       "PostgreSQL Index ve Query Optimization İpuçları",
			Slug:        "postgresql-index-ve-query-optimization-ipuclari",
			Content:     "PostgreSQL veritabanında performanslı sorgular yazmak ve doğru indeks stratejileri kullanmak için pratik ipuçları.",
			Excerpt:     "PostgreSQL performans optimizasyonu ve indeks stratejileri",
			Status:      "published",
			Categories:  "Database,PostgreSQL,Performance",
			Tags:        "postgresql,database,performance,optimization",
			ViewCount:   156,
			AuthorID:    admin.ID,
			PublishedAt: &time.Time{},
		},
		{
			Title:       "Jeoloji Mühendisliğinden Yazılıma Geçiş: Nereden Başladım?",
			Slug:        "jeoloji-muhendisliginden-yazilima-gecis-nereden-basladim",
			Content:     "Alan dışı bir bölümden yazılım geliştirme dünyasına geçiş hikayem. Zorluklarla nasıl başa çıktım ve hangi yolları izledim?",
			Excerpt:     "Kariyer değişikliği hikayesi ve yazılım öğrenme süreci",
			Status:      "published",
			Categories:  "Career,Personal",
			Tags:        "career,geology,software,transition",
			ViewCount:   312,
			AuthorID:    admin.ID,
			PublishedAt: &time.Time{},
		},
	}

	for _, post := range posts {
		if post.Status == "published" {
			now := time.Now()
			post.PublishedAt = &now
		}
		db.Create(&post)
	}
}

func seedProjects(db *gorm.DB) {
	// Check if projects already exist
	var count int64
	db.Model(&Project{}).Count(&count)
	if count > 0 {
		return
	}

	projects := []Project{
		{
			Title:        "EBA SSO Integration Infrastructure",
			Slug:         "eba-sso-integration-infrastructure",
			Description:  "T.C. Millî Eğitim Bakanlığı için OAuth2 tabanlı kimlik doğrulama sistemi",
			Content:      "Spring Authorization Server kullanarak EBA platformu için Single Sign-On (SSO) entegrasyonu geliştirdim. Sistem Redis ile session yönetimi, Kafka ile event tracking ve Docker ile containerization kullanıyor.",
			Status:       "active",
			Technologies: `["Spring Boot", "Spring Authorization Server", "OAuth2", "Redis", "Kafka", "Docker", "PostgreSQL", "CI/CD"]`,
			GithubURL:    "",
			DemoURL:      "",
			StartDate:    time.Now().AddDate(0, -2, 0),
		},
		{
			Title:        "Forum Website Project",
			Slug:         "forum-website-project",
			Description:  "Kapsamlı forum platformu - Freelance proje",
			Content:      "Spring Boot ve Thymeleaf kullanarak tam özellikli forum sitesi geliştirdim. Kullanıcı yönetimi, konu/mesaj sistemi, moderasyon araçları ve admin paneli içeriyor.",
			Status:       "completed",
			Technologies: `["Spring Boot", "Thymeleaf", "PostgreSQL", "Spring Security", "Docker", "Kubernetes", "Maven"]`,
			GithubURL:    "https://github.com/sametegeasik/forum-website",
			DemoURL:      "",
			StartDate:    time.Now().AddDate(0, -6, 0),
			EndDate:      func() *time.Time { t := time.Now().AddDate(0, -1, 0); return &t }(),
		},
		{
			Title:        "E-Commerce Platform",
			Slug:         "e-commerce-platform",
			Description:  "BilgeAdam Academy mezuniyet projesi",
			Content:      "Spring Boot ile e-ticaret platformu geliştirdim. Ürün kataloğu, sepet yönetimi, sipariş takibi ve ödeme entegrasyonu özellikleri bulunuyor.",
			Status:       "completed",
			Technologies: `["Spring Boot", "Spring Data JPA", "MySQL", "Spring Security", "Bootstrap", "Maven"]`,
			GithubURL:    "https://github.com/sametegeasik/ecommerce-platform",
			DemoURL:      "",
			StartDate:    time.Now().AddDate(-1, 0, 0),
			EndDate:      func() *time.Time { t := time.Now().AddDate(0, -8, 0); return &t }(),
		},
		{
			Title:        "Entertainment QR Web App",
			Slug:         "entertainment-qr-web-app",
			Description:  "Kafe ve restoranlar için QR kod tabanlı eğlence uygulaması",
			Content:      "Müşterilerin QR kod ile erişebileceği basit web uygulaması. Rastgele kahve önerileri, oyunlar ve eğlenceli içerikler sunuyor.",
			Status:       "completed",
			Technologies: `["HTML", "CSS", "JavaScript", "Bootstrap", "QR Code API"]`,
			GithubURL:    "https://github.com/sametegeasik/qr-entertainment-app",
			DemoURL:      "https://entertainment-qr.netlify.app",
			StartDate:    time.Now().AddDate(0, -4, 0),
			EndDate:      func() *time.Time { t := time.Now().AddDate(0, -3, 0); return &t }(),
		},
	}

	for _, project := range projects {
		db.Create(&project)
	}
}

func seedMessages(db *gorm.DB) {
	// Check if messages already exist
	var count int64
	db.Model(&Message{}).Count(&count)
	if count > 0 {
		return
	}

	messages := []Message{
		{
			Name:    "Ahmet Yılmaz",
			Email:   "ahmet.yilmaz@example.com",
			Subject: "Spring Boot Projesi Hakkında",
			Message: "Merhaba Samet, Spring Boot ile ilgili yazdığınız yazıları takip ediyorum. Özellikle OAuth2 konusu çok faydalı. Bir proje için danışmanlık hizmeti alabilir miyiz?",
			IsRead:  false,
		},
		{
			Name:    "Fatma Kaya",
			Email:   "fatma.kaya@example.com",
			Subject: "Kariyer Değişikliği Tavsiyesi",
			Message: "Jeoloji mühendisliğinden yazılıma geçiş hikayenizi çok beğendim. Ben de benzer bir süreçteyim. Hangi kaynakları önerirsiniz?",
			IsRead:  false,
		},
		{
			Name:    "Mehmet Demir",
			Email:   "mehmet.demir@example.com",
			Subject: "İşbirliği Teklifi",
			Message: "Selam, LinkedIn profilinizi gördüm. Elimizde bir e-ticaret projesi var, backend tarafında yardıma ihtiyacımız var. Görüşebilir miyiz?",
			IsRead:  true,
		},
		{
			Name:    "Ayşe Özkan",
			Email:   "ayse.ozkan@example.com",
			Subject: "Blog Yazısı Teşekkürü",
			Message: "Redis ile ilgili yazınız çok açıklayıcıydı. Projemde uyguladım ve gerçekten performans artışı gözlemledim. Teşekkürler!",
			IsRead:  false,
		},
		{
			Name:    "Can Yıldız",
			Email:   "can.yildiz@example.com",
			Subject: "Mikroservis Mimarisi Sorusu",
			Message: "Kafka ile mikroservisler arası iletişim konusunda bir sorum var. Hangi durumda sync, hangi durumda async iletişim tercih ediyorsunuz?",
			IsRead:  true,
		},
	}

	for _, message := range messages {
		db.Create(&message)
	}
}

func seedComingSoonPosts(db *gorm.DB) {
	// Check if coming soon posts already exist
	var count int64
	db.Model(&ComingSoon{}).Count(&count)
	if count > 0 {
		return
	}

	// Create sample coming soon posts
	date1 := time.Now().AddDate(0, 0, 7)  // 1 hafta sonra
	date2 := time.Now().AddDate(0, 0, 14) // 2 hafta sonra
	date3 := time.Now().AddDate(0, 1, 0)  // 1 ay sonra
	date4 := time.Now().AddDate(0, 1, 15) // 1.5 ay sonra
	date5 := time.Now().AddDate(0, 2, 0)  // 2 ay sonra

	comingSoonPosts := []ComingSoon{
		{
			Title:         "Kubernetes ile Microservices Deployment",
			Description:   "Kubernetes üzerinde microservices mimarisinin nasıl kurulacağı, deployment stratejileri ve monitoring konularını ele alacağız.",
			Category:      "DevOps",
			Tags:          "kubernetes,microservices,devops,deployment",
			Status:        "planned",
			Priority:      5,
			EstimatedDate: &date1,
		},
		{
			Title:         "Apache Kafka ile Event-Driven Architecture",
			Description:   "Event-driven architecture patterns, Kafka ecosystem ve gerçek dünya use case'leri üzerinde durulacak.",
			Category:      "Architecture",
			Tags:          "kafka,event-driven,architecture,messaging",
			Status:        "in_progress",
			Priority:      4,
			EstimatedDate: &date2,
		},
		{
			Title:         "GraphQL ile Modern API Tasarımı",
			Description:   "GraphQL'in avantajları, Apollo Server kurulumu ve React entegrasyonu ile modern API geliştirme teknikleri.",
			Category:      "API",
			Tags:          "graphql,api,react,apollo",
			Status:        "planned",
			Priority:      3,
			EstimatedDate: &date3,
		},
		{
			Title:         "Machine Learning ile Öneri Sistemi",
			Description:   "Python ve TensorFlow kullanarak kullanıcı davranışlarına dayalı öneri algoritmaları geliştirme.",
			Category:      "AI/ML",
			Tags:          "machine-learning,tensorflow,python,recommendation",
			Status:        "planned",
			Priority:      2,
			EstimatedDate: &date4,
		},
		{
			Title:         "Blockchain ile Decentralized Identity",
			Description:   "Blockchain teknolojisi kullanarak merkezi olmayan kimlik doğrulama sistemlerinin nasıl kurulacağı.",
			Category:      "Blockchain",
			Tags:          "blockchain,identity,decentralized,security",
			Status:        "planned",
			Priority:      1,
			EstimatedDate: &date5,
		},
	}

	for _, post := range comingSoonPosts {
		db.Create(&post)
	}
}
