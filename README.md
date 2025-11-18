# 📄 Конфигурация проекта

Проект использует два источника настроек:

* **`.yaml`**
* **`.env`**

---

## 📁 Пример `config.yaml`

```yaml
environment: development

telegram:
  token: "your-telegram-token-here"
  poller_timeout: 10s

postgres:
  host: localhost
  port: 5432
  user: botik
  password: botik
  name: botik
  sslmode: disable

  max_conns: 25
  min_conns: 5
  max_conn_lifetime: 1h
  max_conn_idle_time: 30m
  connect_timeout: 10s
  query_timeout: 30s

app:
  shutdown_timeout: 30s
  request_timeout: 15s
```

---

## 📄 Пример `.env`

```env
ENVIRONMENT=development

BOT_TOKEN="your-telegram-token-here"

DB_HOST=localhost
DB_PORT=5432
DB_USER=botik
DB_PASSWORD=botik
DB_NAME=botik
DB_SSLMODE=disable
```

---