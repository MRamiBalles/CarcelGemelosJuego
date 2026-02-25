# Referencia de API y WebSockets (Fase F6)
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 1.0 - Confidencial (Engine V2.1)

---

## 1. Introducción
Este documento define el contrato de integración entre el **Cognition Engine de Los Gemelos (Go Backend)** y cualquier cliente autorizado (Frontend Next.js Panel, Apps de Audiencia, o Simuladores). El sistema provee un canal de baja latencia vía WebSockets para los comandos del simulador en tiempo real, y una API REST autoritativa para integraciones externas y back-office.

---

## 2. Protocolo WebSocket (`/ws`)

### 2.1. Conexión y Autenticación
*   **Endpoint:** `ws://<host>:<port>/ws`
*   **Protocolo:** JSON. 
*   **Comportamiento:** La conexión es dual. Para administradores recibe telemetría continua de *todos* los eventos en tiempo real. Los jugadores (presos) usan el mismo ducto emitiendo paquetes `PlayerAction`.

### 2.2. Estructura de Mensaje Entrante (Cliente -> Servidor)
Cualquier solicitud que mute el estado debe respetar este formato:

```json
{
  "type": "NOMBRE_DEL_COMANDO",
  "prisoner_id": "P001",
  "payload": {
    "key": "value"
  }
}
```

### 2.3. Comandos Soportados (`PlayerAction`)

| Tipo de Acción (`type`) | Payload Esperado | Reglas / Restricciones | Resultado Emitido |
| :--- | :--- | :--- | :--- |
| `EAT` | `{"item_type": "RICE"}` | El ítem debe existir en el inventario del remitente. Si está Exhausto, la comida penaliza en vez de curar. | `EventTypeItemConsumed`, restaura `Hambre`. |
| `TOILET` | `{}` | Reduce Dignidad en 10. Si el Cellmate está observando, este sufre daño de tensión. | `EventTypeToiletUse` |
| `STEAL` | `{"target_cell": "CELL_A"}`o target directo | Requiere Stamina alta. 80% de fallo si se está en `StateExhausted`. Destruye la lealtad de la víctima. | `EventTypeSteal` o `EventTypeBetrayal`. |
| `SNITCH` | `{"target_id": "P002", "action": "SUSPICIOUS"}` | Evalúa el Inventario de la víctima. Si el "chivatazo" es real y hay contrabando, roba el bote (Pot). | `EventTypeSnitch` y transferencia de fondos. |
| `USE_RED_PHONE` | `{}` | Solo admitido si el Teléfono está sonando. Sistema de Ruleta Rusa (+Hype, +HP o daño crítico). | `EventTypeRedPhoneAnswer` |
| `MEDITATE` | `{"target_cells": ["CELL_B"]}` | **Exclusivo ArchetypeMystic.** Congela drenaje propio, aplica debuff a celdas objetivo (vampirismo de Stamina). | `EventTypeMeditate` e induce aislamientos. |
| `USE_ORACLE` | `{"target_id": "P003"}` | **Exclusivo ArchetypeMystic.** Golpea críticamente la `Lealtad` del receptor, extrayendo inteligencia del sistema. | `EventTypeOracleUse` |

### 2.4. Estructura de Mensaje Saliente (Servidor -> Cliente)
El servidor emite `GameEvent` inmutables provenientes del **EventLog (VAR)**.
```json
{
  "id": "evt_0001",
  "timestamp": "2026-03-15T10:00:00Z",
  "type": "STEAL",
  "actor_id": "P001",
  "target_id": "P002",
  "payload": "Stolen 1 RICE",
  "is_revealed": true
}
```

---

## 3. Integración REST API (Endpoints)

La API REST actúa sobre `/api/*` y es el punto de entrada para integraciones Web3, extensiones de Twitch o Paneles de Control de "La Audiencia" y Administradores.

### 3.1. VAR Replay (Event Sourcing)
*   **`GET /api/var/replay`**: Retorna la secuencia completa inmutable del universo desde el ID 0. Ideal para reconstrucción del cliente tras reconexiones.
*   **`GET /api/var/event?id={event_id}`**: Busca un evento específico.

### 3.2. Twins Overrides (Modo Dios / Panóptico)
*   **`POST /api/trigger-oracle`**:
    *   *Body:* `{"target": "P001", "message": "El fin se acerca"}`
    *   *Acción:* La IA inyecta un secreto destructivo en la mente del target.
*   **`POST /api/trigger-torture`**:
    *   *Body:* `{"soundId": "BABY_CRY_01"}`
    *   *Acción:* Salta el sistema de audio forzoso para los clientes físicos.

### 3.3. Audiencia y Twitch Polling
Para delegar la justicia a la audiencia en directo:
*   **`POST /api/poll/create`**: Inicializa una encuesta con X opciones.
*   **`POST /api/poll/vote`**: Computa y encripta el voto. *Rate-limited* por IP o token de espectador.
*   **`POST /api/audience/vote_expel`**:
    *   *Body:* `{"prisoner_id": "P006"}`
    *   *Acción:* El botón rojo nuclear. Emite un `EventTypeAudienceExpulsion` que fuerza al motor a purgar al preso, aplicándole permanentemente el estado `StateDead` (inmunidad a cálculos posteriores) hasta que se realice el reseteo de la base de datos.

---
*Fin del Documento de Referencia API. La utilización indebida o abuso de los endpoints no documentados puede disparar las directivas MAD (Mutual Assured Destruction) de la IA Gestora.*
