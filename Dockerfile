FROM golang:1.13 AS builder
RUN mkdir -p /src/github.com/ibrokethecloud/enforcer
COPY . /src/github.com/ibrokethecloud/enforcer
RUN cd /src/github.com/ibrokethecloud/enforcer \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o enforcer 

## Using upstream aquasec kube-bench and layering it up
FROM aquasec/trivy:0.13.0
COPY --from=builder /src/github.com/ibrokethecloud/enforcer/enforcer /usr/bin/enforcer
RUN mkdir /certs/
WORKDIR /
