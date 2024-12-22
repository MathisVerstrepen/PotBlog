![ReadMe Banner](https://raw.githubusercontent.com/MathisVerstrepen/github-visual-assets/refs/heads/main/banner/PotBlog.png)

## Prerequisites

- Go 1.21 or later
- templ CLI tool
- Docker (optional, for containerization)

## Getting Started

1. Clone this repository

2. Copy the `.env.example` file to `.env`:
   ```
   cp .env.example .env
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Install Air for live reload:
   ```
   go install github.com/air-verse/air@latest
   alias air="~/go/bin/air"
   ```

5. Install templ CLI tool:
   ```
   go install github.com/a-h/templ/cmd/templ@latest
   ```

6. Install tailwind dependencies:
   ```
   npm install
   ```

7. Run the development server:
   ```
   air
   ```

## Project Structure

```
.
├── assets
│   └── css
│   └── favicon
│   └── images
├── components
├── handlers
├── models
├── services
├── main.go
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
└── .air.toml
└── .env
```

## Development

- Use `air` to start the development server with live reload.
- Add new routes in `main.go`.
- Create handlers in the `/handlers` directory.
- Add templ components in the `/components` directory.

## Building for Production

1. Build the binary:
   ```
    go build -o bin/server
   ```

2. Run the production server:
   ```
    ./bin/server
   ```

## Docker Support

Build the Docker image:
```
docker compose build
```

Run the container:
```
docker compose up
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.