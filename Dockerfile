FROM golang:1.17 as builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/aramase/azure-appconfig-csi-provider
ADD . .
RUN go build -o /bin/azure-appconfig-csi-provider

FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /bin/azure-appconfig-csi-provider .

ENTRYPOINT ["/azure-appconfig-csi-provider"]
