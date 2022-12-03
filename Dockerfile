FROM golang:1.19-alpine
WORKDIR /app
COPY ./ ./
RUN go build -o /service
EXPOSE 8080
CMD [ "/service" ]