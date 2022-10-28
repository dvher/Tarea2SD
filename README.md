# Tarea 2 Kafka

### 1. Construccion del entorno

```
docker compose up -d
```

### 2. Poner a funcionar servidor

```
go run cmd/server/server.go
```

### 3. Correr cliente
Para probar las funcionalidades de la app.

```
go run cmd/client/client.go
```
### 4. Para ver los resultados


Procesamiento de Coordenadas

```
go run cmd/coordinates/coordinates.go
```

Procesamiento de Ventas

```
go run cmd/ventas/ventas.go
```

Procesamiento de Stock

```
go run cmd/stock/stock.go
```

Procesamiento de Miembros

```
go run cmd/member/member.go
```
