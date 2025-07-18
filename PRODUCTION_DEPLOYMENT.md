# Production Deployment Guide

## 🚀 Production'a Çıkma Rehberi

### 1. Gereksinimler
- Go 1.21 veya üstü
- SSL sertifikası (Let's Encrypt önerilen)
- Domain adı
- Linux server (Ubuntu/CentOS)

### 2. Hızlı Başlangıç

```bash
# 1. Proje dosyalarını sunucuya kopyalayın
scp -r blog-sitesi/ user@your-server:/var/www/

# 2. Sunucuya bağlanın
ssh user@your-server

# 3. Proje dizinine gidin
cd /var/www/blog-sitesi/backend

# 4. Start script'i çalıştırın
./start.sh
```

### 3. Manuel Kurulum

#### Environment Variables Ayarlama

```bash
# .env dosyası oluşturun
cat > .env << EOF
JWT_SECRET=$(openssl rand -base64 32)
GIN_MODE=release
DATABASE_URL=./blog.db
PORT=8081
PRODUCTION_DOMAIN=https://yourdomain.com
BASE_DIR=/var/www/blog-sitesi
EOF
```

#### Build ve Çalıştırma

```bash
# Build
go build -o blog-backend .

# Çalıştırma
./blog-backend
```

### 4. Systemd Service Kurulumu

```bash
# Service dosyası oluşturun
sudo tee /etc/systemd/system/blog-backend.service > /dev/null <<EOF
[Unit]
Description=Blog Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/blog-sitesi/backend
ExecStart=/var/www/blog-sitesi/backend/blog-backend
Restart=always
RestartSec=10
Environment=JWT_SECRET=$(openssl rand -base64 32)
Environment=GIN_MODE=release
Environment=DATABASE_URL=./blog.db
Environment=PORT=8081
Environment=PRODUCTION_DOMAIN=https://yourdomain.com
Environment=BASE_DIR=/var/www/blog-sitesi

[Install]
WantedBy=multi-user.target
EOF

# Service'i etkinleştirin
sudo systemctl daemon-reload
sudo systemctl enable blog-backend.service
sudo systemctl start blog-backend.service

# Durumu kontrol edin
sudo systemctl status blog-backend.service
```

### 5. Nginx Reverse Proxy

```nginx
# /etc/nginx/sites-available/blog-backend
server {
    listen 80;
    server_name yourdomain.com www.yourdomain.com;
    
    # SSL için Let's Encrypt
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com www.yourdomain.com;
    
    # SSL sertifikaları
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    
    # Güvenlik başlıkları
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    # Gzip sıkıştırma
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
    
    # Backend proxy
    location /api/ {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # Static files
    location / {
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 6. SSL Sertifikası (Let's Encrypt)

```bash
# Certbot yükleyin
sudo apt update
sudo apt install certbot python3-certbot-nginx

# Sertifika alın
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Otomatik yenileme
sudo crontab -e
# Bu satırı ekleyin:
0 12 * * * /usr/bin/certbot renew --quiet
```

### 7. Firewall Ayarları

```bash
# UFW firewall
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Port 8081'i sadece localhost'a açın
sudo ufw allow from 127.0.0.1 to any port 8081
```

### 8. Monitoring ve Logs

```bash
# Service loglarını görüntüleyin
sudo journalctl -u blog-backend.service -f

# Nginx logları
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 9. Backup

```bash
# Veritabanı yedekleme scripti
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
cp /var/www/blog-sitesi/backend/blog.db /backup/blog_$DATE.db

# Crontab ekleyin
0 2 * * * /path/to/backup-script.sh
```

### 10. Güvenlik Kontrolleri

- [x] JWT_SECRET güvenli şekilde ayarlandı
- [x] Rate limiting aktif
- [x] CORS kısıtlı
- [x] Security headers ayarlandı
- [x] Input validation aktif
- [x] Error handling güvenli
- [x] HTTPS zorlaması
- [x] Firewall konfigürasyonu

### 11. Environment Variables Tablosu

| Variable | Açıklama | Varsayılan |
|----------|----------|------------|
| JWT_SECRET | JWT token için güvenli anahtar | random |
| GIN_MODE | Gin framework modu | release |
| DATABASE_URL | SQLite veritabanı yolu | ./blog.db |
| PORT | Server port | 8081 |
| PRODUCTION_DOMAIN | Domain adı | - |
| BASE_DIR | Statik dosyalar dizini | /var/www/blog-sitesi |

### 12. Troubleshooting

```bash
# Service durumu
sudo systemctl status blog-backend.service

# Logları kontrol edin
sudo journalctl -u blog-backend.service -n 50

# Port kontrolü
sudo netstat -tlpn | grep 8081

# Process kontrolü
ps aux | grep blog-backend
```

### 13. Güncelleme

```bash
# Yeni versiyon deploy
cd /var/www/blog-sitesi/backend
git pull origin main
go build -o blog-backend .
sudo systemctl restart blog-backend.service
```

---

## 🎯 Özet

1. `./start.sh` ile hızlı başlangıç yapabilirsiniz
2. JWT_SECRET otomatik olarak güvenli şekilde oluşturulur
3. Tüm güvenlik önlemleri aktif
4. Production'da HTTPS zorunlu
5. Nginx reverse proxy önerilen
6. Systemd ile otomatik başlatma
7. Let's Encrypt ile ücretsiz SSL

**Önemli:** Production'da mutlaka HTTPS kullanın ve JWT_SECRET'ı güvenli tutun! 