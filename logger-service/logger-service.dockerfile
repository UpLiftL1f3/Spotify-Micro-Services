#build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

# coping over the executable from /app/loggerServiceApp to /app we made above
COPY loggerServiceApp /app

# start it
CMD ["/app/loggerServiceApp" ]