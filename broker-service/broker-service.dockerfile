#build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

# coping over the executable from /app/brokerApp to /app we made above
COPY . /app

CMD ["/app/brokerApp" ]