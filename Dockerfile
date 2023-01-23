# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

RUN mkdir /build

ADD receipts /build/

WORKDIR /build
