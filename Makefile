.PHONY: build
build:
	go build -o ai-course cmd/ai-course/main.go

.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -o ai-course-linux cmd/ai-course/main.go
	GOOS=windows GOARCH=amd64 go build -o ai-course.exe cmd/ai-course/main.go
	GOOS=darwin GOARCH=amd64 go build -o ai-course-mac cmd/ai-course/main.go