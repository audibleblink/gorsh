# Dockerfile for creating a container in which windows/machos
# binaries can be easily compiled if cgo is necessary
FROM dockercore/golang-cross:1.11.5
RUN go get github.com/gobuffalo/packr/...


