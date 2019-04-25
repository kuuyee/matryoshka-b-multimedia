FROM scratch
WORKDIR /app
COPY ./build/multimedia .
COPY ./multimedia.yml .

CMD ["./multimedia","-c","multimedia.yml"]