FROM golang:latest 

ADD . /go/src/gophermarket/
ADD start.sh /usr/local/bin
COPY db /usr/local/bin/db
RUN chmod a+x /usr/local/bin/start.sh
RUN cd /go/src/gophermarket && rm go.mod && go mod init && go mod tidy
RUN cd /go/src/gophermarket/cmd/gophermart && go build -o /usr/local/bin/gophermarket
RUN cp /go/src/gophermarket/cmd/accrual/accrual_linux_amd64 /usr/local/bin/
RUN chmod a+x /usr/local/bin/accrual_linux_amd64
WORKDIR /usr/local/bin
EXPOSE 8080
EXPOSE 8000
CMD ["/usr/local/bin/start.sh"]