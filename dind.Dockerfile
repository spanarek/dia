FROM hashicorp/vault:1.12.2 AS vault
FROM hashicorp/consul-template:0.30.0 AS consul

MAINTAINER Sergey Panarek sergey@panarek.ru
FROM docker:20.10.23-dind-alpine3.17

COPY --from=vault /bin/vault /usr/local/bin/vault
COPY --from=consul /bin/consul-template /usr/local/bin/consul-template
COPY docker/dind/cert.ctmpl /usr/local/share/
COPY docker/dind/dia-sign.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/dia-sign.sh
RUN apk add grep bash
