FROM node:10-alpine
RUN mkdir -p /hyperledger/caliper/workspace && npm -g config set user root && npm i -g @hyperledger/caliper-cli@0.3.2 node-gyp
RUN apk add python3 gcc g++ make libffi-dev openssl-dev libtool
RUN caliper bind --caliper-bind-sut fabric:1.4.4 --caliper-bind-args=-g