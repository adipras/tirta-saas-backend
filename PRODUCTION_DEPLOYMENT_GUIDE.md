# 🚀 Tirta-SaaS Production Deployment Guide

## 📋 Overview

This guide provides step-by-step instructions for deploying the Tirta-SaaS water utility billing backend to production environments. The system is enterprise-ready with comprehensive security, monitoring, and performance optimizations.

## ✅ Pre-Deployment Checklist

### System Requirements
- [ ] **Go Version**: 1.24.2 or higher
- [ ] **MySQL**: 8.0+ (recommended) or 5.7+
- [ ] **Memory**: Minimum 2GB RAM (4GB recommended)
- [ ] **CPU**: 2+ cores recommended
- [ ] **Storage**: 20GB+ available disk space
- [ ] **Network**: HTTPS/TLS capabilities

### Security Prerequisites
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules configured (only necessary ports open)
- [ ] Database access restricted to application servers
- [ ] Strong passwords/keys for all services
- [ ] Backup and recovery procedures established

---

## 🔧 Environment Configuration

### 1. Environment Variables

Create a `.env` file with the following required variables:

```bash
# Database Configuration
DB_HOST=your-database-host
DB_PORT=3306
DB_USER=your-db-username
DB_PASS=your-strong-db-password
DB_NAME=tirta_saas_production

# Application Configuration
PORT=8080
GIN_MODE=release
APP_VERSION=1.0.0
ENV=production

# Security Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-min-64-chars-recommended

# Optional: Advanced Configuration
MAX_DB_CONNECTIONS=100
DB_MAX_IDLE=10
DB_MAX_LIFETIME=300s
LOG_LEVEL=INFO
```

### 2. Database Setup

#### MySQL Configuration (recommended settings)
```sql
-- Create database and user
CREATE DATABASE tirta_saas_production CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'tirta_user'@'%' IDENTIFIED BY 'your-strong-password';
GRANT ALL PRIVILEGES ON tirta_saas_production.* TO 'tirta_user'@'%';
FLUSH PRIVILEGES;

-- Recommended MySQL settings for production
SET GLOBAL innodb_buffer_pool_size = 1073741824; -- 1GB
SET GLOBAL max_connections = 200;
SET GLOBAL innodb_log_file_size = 268435456; -- 256MB
```

#### Database Optimization
```sql
-- Add recommended indexes (automatically created by migration)
-- The application includes 35+ optimized indexes for performance

-- Monitor slow queries
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;
SET GLOBAL log_queries_not_using_indexes = 'ON';
```

---

## 🐳 Deployment Options

### Option 1: Docker Deployment (Recommended)

#### 1. Create Dockerfile
```dockerfile
# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tirta-saas-backend .

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/tirta-saas-backend .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./tirta-saas-backend"]
```

#### 2. Create docker-compose.yml
```yaml
version: '3.8'

services:
  tirta-saas-backend:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=tirta_user
      - DB_PASS=secure_password
      - DB_NAME=tirta_saas_production
      - JWT_SECRET=your-super-secure-jwt-secret-key
      - GIN_MODE=release
      - PORT=8080
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=tirta_saas_production
      - MYSQL_USER=tirta_user
      - MYSQL_PASSWORD=secure_password
    volumes:
      - mysql_data:/var/lib/mysql
      - ./mysql-config:/etc/mysql/conf.d
    restart: unless-stopped
    command: --default-authentication-plugin=mysql_native_password

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - tirta-saas-backend
    restart: unless-stopped

volumes:
  mysql_data:
```

#### 3. Deploy with Docker
```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f tirta-saas-backend

# Scale application (if needed)
docker-compose up --build -d --scale tirta-saas-backend=3
```

### Option 2: Direct Server Deployment

#### 1. Build Application
```bash
# On development machine
go mod tidy
go build -o tirta-saas-backend

# Transfer to production server
scp tirta-saas-backend user@production-server:/opt/tirta-saas/
scp .env user@production-server:/opt/tirta-saas/
```

#### 2. Create systemd Service
```ini
# /etc/systemd/system/tirta-saas.service
[Unit]
Description=Tirta-SaaS Backend Service
After=network.target mysql.service
Requires=mysql.service

[Service]
Type=simple
User=tirta-saas
Group=tirta-saas
WorkingDirectory=/opt/tirta-saas
ExecStart=/opt/tirta-saas/tirta-saas-backend
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/tirta-saas

[Install]
WantedBy=multi-user.target
```

#### 3. Enable and Start Service
```bash
# Create user and directory
sudo useradd -r -s /bin/false tirta-saas
sudo mkdir -p /opt/tirta-saas
sudo chown tirta-saas:tirta-saas /opt/tirta-saas

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable tirta-saas
sudo systemctl start tirta-saas

# Check status
sudo systemctl status tirta-saas
sudo journalctl -u tirta-saas -f
```

---

## 🔒 Security Configuration

### 1. Reverse Proxy (Nginx)
```nginx
# /etc/nginx/sites-available/tirta-saas
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL Configuration
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;

    # Security Headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains";

    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=auth:10m rate=5r/s;

    location / {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /auth/ {
        limit_req zone=auth burst=10 nodelay;
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. Firewall Configuration
```bash
# UFW (Ubuntu)
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw deny 8080  # Block direct access to app
sudo ufw enable

# iptables alternative
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP
```

---

## 📊 Monitoring & Logging

### 1. Health Check Endpoints
```bash
# Health check (comprehensive)
curl https://your-domain.com/health

# Readiness check (Kubernetes/Docker)
curl https://your-domain.com/ready

# Liveness check (simple)
curl https://your-domain.com/alive

