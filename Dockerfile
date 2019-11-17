# Onedown Dockerfile. Dockerize me baby!
FROM golang:1.13.4-buster

# Copy application files
ADD . src/github.com/TenjouUtena/onedown

# Dependency stuff for backend
RUN go get github.com/golang/dep/cmd/dep

# Backend setup
WORKDIR src/github.com/TenjouUtena/onedown/backend
RUN dep ensure
RUN go install github.com/TenjouUtena/onedown/backend

# Frontend setup
WORKDIR src/github.com/TenjouUtena/onedown/frontend/onedown
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt install nodejs
RUN npm install

EXPOSE 3000
