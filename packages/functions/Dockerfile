# https://docs.docker.com/language/golang/build-images/
FROM golang:1.22-bullseye AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN env GOOS=linux GOARCH=amd64 \
    go build -o /main ./endpoints/search.go

# Install chromium
FROM public.ecr.aws/lambda/provided:al2 as chromium

# install brotli, so we can decompress chromium
RUN yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm && \
    yum -y install brotli && \
    yum clean all

# download chromium
RUN yum -y install wget && \
    wget --progress=dot:giga https://raw.githubusercontent.com/alixaxel/chrome-aws-lambda/master/bin/chromium.br -O /chromium.br && \
    yum clean all

# decompress chromium
RUN brotli -d /chromium.br


# copy artifacts to a clean image
FROM  public.ecr.aws/lambda/provided:al2

# install chromium dependencies
RUN yum -y install \
    libX11 \
    nano \
    unzip \
    wget \
    xclock \
    xorg-x11-xauth \
    xterm && \
    yum clean all

# copy in chromium from chromium stage
COPY --from=chromium /chromium /opt/chromium

# grant our program access to chromium
RUN chmod 777 /opt/chromium

# copy in lambda fn from build stage
COPY --from=build /main /main

ENTRYPOINT ["/main"]