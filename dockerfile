# Build Stage
FROM golang:1.21.5 AS build
WORKDIR /app
COPY . .
RUN go build -o main .

# Final Stage
FROM golang:1.21.5
WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8000
CMD ["./main"]
