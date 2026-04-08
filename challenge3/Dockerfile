# pakai image Go
FROM golang:1.26.0-alpine

# set working dir
WORKDIR /app

# copy go mod dulu (biar cache)
# COPY go.mod go.sum ./
COPY . .
RUN go mod tidy

# copy semua file
COPY . .

# build binary
RUN go build -o main ./cmd

# expose port
EXPOSE 8000

# run app
CMD ["./main"]