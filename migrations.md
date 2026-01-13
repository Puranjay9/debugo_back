Below is a **clean, production-ready `README.md`** you can drop directly into your repo.
It documents the **migration workflow you’ve built**, step by step, without magic.

---

# Database Migrations (Go + Postgres + Docker)

This project uses **PostgreSQL migrations** powered by
[`golang-migrate`](https://github.com/golang-migrate/migrate) with **Docker**.

The goal is to:

* keep database schema changes **versioned**
* make migrations **safe and reproducible**
* separate **application code** from **schema management**

---

## 🧠 Migration Philosophy

* Migrations are **SQL files**, tracked in Git
* Schema changes are **explicit**, not automatic
* Go application **does NOT manage migrations**
* Docker is used to apply migrations consistently

This is the **industry-standard approach** for Go backends.

---

## 📁 Project Structure

```text
.
├── migrations/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   └── ...
├── scripts/
│   └── migrate.sh
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

---

## 🛠️ Required Tools

### 1️⃣ On your local machine (WSL / Linux)

You must install the **migrate CLI** to create migration files:

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64 \
  -o /usr/local/bin/migrate
chmod +x /usr/local/bin/migrate
```

Verify installation:

```bash
migrate -version
```

---

### 2️⃣ Docker (no manual install)

Migrations are **applied using Docker** via the official image:

```yaml
image: migrate/migrate
```

Docker will pull this automatically.

---

## 🐘 Database Configuration

Postgres runs in Docker using this service:

```yaml
go_db:
  image: postgres:12
  environment:
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: postgres
  ports:
    - "5432:5432"
```

Connection string used by migrations:

```text
postgres://postgres:postgres@go_db:5432/postgres?sslmode=disable
```

> ⚠️ `go_db` works **only inside Docker**
> Use `localhost` if connecting from host

---

## 🧩 Migration Workflow (Step by Step)

### 1️⃣ Create a new migration

Run on **your host machine**:

```bash
./scripts/migrate.sh add_users_table
```

This creates two files:

```text
migrations/
├── 000002_add_users_table.up.sql
└── 000002_add_users_table.down.sql
```

---

### 2️⃣ Write migration SQL

#### `000002_add_users_table.up.sql`

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
```

#### `000002_add_users_table.down.sql`

```sql
DROP TABLE users;
```

---

### 3️⃣ Apply migrations

When prompted, press **ENTER**
or manually run:

```bash
docker compose run migrate
```

This will:

* read migration files
* apply unapplied versions
* record progress in `schema_migrations`

---

## 🔁 Rollback Migrations

To roll back the last migration:

```bash
docker compose run migrate down 1
```

To roll back everything:

```bash
docker compose run migrate down
```

---

## 🧾 How Migration State Is Tracked

`golang-migrate` automatically creates:

```sql
schema_migrations
```

This table stores:

* migration version
* dirty state (failed migrations)

⚠️ Never edit this table manually.

---

## 🚫 What This Project Does NOT Do

* ❌ No automatic schema detection
* ❌ No ORM auto-migration
* ❌ No migrations inside Go code
* ❌ No runtime schema changes

This is intentional and **safer for production**.

---

## ✅ Best Practices

✔ One logical change per migration
✔ Always write a `down.sql`
✔ Never edit applied migrations
✔ Commit migrations to Git
✔ Run migrations before app start

---

## 🧪 Common Commands

```bash
# Create migration
./scripts/migrate.sh add_feature_x

# Apply migrations
docker compose run migrate

# Roll back last migration
docker compose run migrate down 1
```

---

## 📌 Summary

* `migrate create` → runs on **host**
* `migrate up/down` → runs in **Docker**
* SQL files are the **source of truth**
* Go app stays **schema-agnostic**


