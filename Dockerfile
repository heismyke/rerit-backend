FROM node:22-alpine AS frontend
WORKDIR /src
COPY package*.json ./
RUN npm ci
COPY index.html tsconfig*.json vite.config.ts ./
COPY public ./public
COPY src ./src
RUN npm run build

FROM golang:1.23-alpine AS backend
WORKDIR /src/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/rerit-server ./cmd/rerit-server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/rerit-server-healthcheck ./cmd/rerit-server-healthcheck

FROM scratch

WORKDIR /app
COPY --from=backend /out/rerit-server /app/rerit-server
COPY --from=backend /out/rerit-server-healthcheck /app/rerit-server-healthcheck
COPY --from=frontend /src/dist /app/dist

ENV RERIT_ADDR=:8095
ENV RERIT_DATA_PATH=/data/rerit.json
ENV RERIT_STATIC_DIR=/app/dist

EXPOSE 8095
VOLUME ["/data"]

ENTRYPOINT ["/app/rerit-server"]
