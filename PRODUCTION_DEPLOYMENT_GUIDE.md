# 🚀 Production Deployment Rehberi

## 📋 Ön Gereksinimler

### Sunucu Gereksinimleri:
- **OS**: Ubuntu 20.04+ / CentOS 8+ / Debian 11+
- **RAM**: Minimum 2GB (Önerilen: 4GB+)
- **CPU**: 2 Core (Önerilen: 4 Core+)
- **Disk**: 20GB+ boş alan
- **Domain**: sametegeasik.com (DNS ayarları yapılmış)

### Yazılım Gereksinimleri:
- **Docker** (önerilen) veya **Go 1.21+** + **Nginx**
- **SSL Sertifikası** (Let's Encrypt)
- **Git**

## 🐳 Docker ile Deployment (Önerilen)

### 1. Sunucuya Dosyaları Yükleyin
```bash
# Sunucuya SSH ile bağlanın
ssh user@your-server-ip

# Proje dizinini oluşturun
mkdir -p /var/www/blog-sitesi
cd /var/www/blog-sitesi

# Dosyaları yükleyin (scp, git clone veya rsync ile)
git clone https://github.com/SEA6153/blog.sametegeasik.git .
```

### 2. Docker Compose ile Başlatın
```bash
# Docker Compose ile başlatın
docker-compose up -d

# Logları kontrol edin
docker-compose logs -f
```

### 3. SSL Sertifikası Kurulumu
```bash
# Certbot kurulumu
sudo apt update
sudo apt install certbot python3-certbot-nginx

# SSL sertifikası alın
sudo certbot --nginx -d sametegeasik.com -d www.sametegeasik.com

# Otomatik yenileme için cron job ekleyin
sudo crontab -e
# Aşağıdaki satırı ekleyin:
# 0 12 * * * /usr/bin/certbot renew --quiet
```

## 🔧 Manuel Deployment (Docker Olmadan)

### 1. Go Backend Kurulumu
```bash
# Go kurulumu
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Backend'i build edin
cd /var/www/blog-sitesi/backend
go build -o blog-backend .
```

### 2. Systemd Service Oluşturun
```bash
# Service dosyasını kopyalayın
sudo cp blog-backend.service /etc/systemd/system/

# Service'i etkinleştirin ve başlatın
sudo systemctl enable blog-backend
sudo systemctl start blog-backend
sudo systemctl status blog-backend
```

### 3. Nginx Kurulumu
```bash
# Nginx kurulumu
sudo apt update
sudo apt install nginx

# Konfigürasyonu kopyalayın
sudo cp nginx-blog.conf /etc/nginx/sites-available/sametegeasik.com
sudo ln -s /etc/nginx/sites-available/sametegeasik.com /etc/nginx/sites-enabled/

# Nginx'i yeniden başlatın
sudo systemctl reload nginx
```

## 🔒 Güvenlik Ayarları

### 1. Firewall Konfigürasyonu
```bash
# UFW kurulumu
sudo apt install ufw

# Gerekli portları açın
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

### 2. Fail2ban Kurulumu
```bash
# Fail2ban kurulumu
sudo apt install fail2ban

# Konfigürasyon
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

### 3. Güvenlik Başlıkları
Nginx konfigürasyonunda zaten mevcut:
- X-Frame-Options
- X-XSS-Protection
- X-Content-Type-Options
- Content-Security-Policy

## 📊 Monitoring ve Logging

### 1. Log Dosyaları
```bash
# Nginx logları
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# Backend logları
sudo journalctl -u blog-backend -f
```

### 2. Performance Monitoring
```bash
# Sistem kaynakları
htop
df -h
free -h

# Port durumu
netstat -tlnp
```

## 🔄 Backup Stratejisi

### 1. Database Backup
```bash
# Otomatik backup script
#!/bin/bash
BACKUP_DIR="/var/backups/blog"
DATE=$(date +%Y%m%d_%H%M%S)

# Database backup
cp /var/www/blog-sitesi/backend/blog.db $BACKUP_DIR/blog_$DATE.db

# Frontend backup
tar -czf $BACKUP_DIR/frontend_$DATE.tar.gz /var/www/blog-sitesi/*.html

# 30 günden eski backup'ları sil
find $BACKUP_DIR -name "*.db" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

### 2. Cron Job ile Otomatik Backup
```bash
# Crontab'a ekleyin
0 2 * * * /var/www/blog-sitesi/backup.sh
```

## 🚀 Deployment Scripti Kullanımı

### 1. Script'i Çalıştırın
```bash
# Deployment script'ini çalıştırın
./deploy.sh
```

### 2. Oluşturulan Dosyaları Kontrol Edin
- `blog-backend.service` - Systemd service dosyası
- `nginx-blog.conf` - Nginx konfigürasyonu

## 📈 Performance Optimizasyonu

### 1. Nginx Caching
```nginx
# Static dosyalar için cache
location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

### 2. Gzip Sıkıştırma
```nginx
# Nginx konfigürasyonuna ekleyin
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
```

### 3. Database Optimizasyonu
```sql
-- SQLite için VACUUM
VACUUM;
ANALYZE;
```

## 🔍 Troubleshooting

### Yaygın Sorunlar:

#### 1. Backend Başlamıyor
```bash
# Logları kontrol edin
sudo journalctl -u blog-backend -f

# Port kullanımını kontrol edin
sudo netstat -tlnp | grep 8081
```

#### 2. Nginx Hatası
```bash
# Konfigürasyon testi
sudo nginx -t

# Logları kontrol edin
sudo tail -f /var/log/nginx/error.log
```

#### 3. SSL Sertifikası Sorunu
```bash
# Sertifika durumunu kontrol edin
sudo certbot certificates

# Sertifikayı yenileyin
sudo certbot renew
```

## 📞 Destek

Sorun yaşarsanız:
1. Log dosyalarını kontrol edin
2. Sistem kaynaklarını kontrol edin
3. Network bağlantısını test edin
4. Konfigürasyon dosyalarını doğrulayın

---

**Not:** Bu rehber production ortamı için hazırlanmıştır. Test ortamında önce deneyin. 