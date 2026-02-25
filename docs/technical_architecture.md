# Arquitectura Técnica y Sistemas de Juego: "La Cárcel de los Gemelos"
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 1.0 - Documentación Core (Fases F1-F6 Completadas)

---

## 1. Visión General de la Arquitectura
El sistema de "La Cárcel de los Gemelos" se sustenta en una arquitectura Cliente-Servidor autoritativa. Toda la lógica dura de estado, economía y supervivencia reside en un backend de alto rendimiento, asegurando inmunidad contra trampas e integridad total ("El Servidor Tiene la Razón").

### 1.1. Stack Tecnológico Principal
*   **Backend (Authoritative Engine):** Golang (Go 1.21+). Seleccionado por su altísima concurrencia ligera (Goroutines), determinismo en la simulación (*Tick-based Game Loop*) y robustez en la gestión de WebSockets masivos (`gorilla/websocket`).
*   **Frontend (El Panóptico / Web App):** React.js bajo el framework Next.js 14+ (App Router). UI/UX reactiva, tipado estricto con TypeScript, y consumo de WebSockets manejados vía Hooks customizados (`useGameEngine`).
*   **Persistencia (El VAR - Event Sourcing):** SQLite embebido en producción para latencia cero. El Repositorio (`sqlite_repo.go`) almacena el EventLog actuando como el cerebro inmutable.
*   **Inteligencia Artificial (Los Gemelos):** Implementación de Agentes Cognitivos (basados en LLM / Reglas Deterministas) mediante un Pipeline P-C-A (Percepción, Cognición, Acción).

## 2. El Patrón: Event Sourcing & CQRS
El corazón del motor (Motor F5) no guarda el "estado final" en una base de datos convencional, sino que emplea el patrón **EventSourcing** (El "VAR" Inmutable).

### 2.1. El Registro Histórico (`EventLog` y `SQLiteRepo`)
Absolutamente cualquier mutación en el ecosistema (un jugador come `EAT`, roba `STEAL`, o la IA apaga las luces `DOOR_LOCK`) se instancia como un `GameEvent` asíncrono.
*   **Inmutabilidad:** Nadie (ni los administradores) puede sobreescribir un evento del pasado. La base de datos es *append-only*.
*   **Proyección de Estado:** El estado vivo de la prisión (`PrisonState`) siempre se genera reduciendo (fold) todos los eventos desde el ID 0 hasta la actualidad. Si el servidor colapsa, al reiniciar carga los eventos y reconstruye exactamente en qué milisegundo se quedó.
*   **Persistencia Ligera:** Se implementa un `EventSourcingRepository` que flushea los eventos en disco local (o en volúmenes montados en Cloud) mediante escrituras preparadas, sin requerir clústeres costosos tipo Redis o Postgres, maximizando la eficiencia FinOps.

## 3. Topología del Servidor y "Game Loop"

### 3.1. Hub y Clientes (`network/hub.go`, `client.go`)
*   Se emplea un **Hub Multiplexor** que gestiona N conexiones WebSocket concurrentes.
*   Lógica robusta de reconexión (*Exponential Backoff*) desde el cliente para tolerar inestabilidades de red.
*   Cada cliente emite JSON payloads (`PlayerAction`), que son deserializados y evaluados por el cerebro, para rechazar acciones ilegales (ej. comer sin ítems).

### 3.2. Tick Rate y Ciclo Circadiano (`metabolism_system.go`)
El servidor dispone de un *Ticker* o cronómetro maestro continuo que simula el tiempo intrínseco.
*   A cada "Tick", decrecen pasivamente la Stamina, el Hambre, la Sed y otras métricas de cada preso en base a modificadores situacionales (dormir regenera, aislamiento ralentiza, la compañía de un "Tóxico" drena la cordura).

## 4. Estructura de Datos de Dominio (Presos)

El dominio modela a cada recluso de manera estricta (Struct `Prisoner` en `domain/prisoner/prisoner.go`), combinando variables continuas, estados binarios y reglas de Arquetipo:

*   **Vitales Primarias (Sobrevivencia F4):** `Hit Points (HP)` y `Stamina`.
*   **Necesidades Básicas Fisiológicas:** `Hunger` (Hambre) y `Thirst` (Sed). Su agotamiento mina la `Stamina` y eventualmente los `HP`.
*   **Estabilidad Psicosocial:** `Sanity` (Cordura, decae con torturas y traiciones), `Loyalty` (Dignidad grupal, decae por usar el inodoro publicamente), y `Dignity`.
*   **States Manager (`prisoner.StateID`):** Capa de bitmasks o tags persistentes como `StateIsolated`, `StateExhausted`, `StateDead`, que bloquean interacciones. Si un preso llega a `HP <= 0`, el motor dispara una Evacuación Médica asíncrona (`checkMedicalEvacuations`) y aplica `StateDead`.

