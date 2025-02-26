# RustySearch Go Server

This project is a Go server that includes functionality for searching and returning results from files. It is designed to be efficient and easy to use.

## Project Structure

```
RustySearch
└── go-server
    ├── .github
    │   └── workflows
    │       └── go-test.yml
    ├── domain
    │   └── model
    │       └── search_test.go
    └── README.md
```

## Getting Started

To get started with the project, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/RustySearch.git
   cd RustySearch/go-server
   ```

2. **Install Go:**
   Make sure you have Go installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

3. **Install dependencies:**
   Run the following command to install the necessary dependencies:
   ```bash
   go mod tidy
   ```

4. **Run tests:**
   To run the tests, use the following command:
   ```bash
   go test ./... -v
   ```

## GitHub Actions

This project uses GitHub Actions for continuous integration. The workflow is defined in `.github/workflows/go-test.yml`, which runs tests automatically on pushes and pull requests to the main branch.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.