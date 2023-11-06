FROM golang:latest

WORKDIR /app
EXPOSE 1323

ENV PATH="${PATH}:/go/bin"
RUN echo $PATH && ls -l /go/bin

RUN set -x && go install github.com/cosmtrek/air@latest
RUN ls -l /go/bin > /tmp/ls_result.txt && cat /tmp/ls_result.txt
RUN find / -name air > /tmp/find_result.txt && cat /tmp/find_result.txt


COPY go.mod go.sum ./

RUN go mod download

RUN ls /go/bin
COPY . .
RUN ls /go/bin

EXPOSE 1323

CMD ["air"]