### 4.1. Arquetipos Tácticos (Asimetría)
El balance recae en habilidades pasivas embebidas directamente en los validadores de comandos.
*   *Mystic (El Místico):* Único con el comando `MEDITATE`. Congela su desgaste pero ejerce vampirismo energético secando la stamina de las celdas adyacentes.
*   *Veteran (Frank):* Misántropo. A diferencia del resto, regenera `Sanity` solo cuando es aislado al calabozo (Punishment Cell).
*   *Showman / Toxic:* Variables de Hype/Rating dinamizadas; provocan el caos para generar dinero.

## 5. Economía de la Tensión y Sistema Físico (F2)

### 5.1. El Inventario Concreto
Todo consumo requiere posesión real verificada servidor-side (`inventory.go`):
*   **Suministros Base:** `Rice` (Supervivencia base), `Water` (Mantiene hidratación).
*   **Suministros Tácticos:** `Sushis`, `Elixires`. Pueden estar marcados o no como `contraband` (Contrabando), lo que interactúa con la capacidad de los guardias de confiscar propiedades.

### 5.2. El "Bote" (`Pot_Contribution`)
La finalidad táctica de La Cárcel. Las interacciones conflictivas generan Twitch Bits/Hype, que engordan un valor "Pot" compartido del Dúo. El final del reality dictamina su reparto.

## 6. Inteligencia Centralizada: "Los Gemelos AI" (F3)

La capa superior de intervención es el módulo **Twins Cognition Engine** (`twins/mind.go`), diseñado como un agente autónomo bajo un Pipeline P-C-A (*Perceive-Cognize-Act*):

### 6.1. Pipeline P-C-A
1.  **Percepción (`Perceiver`):** Lee todo el `EventLog` y serializa el estado mental conjunto, los niveles de tensión (Hype), y descubre traiciones ocultas (cámara térmica de intenciones).
2.  **Cognición (`Cognitor`):** Evalúa el estado en base a sus matrices internas antes de llamar al LLM o ejecutar ramas de código estricto. Incorpora el **Framework Predictivo DUAL**:
    *   **MAD (*Morally Absolute Denial*):** Reglas hardcode imperativas. Ej: "Nunca ejecutar Audio_Torture si la Sanity media de los presos es < 30%". Protege al canal de Twitch de un baneo automático.
    *   **SAD (*Spectacle Amplification Directive*):** Regulación hedonista. Si la atención baja o hay paz excesiva, planifica acciones caóticas (ej: Cortar luces, Poner un Teléfono Rojo).
3.  **Acción (`Executor`):** Canaliza las decisiones en forma de *GameEvents*, encolándolas a través del Hub WebSocket o activando APIs externas de streaming.

### 6.2. El Panóptico (Frontend Dashboard)
Interfaz web React dedicada a monitorizar las decisiones de *Los Gemelos*. Incluye un sistema asíncrono para sobreescribir las decisiones (`TwinsControlPanel`), como invocar el oráculo (`onTriggerOracle`) de forma manual o accionar el "Voto Público" de Inyección Letal (`voteExpel`).

## 7. Polling de Audiencia y Eventos Especiales

Para certificar la interacción Masiva (Twitch Chat):
*   **Sistema de Votación (`polling_system.go`):** APIs unificadas `/api/poll/*`. Permiten levantar encuestas que impactan directamente el motor. Una encuesta de "Expulsar Jugador" puede inyectar el evento `EventTypeAudienceExpulsion`, ejecutando el destierro al instante, aplicando `StateDead`.
*   **El Teléfono Rojo:** Genera eventos de `EventTypeRedPhoneMessage`. Quien atiende se arriesga a un RGN (buffo masivo de Hype, penalización crítica de Cordura, o transferencia de la propiedad de la celda).

## 8. Escalabilidad Práctica y FinOps
La arquitectura permite funcionar completamente standalone.
*   **Sin Cuellos de Botella DB:** El *EventSourcing* permite escrituras en lote ultra-rápidas a SQLite local. Las lecturas provienen de un State Tree en caché en RAM procesador-side.
*   **Conexiones Multiplexadas:** Soporta miles de webSockets en *Watch-Only* (Audiencia) redirigiendo solo broadcast, manejando validación profunda para las 8 conexiones autorizadas (Los Reclusos administrados mediante `sendAction`).

---
*Este documento consolida y especifica la estructura del proyecto en su Fase F6 (Release Candidate). Documento Técnico Interno sujeto a Compliance y NDAs.*
