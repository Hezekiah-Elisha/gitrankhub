**Gitrankhub**

Simple Go (Gin) API service. This README covers setup, environment, running locally, and common project workflows.

**Project Structure**
- go.mod: Module and dependencies
- main.go: App entrypoint
- config/config.go: App configuration loading
- handlers/usersHandler.go: User-related HTTP handlers
- models/user.go: Data models

**Prerequisites**
- Go 1.20+ installed
- Git installed

**Environment**
- Create a `.env` file in the project root with values like:
	```dotenv
	DB_USERNAME=your_user
	DB_PASSWORD=your_password
	DB_DATABASE=gitrankhub_db
	DB_HOST=localhost
	DB_PORT=3306
	```
- Ensure `.env` is ignored by Git (already in `.gitignore`). If it was ever committed before ignoring, remove it from tracking with:
	```bash
	git rm --cached .env
	```

**Install Dependencies**
- From the project root, install required dependencies:
	```bash
	go get .
	```

**Run Locally**
- Start the server:
	```bash
	go run .
	```

**Build**
- Produce a binary (Windows shown):
	```bash
	go build -o bin/gitrankhub.exe .
	```

**Common Commands**
- Tidy modules: `go mod tidy`
- Run tests: `go test ./...`
- Lint/format: `gofmt -s -w .`

**Notes**
- Dependencies were installed using `go get .`.
- The project is run locally using `go run .`.
- `.env` was removed from Git tracking using `git rm --cached .env`.
