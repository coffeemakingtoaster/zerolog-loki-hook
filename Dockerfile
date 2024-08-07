FROM golang:1.22.1-alpine3.18 AS build 

WORKDIR /build

COPY . .

RUN go build -o zerolog-loki 

## -- RUNTIME STAGE --
FROM alpine:3.18.4 AS runtime

WORKDIR /app

ARG USER=docker
ARG UID=5432
ARG GID=5433

# Create user for execution

#User group has same name as user
RUN addgroup -g $GID $USER 

RUN adduser \
	--disabled-password \
	--gecos "" \
	--ingroup "$USER" \
	--no-create-home \
	--uid "$UID" \
	"$USER"

# Copy build with permissions
COPY --from=build --chown=$USER:$USER /build/zerolog-loki /app/zerolog-loki

# Ensure that backend can be run
RUN chmod +x /app/zerolog-loki

USER $USER 

CMD ["/app/zerolog-loki"]
