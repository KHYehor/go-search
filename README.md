# 🔍 go-search

## 📓 Description

`go-search` is a pet project practice for self-development and improving skills. This is a high-performance text search engine written in Go, optimized for processing large files with efficient memory usage and parallel processing. 


---

## 📦 General Description

This project allows you to upload a file and search for multiple keywords giving as output array of coordinates `[line_number, column_number]`. 
It is designed to handle large-scale text processing using Go’s concurrency primitives, making it suitable for backend systems, developer tools, or any application where high-speed text search is required.

---

## 🏗️ Architecture

The system is built with the following key components:

- **Gin Web Framework** – lightweight HTTP server with middleware support.
- **Buffered File Scanner** – optimized file reading with large buffer size (1MB+ per line).
- **Sharded Concurrent Indexer** – search logic split across CPU cores using goroutines and channels.
- **JSON API** – accepts POST requests with file + keywords, returns search results.
- **Jobs Manager** – handles background task execution with request tracking.

---

## ⚙️ Features

---

## 📋 Results

---

## 📈 Future Improvements

---
