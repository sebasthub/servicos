FROM fedora

WORKDIR /home/

COPY ./ ./

RUN sudo dnf update -y

RUN sudo dnf install golang -y

RUN go get .

EXPOSE 8080

CMD [ "go", "run", "." ]
