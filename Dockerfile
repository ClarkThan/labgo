FROM --platform=$TARGETPLATFORM golang:alpine as builder
COPY . /opt/mybuild/
WORKDIR /opt/mybuild
RUN go build -o labgo

FROM --platform=$TARGETPLATFORM alpine
WORKDIR /app
COPY --from=builder  /opt/mybuild/labgo /app/labgo
CMD ./labgo

# https://github.com/lmk123/docker-buildx-test
# docker build -t labgo:aarch64 .
# https://www.cnblogs.com/lisongyu/p/16212037.html
# docker build --platform linux/amd64 -t labgo:x86_64 .