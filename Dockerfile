FROM gcr.io/distroless/static

COPY ./_output/azure-appconfig-csi-provider .

LABEL maintainers="aramase"
LABEL description="Secrets Store CSI Driver Provider Azure AppConfig"

ENTRYPOINT [ "/azure-appconfig-csi-provider" ]
