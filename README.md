# Fixer

Fixer is a high-performance, scalable file search engine project that leverages the strengths of both Rust and Go to provide efficient search capabilities.

## 🎯 Project Goals

- Implementation of a high-performance file search engine using Rust
- Provision of a stable HTTP server using Go
- Flexible search functionality through fuzzy matching
- Scalability through microservices architecture

## 🏗 System Architecture

The project consists of two main components:

1. **Rust Search Engine** (`rust-search/`)
   - File system traversal and search index creation
   - Flexible search functionality using fuzzy-matcher
   - High-performance search processing

2. **Go HTTP Server** (`go-server/`)
   - RESTful API provision
   - Integration with Rust search engine
   - Client request handling

## 🚀 Getting Started

### Prerequisites
- Rust 1.70 or higher
- Go 1.20 or higher
- Docker (optional)

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/fixer.git
   cd fixer
   ```

2. Build the Rust component:
   ```bash
   cd rust-search
   cargo build --release
   ```

3. Start the Go server:
   ```bash
   cd ../go-server/cmd
   go mod tidy
   go run main.go
   ```

### Docker Deployment

```bash
docker-compose up
```

## 🔍 Key Features

- Recursive file system traversal
- Fuzzy matching for filename search
- Full-text search within files
- Search interface through RESTful API

## 🛠 API Endpoints

- `GET /search?q={query}` - Search by filename
- `GET /search/content?q={query}` - Search within file contents
- `GET /status` - Check service status

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
