# golang_practice

A Golang (Gin) practice repository focused on **middleware design**, **request/response logging**, and **safe handling of different content types**.

This project is created **for learning and testing purposes**.

---

## 📌 Purpose

- Practice Gin framework fundamentals
- Implement custom middleware
- Log request & response activities
- Handle multipart file uploads safely
- Follow real-world backend best practices

---

## 🛠 Tech Stack

- Go (Golang)
- Gin Web Framework
- HTTP / REST APIs
- Middleware & Logging

---

## 🧠 Golang Areas Covered

This repository focuses on the following **core Golang areas**:

### 🔹 Language Fundamentals
- Structs & methods
- Interfaces
- Packages & modules
- Error handling
- Pointers & memory basics

### 🔹 Concurrency
- Goroutines
- Asynchronous operations
- Non-blocking middleware logic

### 🔹 HTTP & Web
- net/http basics
- Gin framework
- Request lifecycle
- Middleware chaining
- Context handling (`gin.Context`)

### 🔹 Middleware Design
- Request interception
- Response interception
- Custom `ResponseWriter`
- Header & payload inspection
- Latency measurement

### 🔹 Data Handling
- JSON encoding/decoding
- Safe request body reading
- Multipart form-data handling
- Metadata extraction

### 🔹 Performance & Safety
- Avoiding memory leaks
- Handling large payloads
- Skipping binary/file logging
- Asynchronous DB writes

---

## 🚫 Ignoring Logs for File Upload Requests

When logging request/response data in middleware or logger,  
**multipart/form-data requests are intentionally ignored**.

### ❓ Why?

- File uploads may contain large binary data
- Reading request body loads entire file into memory
- Can cause high RAM usage or crashes
- Logging file content is unsafe and unnecessary
- Gin already streams files efficiently

---

## ✅ Safe Multipart Detection

```go
contentType := c.GetHeader("Content-Type")

if strings.Contains(contentType, "multipart/form-data") {
    // Handle file upload safely
    // Do NOT read request body
}
## set env configuration based on system 
// like for local linux ubuntu just run "export APP_ENV=dev"