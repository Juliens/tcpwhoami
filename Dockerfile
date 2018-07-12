FROM golang AS build
RUN  go get -u golang.org/x/vgo
COPY .  /go/src/github.com/juliens/tcpwhoami

WORKDIR /go/src/github.com/juliens/tcpwhoami

#RUN go generate
RUN CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix cgo -o main .


FROM scratch

EXPOSE 8080
COPY --from=build  /go/src/github.com/juliens/tcpwhoami/main /
ENTRYPOINT [ "/main" ]
