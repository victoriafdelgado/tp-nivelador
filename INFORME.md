# Protocolo de Comunicación
Se utiliza un protocolo con patrón TLV (Type, Length, Value).
Cada paquete está compuesto por dos partes: un header (tamaño fijo) y un payload (tamaño variable).

En el header:
- Código de tipo de mensaje (1 byte): Identifica qué representa el mensaje.
- Largo del mensaje (4 bytes): Es el largo del payload en bytes, en orden BigEndian.

El payload:
- Payload (tamaño variable): Es el mensaje a enviar en bytes.

| Tipo | Largo | Mensaje |
|---|---|---|
| 1 byte | 4 bytes | N bytes (según longitud) |

# Implementación de concurrencia

## Uso de threads
La concurrencia se implementó usando un modelo basado en threads (utilizando la librería `threading` de Python). Se le asignó un hilo a cada cliente, y cada uno procesa sus propias apuestas concurrentemente.

Este programa realiza principalmente operaciones I/O-bound (enviar/recibir datos por sockets y leer/escribir archivos), que ocurren por fuera del GIL. El GIL no sería ideal para operaciones CPU-bound.

## Métodos de sincronización
- **Condition**: se utiliza para esperar a las agencias hasta que se haya alcanzado el quórum mínimo (`AGENCY_QUORUM_MIN`, configurado como variable de entorno) de agencias que terminaron de enviar sus apuestas. Cuando se alcanza esta cantidad de agencias, se procesan y devuelven los ganadores. Los threads esperan bloqueados en `wait()`.
- **Lock de storage**: se utiliza para proteger el archivo compartido de escritura (`store_bets`) y lectura (`load_bets`). De esta manera se asegura que nunca va a haber una lectura o escritura concurrente en el archivo, evitando que se realicen dos escrituras al mismo tiempo o una lectura mientras hay una escritura en curso.