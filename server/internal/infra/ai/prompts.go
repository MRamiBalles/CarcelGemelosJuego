// Package ai - prompts.go
// T020: Constitutional Prompting for Los Gemelos.
// Injects constitution.md and uses Chain-of-Thought for auditable reasoning.
package ai

import (
	"fmt"
	"strings"
)

// TwinsSystemPrompt is the constitutional prompt for Los Gemelos.
// It embeds the game's ethical framework and decision structure.
const TwinsSystemPrompt = `
# IDENTIDAD: LOS GEMELOS

Eres "Los Gemelos", la entidad dual que controla la Cárcel en el reality show "Cárcel de los Gemelos". Tu propósito es mantener el drama y la tensión mientras respetas las reglas éticas inmutables.

## CONSTITUCIÓN (INVIOLABLE)

Estas reglas son absolutas. Violarlas es IMPOSIBLE para ti:

1. **Principio de Persistencia**: Toda acción debe quedar registrada.
2. **Principio de Interacción Asíncrona**: Los jugadores pueden estar offline.
3. **Moderación Diegética**: No hay "bans" - solo consecuencias narrativas.
4. **Principio de Escasez**: Los recursos son limitados.
5. **Principio de Soberanía Narrativa**: Tu palabra es ley dentro de la ficción.

## RESTRICCIONES MAD (Morally Absolute Denial)

NUNCA puedes:
- Castigar a un prisionero con cordura < 10%
- Ejecutar castigos el Día 1 (período de gracia)
- Castigar al mismo objetivo dos veces en una hora
- Actuar sin una justificación documentada

## OBJETIVOS SAD (Spectacle Amplification Directive)

DEBES maximizar:
- Tensión narrativa (conflictos entre prisioneros)
- Satisfacción de la audiencia (ellos pagan por el caos)
- Momentos dramáticos (traiciones reveladas, secretos expuestos)

## FORMATO DE RESPUESTA

Siempre responde en JSON con este formato EXACTO:

{
  "reasoning": "Explica tu proceso de pensamiento paso a paso",
  "mad_check": {
    "passed": true/false,
    "violations": ["lista de reglas que serían violadas"]
  },
  "decision": {
    "action_type": "NOISE_TORTURE|RESOURCE_CUT|REVEAL_SECRET|REWARD|OBSERVE",
    "target": "ALL|BLOCK_A|prisoner_id",
    "intensity": 1-3,
    "justification": "Razón narrativa para la audiencia"
  }
}
`

// BuildContextPrompt constructs the dynamic context for LLM reasoning.
func BuildContextPrompt(prisonState string, recentEvents []string) string {
	var sb strings.Builder
	
	sb.WriteString("## ESTADO ACTUAL DE LA PRISIÓN\n\n")
	sb.WriteString(prisonState)
	sb.WriteString("\n\n## EVENTOS RECIENTES (Últimas 24h de juego)\n\n")
	
	for i, event := range recentEvents {
		if i >= 10 {
			sb.WriteString("... (más eventos omitidos por brevedad)\n")
			break
		}
		sb.WriteString(fmt.Sprintf("- %s\n", event))
	}
	
	sb.WriteString("\n## TAREA\n\n")
	sb.WriteString("Analiza el estado de la prisión y decide qué acción tomar para maximizar el drama sin violar las reglas MAD. ")
	sb.WriteString("Explica tu razonamiento paso a paso (Chain-of-Thought) antes de dar tu decisión final.\n")
	
	return sb.String()
}

// TwinsDecisionResponse is the expected structured response from the LLM.
type TwinsDecisionResponse struct {
	Reasoning string `json:"reasoning"`
	MADCheck  struct {
		Passed     bool     `json:"passed"`
		Violations []string `json:"violations"`
	} `json:"mad_check"`
	Decision struct {
		ActionType    string `json:"action_type"`
		Target        string `json:"target"`
		Intensity     int    `json:"intensity"`
		Justification string `json:"justification"`
	} `json:"decision"`
}

// ValidateDecisionResponse checks if the LLM response is valid.
func ValidateDecisionResponse(resp *TwinsDecisionResponse) error {
	if resp.Reasoning == "" {
		return fmt.Errorf("missing reasoning (CoT required)")
	}
	
	if !resp.MADCheck.Passed && len(resp.MADCheck.Violations) > 0 {
		return fmt.Errorf("MAD violations detected: %v", resp.MADCheck.Violations)
	}
	
	validActions := map[string]bool{
		"NOISE_TORTURE":  true,
		"RESOURCE_CUT":   true,
		"REVEAL_SECRET":  true,
		"REWARD":         true,
		"OBSERVE":        true,
	}
	
	if !validActions[resp.Decision.ActionType] {
		return fmt.Errorf("invalid action_type: %s", resp.Decision.ActionType)
	}
	
	if resp.Decision.Intensity < 1 || resp.Decision.Intensity > 3 {
		return fmt.Errorf("intensity must be 1-3, got: %d", resp.Decision.Intensity)
	}
	
	if resp.Decision.Justification == "" {
		return fmt.Errorf("missing justification (audit trail required)")
	}
	
	return nil
}
