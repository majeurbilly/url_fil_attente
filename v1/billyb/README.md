# billyb

Billyb(ackend) est un backend ultra-simple en Go.

Le service contient une liste accessible au endpoint `/api/v1/items`. 

Par `GET`, il retourne le contenu de la liste
Par `POST`, il ajoute a la liste le contenu de la requete.

![logo](assets/logo.png)

## Stack

- [go-gin](https://github.com/gin-gonic/gin) 
- [otel](https://github.com/open-telemetry/opentelemetry-go)
- [prometheus](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
- [testify](https://github.com/stretchr/testify)
- [viper](https://github.com/spf13/viper)
- [zap](https://github.com/uber-go/zap)

## Usage

Dev-local (sans install)

```bash
go run cmd/api/main.go
```

Install (pour production)

```bash
go install cmd/api/main.go
GIN_MODE=release main
```

Envoyer un item dans la liste (sans le frontend)

```bash
curl -X POST http://localhost:8080/api/v1/items -H "Content-Type: application/json" -d '{"value": "<NOM DE MON ITEM>"}
```

Pour voir la liste des elements

```bash
curl http://localhost:8080/api/v1/items | jq
```


Pour voir les metriques

```bash
curl http://localhost:8080/metrics
```


Pour voir le healthcheck status

```bash
curl http://localhost:8080/health | jq
```