# Detailed metrics
curl https://your-domain.com/metrics
```

### 2. Log Management
```bash
# View application logs
sudo journalctl -u tirta-saas -f --since "1 hour ago"

# Log rotation (logrotate)
# /etc/logrotate.d/tirta-saas
/opt/tirta-saas/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

### 3. Performance Monitoring
```bash
# Monitor resource usage
htop
iotop
nethogs

# Database monitoring
mysql -u root -p -e "SHOW PROCESSLIST;"
mysql -u root -p -e "SHOW STATUS LIKE 'Threads%';"
```

---

## 🔄 Backup & Recovery

### 1. Database Backup
```bash
#!/bin/bash
# backup-database.sh
DATE=$(date +"%Y%m%d_%H%M%S")
BACKUP_DIR="/opt/backups/tirta-saas"
DB_NAME="tirta_saas_production"

mkdir -p $BACKUP_DIR

# Create database backup
mysqldump -u tirta_user -p$DB_PASS $DB_NAME > $BACKUP_DIR/db_backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/db_backup_$DATE.sql

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: db_backup_$DATE.sql.gz"
```

### 2. Automated Backups (Cron)
```bash
# Add to crontab (crontab -e)
0 2 * * * /opt/tirta-saas/scripts/backup-database.sh >> /var/log/tirta-backup.log 2>&1
```

---

## 🚀 Deployment Workflow

### 1. Zero-Downtime Deployment
```bash
#!/bin/bash
# deploy.sh
set -e

echo "Starting deployment..."

# Build new version
go build -o tirta-saas-backend-new

# Test new binary
./tirta-saas-backend-new --version

# Stop old service
sudo systemctl stop tirta-saas

# Replace binary
mv tirta-saas-backend-new tirta-saas-backend

# Start new service
sudo systemctl start tirta-saas

# Verify deployment
sleep 5
curl -f http://localhost:8080/health || exit 1

echo "Deployment successful!"
```

### 2. Blue-Green Deployment (Docker)
```bash
# Deploy to blue environment
docker-compose -f docker-compose.blue.yml up -d

# Run health checks
curl -f http://blue.your-domain.com/health

# Switch traffic (update load balancer)
# Stop green environment
docker-compose -f docker-compose.green.yml down
```

---

## 📈 Performance Optimization

### 1. Application Tuning
```bash
# Environment variables for production
export GOGC=100
export GOMAXPROCS=4
export GOMEMLIMIT=2GiB
```

### 2. Database Tuning
```sql
-- my.cnf optimizations
[mysqld]
innodb_buffer_pool_size = 2G
innodb_log_file_size = 512M
innodb_flush_log_at_trx_commit = 2
query_cache_size = 256M
max_connections = 200
```

### 3. Connection Pooling
```go
// Configured automatically by the application
// Default settings:
// - Max Open Connections: 100
// - Max Idle Connections: 10
// - Connection Max Lifetime: 5 minutes
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. Database Connection Issues
```bash
# Check database connectivity
mysql -h $DB_HOST -u $DB_USER -p$DB_PASS -e "SELECT 1;"

# Check connection pool stats
curl http://localhost:8080/metrics | jq '.database'
```

#### 2. Memory Issues
```bash
# Check memory usage
free -h
ps aux --sort=-%mem | head

# Check application metrics
curl http://localhost:8080/metrics | jq '.memory'
```

#### 3. Performance Issues
```bash
# Check slow queries
sudo tail -f /var/log/mysql/slow.log

# Monitor request performance
curl http://localhost:8080/metrics | jq '.requests'
```

### Health Check Responses

#### Healthy System
```json
{
  "status": "healthy",
  "timestamp": "2024-01-31T10:00:00Z",
  "checks": {
    "database": {"status": "healthy", "duration": "5ms"},
    "memory": {"status": "healthy"},
    "disk": {"status": "healthy"}
  }
}
```

#### Unhealthy System
```json
{
  "status": "unhealthy",
  "timestamp": "2024-01-31T10:00:00Z",
  "checks": {
    "database": {"status": "unhealthy", "message": "Connection timeout"}
  }
}
```

---

## 🎯 Production Checklist

### Pre-Launch
- [ ] All environment variables configured
- [ ] Database optimized and secured
- [ ] SSL certificates installed
- [ ] Reverse proxy configured
- [ ] Firewall rules applied
- [ ] Health checks passing
- [ ] Backup system tested
- [ ] Monitoring alerts configured
- [ ] Load testing completed
- [ ] Security audit passed

### Post-Launch
- [ ] Monitor error rates
- [ ] Check response times
- [ ] Verify backup schedules
- [ ] Test failover procedures
- [ ] Document runbook procedures
- [ ] Train operations team
- [ ] Set up alerting thresholds
- [ ] Review security logs

---

## 📞 Support & Maintenance

### Monitoring Alerts
Set up alerts for:
- HTTP error rate > 5%
- Response time > 1 second
- Database connections > 80
- Memory usage > 80%
- Disk space < 20%
- Health check failures

### Regular Maintenance
- **Weekly**: Review performance metrics and logs
- **Monthly**: Database optimization and cleanup
- **Quarterly**: Security updates and patches
- **Annually**: Complete security audit

---

## 🎉 Conclusion

The Tirta-SaaS backend is now production-ready with enterprise-grade features:

- ✅ **Security**: HTTPS, rate limiting, input validation
- ✅ **Monitoring**: Health checks, metrics, logging
- ✅ **Performance**: Optimized queries, connection pooling
- ✅ **Reliability**: Error handling, graceful degradation
- ✅ **Scalability**: Horizontal scaling support

For additional support or questions, refer to the API documentation at `/swagger/index.html` or contact the development team.

**🚀 Ready for Production Deployment!**