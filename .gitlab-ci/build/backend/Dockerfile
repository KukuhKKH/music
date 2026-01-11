FROM ${IMAGE_NAMESPACE}/base-music-backend:latest

WORKDIR /app

COPY backend/cmd/web/main .
COPY backend/config/config.toml .

RUN chmod +x main

EXPOSE 3000
CMD ["./main"]
