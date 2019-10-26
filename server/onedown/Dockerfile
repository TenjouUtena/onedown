FROM golang
WORKDIR /app
ADD onedown /app
ADD configuration.json /app/configuration/
CMD ["./onedown", "--config=/app/configuration/"]
