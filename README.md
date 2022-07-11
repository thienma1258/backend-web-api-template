
### BackEnd app

#### Prerequisites

- [golang](https://go.dev/)

## Installation

Run build app

```bash
$ make build
```

Run commands to run app

```bash
$ make run
```

## Testing

Running ping test

```bash
curl --location --request GET 'https://dongpham-challenge.herokuapp.com/ping' \
--header 'Content-Type: application/json' 
```