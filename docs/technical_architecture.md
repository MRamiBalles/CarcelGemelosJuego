# Arquitectura Técnica y Sistemas de Juego: "La Cárcel de los Gemelos"
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 1.0 - Documentación Core (Fases F1-F6 Completadas)

---

## 1. Visión General de la Arquitectura
El sistema de "La Cárcel de los Gemelos" se sustenta en una arquitectura Cliente-Servidor autoritativa. Toda la lógica dura de estado, economía y supervivencia reside en un backend de alto rendimiento, asegurando inmunidad contra trampas e integridad total ("El Servidor Tiene la Razón").

### 1.1. Stack Tecnológico
*   **Backend (Game Engine):** Golang (Go 1.21+). Seleccionado por su altísima concurrencia (Goroutines), ideal para simulaciones de físicas, *ticks* temporales (Game Loop) y webSockets.
*   **Frontend (El Panóptico / Web App):** React.js bajo el framework Next.js 14+ (App Router). UI/UX moderna gobernada por Tailwind CSS y tipado estricto con TypeScript.
*   **Persistencia (El VAR):** SQLite local vía EventSourcing / Arquitectura de Log Inmutable.
*   **Red:** WebSockets (Gorilla/Websocket) persistentes y bi-direccionales + API REST estandarizada para integraciones asíncronas (Twitch, Paneles Administrativos).

## 2. El Patrón: Event Sourcing & CQRS
El corazón del motor (Motor F5) no guarda el "estado final" en una base de datos convencional, sino que emplea el patrón **EventSourcing** (El "VAR" Inmutable).

### 2.1. El Registro Histórico (`EventLog`)
Cada acción en el juego, sin importar si es del usuario (EAT, TOILET, STEAL) o de la inteligencia artificial administradora (NOISE_TORTURE), se convierte en un `GameEvent` asíncrono.
*   Nadie puede eliminar ni sobreescribir un evento del pasado.
*   El estado actual de la prisión se reconstruye desde el Event 0 hasta la actualidad ("Replaying the Tape"), lo que permite resolver disputas, auditar trampas y recuperar el sistema perfectamente si el servidor colapsa.

## 3. Topología del Servidor y "Game Loop"

### 3.1. Hub y Clientes (`network/hub.go`, `client.go`)
*   Se emplea un **Hub Multiplexor** que gestiona N conexiones WebSocket concurrentes.
*   Lógica robusta de reconexión (*Exponential Backoff*) desde el cliente para tolerar inestabilidades de red.
*   Cada cliente emite JSON payloads (`PlayerAction`), que son deserializados y evaluados por el cerebro, para rechazar acciones ilegales (ej. comer sin ítems).

### 3.2. Tick Rate y Ciclo Circadiano (`metabolism_system.go`)
El servidor dispone de un *Ticker* o cronómetro maestro continuo que simula el tiempo intrínseco.
*   A cada "Tick", decrecen pasivamente la Stamina, el Hambre, la Sed y otras métricas de cada preso en base a modificadores situacionales (dormir regenera, aislamiento ralentiza, la compañía de un "Tóxico" drena la cordura).

## 4. Estructura de Datos de Dominio (Presos)

Cada recluso es un objeto complejo (Struct `Prisoner` en `domain/prisoner.go`), definido por Variables de Estado Constante y Arquetipos (F3):
*   **Vitales Primarias (F4):** Hit Points (HP), Stamina.
*   **Necesidades Básicas (F5):** Hunger (Hambre), Thirst (Sed).
*   **Estabilidad Psicosocial:** Sanity (Cordura), Loyalty (Lealtad), Dignity (Dignidad).
*   **Arquetipos:** Cada preso tiene un rol que influye hard-code en sus lógicas.
    *   *Místico (Mystic):* Su habilidad de meditar congela sus *stats* pero roba energía de celdas vecinas (Vampirismo Energético).
    *   *Veterano (Veteran):* Regenera cordura cuando se le castiga en el módulo de asilamiento ("Misántropo").
    *   *Showman / Tóxico:* Dinámicas aceleradas por las interacciones agresivas (Hype/Rating).

## 5. Almacenamiento y Sistema Físico (Inventario F2)

El sistema de recolección abandona las variables booleanas temporales e incorpora Objetos Estrictos (`inventory.go`).
*   **Capacidad Limitada:** Contrabando frente a raciones oficiales. Guardar el objeto incrementa su caché, pero si es de tipo "Ilegal" puede ser requisado.
*   **El Pot / Hype:** Divisa compartida entre los dúos generada por crear polémicas. El "Bote de la Victoria" se reparte (o se roba) en el Día 21.

## 6. Inteligencia Centralizada: "Los Gemelos AI" (F3)

La capa superior de intervención es el módulo *Twins Cognition Engine*. No es sólo un backend reactivo, sino un actor pro-activo.
*   **Core MAD-BAD-SAD:** El algoritmo de Los Gemelos evalúa el `TensionLevel` de la prisión. Si el *engagement* cae, invoca protocolos punitivos.
    *   *MAD (Mutual Assured Destruction):* Bloquea a la IA y a la audiencia ejecutar ciertas crueldades si la cordura media de la prisión baja peligrosamente del 30%, previniendo la catástrofe sistémica y cierres (Baneos de Twitch).
*   **El Panóptico (Interfaz Frontend):** El cliente de los Gemelos provee herramientas en tiempo real para disparar Eventos de Sirenas, Intervenciones Orales (Teléfono Rojo) o expulsar presos por aclamación, sin tener que tocar código ni SSH.

## 7. Escalabilidad Práctica (FinOps)
Para sostener la carga de docenas de eventos concurrentes interactuando simultáneamente con la base de datos de espectadores de streaming:
*   Mantenimiento del `State` en memoria en las tablas hash del propio Binario de Go.
*   Ausencia de escrituras asíncronas bloqueantes pesadas a Bases de Datos Lentas: se persisten "Lotes" de eventos.

---
*Este documento consolida técnicamente la infraestructura en la Fase F6. Propiedad Intelectual Confidencial de M. Ramírez Ballesteros.*
