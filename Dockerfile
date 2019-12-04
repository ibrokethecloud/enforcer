FROM golang:1.12 AS builder
RUN mkdir -p /src/github.com/ibrokethecloud/enforcer
COPY . /src/github.com/ibrokethecloud/enforcer
RUN cd /src/github.com/ibrokethecloud/enforcer \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o enforcer -mod vendor 

## Using upstream aquasec kube-bench and layering it up
FROM aquasec/trivy:latest
COPY --from=builder /src/github.com/ibrokethecloud/enforcer/enforcer /usr/bin/enforcer
RUN mkdir /certs/
COPY docker-entrypoint.sh /
WORKDIR /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD [ "/usr/bin/enforcer","-severity","CRITICAL","-tls-cert-file","/certs/webhook.crt","-tls-key-file","/certs/webhookCA.key","-prune-images","true"  ]
