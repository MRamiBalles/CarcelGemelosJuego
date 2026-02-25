# Game Design Document (GDD): "La Cárcel de los Gemelos"
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 3.0 (Consolidación Fases F1-F6)

---

## 1. Concepto Core: Interactive Reality RPG
"La Cárcel de los Gemelos" redefine el formato de "Prison Sim" transformándolo en un **"Multiplayer Prisoner's Dilemma RPG"**.
*   **Target Principal:** Ecosistema de creadores de contenido (Twitch).
*   **Meta Final:** Sobrevivir 21 días cronometrados reales (Mar 15 - Apr 5).
*   **La Paradoja del Dúo:** Los jugadores conviven y sobreviven en celdas de a dos (Duos), pero **solo un individuo puede reclamar el bote final**.
*   **Los Antagonistas:** La *Convivencia* (El desgaste del roce), la *Audiencia* (Pay-to-Torture) y *Los Gemelos* (La IA omnipresente que modera y castiga).

## 2. Pila Fisiológica y Psicológica (F4 / F5)
Los avatares o reclusos no son inmortales; su ciclo de vida está simulado y debe mantenerse equilibrado:
*   **Hit Points (HP):** Si caen a 0, se aplica el protocolo de Evacuación Médica (`MedEvac`) y el preso es marcado como `StateDead` (expulsado del concurso).
*   **Stamina (Energía):** Se drena diariamente y se consume masivamente en Desafíos o al ejecutar actos de Robo (`STEAL`). Si cae al 10%, el preso entra en `StateExhausted` (Duplica daño recibido, 80% de fallo en acciones). Se recupera durmiendo en las horas de *Lockdown*.
*   **Hambre y Sed (`Hunger` / `Thirst`):** Se drena constante por tick. Alcanzar la inanición bloquea la recuperación de Stamina. Requiere el consumo físico de `RICE` y `WATER` del inventario.
*   **Cordura (`Sanity`):** Barra mental. Decae por torturas acústicas, traiciones, y falta de sueño. Cuando es crítica (< 30%), Los Gemelos pueden activar protocolos MAD para proteger al preso.
*   **Dignidad y Lealtad:** La dignidad se arruina al usar el inodoro (`TOILET`) frente a compañeros. La lealtad se destroza mediante chivatazos (`SNITCH`) y robos cruzados.

## 3. Entorno Hardcore y Rutinas
*   **El Módulo de Aislamiento:** Una celda de tortura aislada. Los Gemelos enviarán a los infractores allí (`StateIsolated`). Bloquea interacciones, anula el sueño y duplica los daños de ruido.
*   **Lockdown (Toque de Queda):** Entre las 00:00 y las 08:00, las celdas se bloquean. Tiempo crítico para regenerar Stamina y evitar hacer ruido.
*   **Visión Cruzada (Line of Sight):** Arquitectura Panóptica. Las celdas enfrentadas (A-D, B-E, C-F) permiten chivatazos (Snitches) reales a la IA si un contrario saca contrabando.

## 4. Asimetría de Arquetipos (Los Participantes)

La jugabilidad recae en pasivas ancladas al código para garantizar roles asimétricos:

### A. El Veterano (Frank Cuesta)
*   **Pasiva "Estómago de Hierro":** Inmunidad pasiva a intoxicaciones o dietas crudas.
*   **Pasiva "Misántropo":** Único arquetipo capaz de regenerar `Sanity` dentro del Módulo de Aislamiento. Necesita castigarse para sanar mentalmente.

### B. Los Tóxicos (Labrador & Ylenia)
*   **Pasiva "Bad Romance":** Sinergia negativa. Generan ingresos explosivos ("Hype") mediante conflictos verbales y agresiones entre ellos. Sin embargo, su proximidad física pasiva drena su Cordura el doble de rápido. Deben acercarse para ganar dinero y huir para no reventar psiquiátricamente.

### C. El Místico (Tartaria)
*   **La Dieta del Aire (Breatharian):** Barra de Hambre bloqueada. Si consume comida sólida, el servidor aplica un penalizador de muerte casi inminente.
*   **Activa "Meditación Vampírica":** Posee el comando `MEDITATE`. Al usarlo, congela su propio desgaste corporal pero aturde de frío y seca la energía (Stamina) de las celdas vecinas.
*   **Activa "Oráculo":** Destruye pasivamente la Lealtad de otro jugador inyectando disonancias.

### D. Los Agentes del Caos
*   **Aída (La Insomne Poltergeist):** Reduce al 50% su necesidad de sueño. Acción activa: generar pánico a media noche anulando las horas de sueño de compañeros.
*   **Dakota (Mecha Corta):** El daño es exponencial cuando su barra de cordura está por debajo del 30%.
*   **Héctor:** Puede transferir pequeñas dosis de contrabando saltándose la auditoría inmediata del sistema de `Snitch` (delay del VAR de 12 horas).

## 5. El Inventario Físico y la Economía de la Audiencia (F2 / F6)

*   **Inventario del Preso:** Los ítems existen físicamente. Contrabando (Cigarrillos, Teléfonos, Sushi) versus Subsistencia (Arroz, Agua).
*   **El Bote ("The Pot" y Hype):** La divisa transaccional. Las visualizaciones y la tensión generan dinero.
*   **La Inyección Letal de la Audiencia (`Audience_Expulsion`):** La comunidad puede votar mediante integraciones de API (`/api/audience/vote_expel`). La IA aplica la justicia del pueblo castigando a muerte (`StateDead`) al preso elegido, vaciando su bote.

## 6. End Game: El Dilema Final
El Día 21 marca el cese cronológico del servidor.
La "Bóveda" (Vault) de la prisión se abre. Los ocupantes supervivientes de la celda acceden al Dilema sobre el "Bote" acumulado:
*   **Robar (Steal):** El que roba se queda el 100%. Si ambos roban, ambos colapsan a 0.
*   **Repartir (Split):** 50/50. Sin embargo, bajo la política de *Los Gemelos*, un Split pacífico puede ser tasado y castigado comercialmente por aburrido, premiando siempre el drama (El formato Reality Show).

---
*Documento Marco Consolidado (Versión 3.0)*
