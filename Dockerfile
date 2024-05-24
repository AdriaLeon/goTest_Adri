FROM golang:1.22.3

WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Install necessary tools for Python
RUN apt-get update && apt-get install -y python3-venv

# Create and activate virtual environment
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install Python dependencies
RUN pip install --upgrade pip
RUN pip install -r lib/requirements.txt

# Build the Go application
RUN go build -o bin/main cmd/main.go

EXPOSE 8080

ENTRYPOINT ["./bin/main"]
