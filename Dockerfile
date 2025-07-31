# ==============================
# 👉 المرحلة 1: بناء الواجهة الأمامية (Frontend)
# ==============================
FROM node:20 AS frontend-builder

WORKDIR /build

COPY web/package.json web/yarn.lock ./
RUN yarn --frozen-lockfile

COPY ./web .
COPY ./VERSION .
RUN DISABLE_ESLINT_PLUGIN='true' VITE_APP_VERSION=$(cat VERSION) npm run build

# ==============================
# 👉 المرحلة 2: بناء الباك اند Go
# ==============================
FROM golang:1.24.2 AS backend-builder

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://proxy.golang.org,direct

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /build/build ./web/build
RUN go build -ldflags "-s -w -X 'done-hub/common.Version=$(cat VERSION)' -extldflags '-static'" -o done-hub

# ==============================
# 👉 المرحلة 3: الصورة النهائية
# ==============================
FROM alpine:latest

# إعداد البيئة
ENV TZ=Asia/Shanghai \
    USER_TOKEN_SECRET=3378844ccf5baae72b6135ea350bd0 \
    SESSION_SECRET=82f209eece5181530e7947fa205b6c10 \
    SQL_DSN=postgres://root:312sCEUn4flh7FeQ06cBZMkaG5RyJ9p8@sjc1.clusters.zeabur.com:32092/zeabur

RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates 2>/dev/null || true

# نسخ البرنامج النهائي
COPY --from=backend-builder /build/done-hub /done-hub

EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/done-hub"]
