# 📖 Blog Sitesi Kullanım Kılavuzu

## 🚀 Hızlı Başlangıç

### Backend Kurulumu
```bash
cd backend
go run .
```
Backend `http://localhost:8081` adresinde çalışacak.

### Frontend Kullanımı
- Ana sayfa: `index.html`
- Admin paneli: `admin.html`
- Blog yazıları: `yazilar.html`
- Projeler: `projeler.html`
- Hakkımda: `hakkimda.html`
- İletişim: `iletisim.html`

## 🔐 Admin Paneli

### Admin Panel Özellikleri
- ✅ Blog yazısı ekleme/düzenleme/silme
- ✅ Proje ekleme/düzenleme/silme
- ✅ Gelen mesajları görüntüleme
- ✅ Site ayarlarını düzenleme
- ✅ Dashboard istatistikleri

## 📝 Blog Yazısı Yönetimi

### Yeni Yazı Ekleme
1. Admin paneline giriş yapın
2. "Yeni Yazı" butonuna tıklayın
3. Başlık, içerik ve diğer bilgileri doldurun
4. "Kaydet" butonuna tıklayın

### Yazı Düzenleme
1. Yazı listesinden düzenlemek istediğiniz yazıya tıklayın
2. Bilgileri güncelleyin
3. "Güncelle" butonuna tıklayın

## 🚀 Proje Yönetimi

### Yeni Proje Ekleme
1. Admin panelinde "Projeler" sekmesine gidin
2. "Yeni Proje" butonuna tıklayın
3. Proje bilgilerini doldurun
4. "Kaydet" butonuna tıklayın

## 📧 İletişim Mesajları

### Mesajları Görüntüleme
1. Admin panelinde "Mesajlar" sekmesine gidin
2. Gelen tüm mesajları görebilirsiniz
3. Mesajları okundu/yanıtlandı olarak işaretleyebilirsiniz

## ⚙️ Site Ayarları

### Ayarları Düzenleme
1. Admin panelinde "Ayarlar" sekmesine gidin
2. Site başlığı, açıklama, sosyal medya linkleri gibi bilgileri güncelleyin
3. "Kaydet" butonuna tıklayın

## 📊 Dashboard

### İstatistikler
Dashboard'da şu bilgileri görebilirsiniz:
- Toplam blog yazısı sayısı
- Toplam proje sayısı
- Okunmamış mesaj sayısı
- Toplam görüntülenme sayısı
- Newsletter abone sayısı

## 🔧 Sorun Giderme

### Backend Çalışmıyor
```bash
# Port kontrolü
lsof -i :8081

# Logları kontrol edin
tail -f backend/logs/app.log
```

### Veritabanı Sıfırlama
```bash
cd backend
rm blog.db
go run .
```

## 📱 Responsive Tasarım

Site tüm cihazlarda çalışır:
- **Desktop:** 1200px+
- **Tablet:** 768px-1199px
- **Mobile:** 320px-767px

## 🔒 Güvenlik

- JWT token tabanlı kimlik doğrulama
- Şifre hashleme
- Input validation
- CORS koruması

## 📞 Destek

Sorun yaşarsanız:
- **E-posta:** blog.sametegeasik@gmail.com
- **GitHub:** [github.com/sametegeasik](https://github.com/sametegeasik)

---

**Not:** Bu kılavuz temel kullanım için hazırlanmıştır. Detaylı teknik bilgiler için `README.md` dosyasını inceleyebilirsiniz. 