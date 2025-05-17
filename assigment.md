# Assigment

## 📝 Instructions

Implement simple microservice (preferably in Go). The service will provide two REST API (accept and provide JSON) endpoints with following definition:

---

### 1. Storing data into DB

**Endpoint:**

```
POST /save
```

**Request body:**

```json
{
  "external_id": "<uuid>",
  "name": "some name",
  "email": "email@email.com",
  "date_of_birth": "2020-01-01T12:12:34+00:00"
}
```

---

### 2. Receiving data from DB

**Endpoint:**

```
GET /{id}
```

**Response:**

```json
{
  "external_id": "<uuid>",
  "name": "some name",
  "email": "email@email.com",
  "date_of_birth": "2020-01-01T12:12:34+00:00"
}
```

---

## ⚙️ Technologie

- You can use any relational DB of your choice (e.g. Postgres, SQLite, etc.).
- For data manipulation any ORM framework (we use GORM for instance).

---

## ✅ Requirements

- Please prepare a **buildable binary** which we can run and test (for example with curl commands).
- Prepare a **Dockerfile** building the image containing the binary.
  - If we run the container from the image we expect it exposes the HTTP server with endpoints above.
- **Bonus**: Create integration test testing both API endpoints.
  - The application is executed as the whole program as close to how it will run in production as possible.
  - Parsing ENVs, arguments, connecting to all dependencies etc.
  - Endpoints are tested via HTTP client from testing code.

---

