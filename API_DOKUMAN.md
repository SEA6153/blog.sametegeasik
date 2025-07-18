# 🔌 API Dokümantasyonu

## 📍 Base URL
```
http://localhost:8081/api
```

## 🔐 Kimlik Doğrulama

### Giriş Yapma
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Yanıt:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com"
  }
}
```

### Token Yenileme
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

## 📝 Blog Yazıları

### Tüm Yazıları Getir
```http
GET /api/posts
```

### Tek Yazı Getir
```http
GET /api/posts/{slug}
```

### Yeni Yazı Ekle (Auth Gerekli)
```http
POST /api/admin/posts
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "Yeni Blog Yazısı",
  "content": "Yazı içeriği...",
  "excerpt": "Kısa özet",
  "status": "published",
  "categories": "Teknoloji",
  "tags": "go,web,api"
}
```

### Yazı Güncelle (Auth Gerekli)
```http
PUT /api/admin/posts/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "Güncellenmiş Başlık",
  "content": "Güncellenmiş içerik..."
}
```

### Yazı Sil (Auth Gerekli)
```http
DELETE /api/admin/posts/{id}
Authorization: Bearer {access_token}
```

## 🚀 Projeler

### Tüm Projeleri Getir
```http
GET /api/projects
```

### Tek Proje Getir
```http
GET /api/projects/{slug}
```

### Yeni Proje Ekle (Auth Gerekli)
```http
POST /api/admin/projects
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "title": "Yeni Proje",
  "description": "Proje açıklaması",
  "content": "Detaylı içerik",
  "status": "active",
  "technologies": "Go,React,SQLite",
  "github_url": "https://github.com/user/project",
  "demo_url": "https://demo.com"
}
```

## 📧 Mesajlar

### Mesaj Gönder
```http
POST /api/messages
Content-Type: application/json

{
  "name": "Ad Soyad",
  "email": "email@example.com",
  "subject": "Konu",
  "message": "Mesaj içeriği"
}
```

### Mesajları Listele (Auth Gerekli)
```http
GET /api/admin/messages
Authorization: Bearer {access_token}
```

## 📊 Dashboard

### İstatistikleri Getir (Auth Gerekli)
```http
GET /api/dashboard/stats
Authorization: Bearer {access_token}
```

**Yanıt:**
```json
{
  "total_posts": 5,
  "total_projects": 3,
  "unread_messages": 2,
  "total_views": 150,
  "published_posts": 4,
  "draft_posts": 1,
  "active_projects": 2,
  "completed_projects": 1,
  "newsletter_subscribers": 10
}
```

## ⚙️ Site Ayarları

### Ayarları Getir
```http
GET /api/settings
```

### Ayarları Güncelle (Auth Gerekli)
```http
PUT /api/admin/settings
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "site_title": "Site Başlığı",
  "site_tagline": "Site Alt Başlığı",
  "email": "contact@example.com",
  "linkedin": "https://linkedin.com/in/user",
  "github": "https://github.com/user"
}
```

## 📧 Newsletter

### Abone Ol
```http
POST /api/newsletter
Content-Type: application/json

{
  "email": "user@example.com"
}
```

## 🔍 Hata Kodları

| Kod | Açıklama |
|-----|----------|
| 200 | Başarılı |
| 400 | Hatalı istek |
| 401 | Kimlik doğrulama gerekli |
| 403 | Yetkisiz erişim |
| 404 | Bulunamadı |
| 500 | Sunucu hatası |

## 📝 Örnek Kullanım

### JavaScript ile API Kullanımı
```javascript
// Giriş yapma
const login = async () => {
  const response = await fetch('/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      username: 'admin',
      password: 'admin123'
    })
  });
  
  const data = await response.json();
  localStorage.setItem('token', data.access_token);
};

// Blog yazılarını getirme
const getPosts = async () => {
  const response = await fetch('/api/posts');
  const posts = await response.json();
  return posts;
};

// Yeni yazı ekleme
const createPost = async (postData) => {
  const token = localStorage.getItem('token');
  const response = await fetch('/api/admin/posts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(postData)
  });
  
  return await response.json();
};
```

## 🔒 Güvenlik Notları

- Tüm admin endpoint'leri JWT token gerektirir
- Token'lar 24 saat geçerlidir
- Rate limiting aktif
- Input validation zorunlu
- CORS koruması mevcut

---

**Not:** Bu API dokümantasyonu temel kullanım için hazırlanmıştır. Detaylı bilgiler için kaynak kodları inceleyebilirsiniz. 