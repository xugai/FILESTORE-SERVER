#务必注意，编译环境是在MAC，运行环境是基于linux，两者环境不同，因此你在编译期间要进行交叉编译，也就是指定操作系统与CPU指令架构的类别amd64/386
# GOOS=linux GOARCH=amd64 go build -o dest src

FROM golang

ADD account /

RUN chmod +x /account

ENTRYPOINT ["sh", "-c", "/account"]

