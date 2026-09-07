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
La concurrencia se implementó utilizando un modelo basado en threads (mediante la librería `threading` de Python). Se le asigna un hilo a cada cliente conectado, para poder procesar sus propias apuestas de forma concurrente.

Se eligio usar threads ya que el programa es principalmente I/O-bound: la mayor parte del tiempo se dedica a enviar y recibir datos por sockets y a leer y escribir archivos, operaciones que liberan el GIL mientras esperan a completarse. En un escenario CPU-bound, en cambio, el GIL sí sería una limitación real, ya que impediría aprovechar múltiples núcleos para cómputo puro.

## Métodos de sincronización
- **Condition**: se utiliza para esperar a las agencias hasta que se haya alcanzado el quórum mínimo (`AGENCY_QUORUM_MIN`, configurado como variable de entorno) de agencias que terminaron de enviar sus apuestas. Cuando se alcanza esta cantidad de agencias, se procesan y devuelven los ganadores. Los threads esperan bloqueados en `wait()`.
- **Lock de storage**: se utiliza para proteger el archivo compartido de escritura (`store_bets`) y lectura (`load_bets`). De esta manera se asegura que nunca va a haber una lectura o escritura concurrente en el archivo, evitando que se realicen dos escrituras al mismo tiempo o una lectura mientras hay una escritura en curso.