# 🚀 Samet Ege Aşık - Blog Sitesi & Backend API

Modern, responsif blog sitesi ve Go ile geliştirilmiş REST API backend'i. OAuth2 kimlik doğrulama, admin paneli ve dinamik içerik yönetimi ile tam özellikli blog platformu.

## 📋 Özellikler

### 🎨 Frontend
- **Modern UI/UX**: Tailwind CSS ile responsive tasarım
- **Çoklu Sayfa**: Ana sayfa, blog yazıları, projeler, hakkımda, iletişim
- **Admin Paneli**: Tam özellikli yönetim arayüzü
- **Dinamik İçerik**: Backend API'den canlı veri çekme
- **Form Validasyonu**: İletişim formu ve admin paneli
- **Mobil Uyumlu**: Tüm cihazlarda optimum görüntüleme

### 🔧 Backend API
- **Go + Gin Framework**: Yüksek performanslı REST API
- **JWT Authentication**: Güvenli kimlik doğrulama
- **SQLite Database**: Hafif ve taşınabilir veritabanı
- **GORM ORM**: Kolay veritabanı yönetimi
- **CORS Support**: Frontend entegrasyonu
- **Auto-Migration**: Otomatik veritabanı schema yönetimi

### 🛡️ Güvenlik
- **JWT Tokens**: Güvenli oturum yönetimi
- **Password Hashing**: bcrypt ile şifre şifreleme
- **Input Validation**: Gin validator ile veri doğrulama
- **CORS Protection**: Güvenli API erişimi

## 🚀 Kurulum

### Backend Kurulumu

1. **Go Kurulumu** (v1.19+)
```bash
# macOS için
brew install go

# Linux için
sudo apt-get install golang-go

# Windows için Go web sitesinden indirin
```

2. **Backend Başlatma**
```bash
cd backend
go mod tidy
go run .
```

Backend `http://localhost:8080` adresinde çalışacak.

### Frontend Kurulumu

1. **HTML Dosyalarını Açma**
```bash
# Ana sayfayı açma
open index.html

# Admin paneli
open admin.html

# İletişim sayfası
open iletisim.html
```

2. **HTTP Server ile Çalıştırma** (Önerilen)
```bash
# Python 3
python -m http.server 8000

# Node.js
npx http-server

# PHP
php -S localhost:8000
```

## 🎯 Kullanım

### Admin Paneli Giriş
- **Kullanıcı Adı**: `admin`
- **Şifre**: `admin123`
- **URL**: `admin.html`

### API Endpoints

#### 🔐 Authentication
- `POST /api/auth/login` - Kullanıcı girişi

#### 📝 Blog Yazıları
- `GET /api/posts` - Tüm yazıları listele
- `GET /api/posts/:slug` - Tek yazı detayı
- `POST /api/admin/posts` - Yeni yazı oluştur (Auth)
- `PUT /api/admin/posts/:id` - Yazı güncelle (Auth)
- `DELETE /api/admin/posts/:id` - Yazı sil (Auth)

#### 🚀 Projeler
- `GET /api/projects` - Tüm projeleri listele
- `GET /api/projects/:slug` - Tek proje detayı
- `POST /api/admin/projects` - Yeni proje oluştur (Auth)
- `PUT /api/admin/projects/:id` - Proje güncelle (Auth)
- `DELETE /api/admin/projects/:id` - Proje sil (Auth)

#### 📧 Mesajlar
- `POST /api/messages` - Yeni mesaj gönder
- `GET /api/admin/messages` - Tüm mesajları listele (Auth)
- `PUT /api/admin/messages/:id` - Mesaj güncelle (Auth)
- `DELETE /api/admin/messages/:id` - Mesaj sil (Auth)

#### 📊 Dashboard
- `GET /api/dashboard/stats` - Dashboard istatistikleri (Auth)

#### ⚙️ Ayarlar
- `GET /api/settings` - Site ayarlarını al
- `PUT /api/admin/settings` - Ayarları güncelle (Auth)

## 📁 Proje Yapısı

```
blog-sitesi/
├── backend/
│   ├── main.go          # Ana server dosyası
│   ├── models.go        # Veritabanı modelleri
│   ├── auth.go          # Authentication fonksiyonları
│   ├── seed.go          # Örnek veri oluşturma
│   ├── go.mod           # Go modül dosyası
│   └── blog.db          # SQLite veritabanı
├── index.html           # Ana sayfa
├── yazilar.html         # Blog yazıları
├── projeler.html        # Projeler sayfası
├── hakkimda.html        # Hakkımda sayfası
├── iletisim.html        # İletişim sayfası
├── admin.html           # Admin paneli
└── README.md            # Bu dosya
```

## 🔧 Geliştirme

### Backend Geliştirme
```bash
cd backend
go mod tidy
go run .
```

### Veritabanı Sıfırlama
```bash
cd backend
rm blog.db
go run .
```

## 🌟 Teknolojiler

### Backend
- **Go 1.19+** - Programlama dili
- **Gin Web Framework** - HTTP router
- **GORM** - ORM kütüphanesi
- **SQLite** - Veritabanı
- **JWT** - Kimlik doğrulama
- **bcrypt** - Şifre şifreleme
- **CORS** - Cross-origin resource sharing

### Frontend
- **HTML5** - Markup
- **Tailwind CSS** - Styling framework
- **JavaScript (ES6+)** - Dinamik işlevsellik
- **Font Awesome** - İkonlar
- **Fetch API** - HTTP istekleri

## 📱 Responsive Tasarım

Site tüm cihazlarda optimize edilmiştir:
- **Desktop**: 1200px+
- **Tablet**: 768px-1199px
- **Mobile**: 320px-767px

## 🔐 Güvenlik Özellikleri

- JWT token tabanlı kimlik doğrulama
- Password hashing (bcrypt)
- Input validation ve sanitization
- CORS protection
- SQL injection koruması (GORM)

## 🚀 Deployment

### Backend Deployment
```bash
# Binary oluşturma
go build -o blog-backend

# Çalıştırma
./blog-backend
```

### Frontend Deployment
- Static hosting (Netlify, Vercel, GitHub Pages)
- Backend URL'ini production API'ye güncelleyin

## 📊 Örnek Veriler

Backend ilk çalıştırıldığında otomatik olarak örnek veriler oluşturulur:
- 5 blog yazısı
- 4 proje
- 5 iletişim mesajı
- 1 admin kullanıcısı

## 🤝 Katkıda Bulunma

1. Fork edin
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📧 İletişim

- **E-posta**: blog.sametegeasik@gmail.com
- **LinkedIn**: [linkedin.com/in/samet-ege-asik](https://linkedin.com/in/samet-ege-asik)
- **GitHub**: [github.com/sametegeasik](https://github.com/sametegeasik)

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

---

**Samet Ege Aşık** - Java Backend Developer @ T.C. Millî Eğitim Bakanlığı YEĞİTEK 