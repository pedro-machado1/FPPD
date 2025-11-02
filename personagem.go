// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"fmt"
	"log"
	"time"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		if jogo.Cliente != nil && jogo.ID != "" {
			var ack bool
			_ = jogo.Cliente.Call("Servidor.DesconectarJogador", jogo.ID, &ack)
		}
		return false
	case "interagir":
		// Apenas mensagem local e exemplo de sinalização lógica (botão/portal)
		jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
		// Caso queira ativar/desativar algo global:
		if jogo.Cliente != nil {
			var ack bool
			_ = jogo.Cliente.Call("Servidor.AtualizarEstadoLogico", EstadoJogo{
				BotaoAtivo:  true,
				PortalAtivo: false,
			}, &ack)
		}
	case "mover":
		// Move localmente (colisão local com mapa do cliente)
		oldX, oldY := jogo.PosX, jogo.PosY
		dx, dy := 0, 0
		switch ev.Tecla {
		case 'w':
			dy = -1
		case 'a':
			dx = -1
		case 's':
			dy = 1
		case 'd':
			dx = 1
		}
		nx, ny := jogo.PosX+dx, jogo.PosY+dy
		if jogoPodeMoverPara(jogo, nx, ny) {
			jogo.chMapa <- MapaUpdate{
				tipo: "Personagem",
				fx:   oldX, fy: oldY,
				tx: nx, ty: ny,
			}
		}

		// Envia posição para o servidor (sequência simples baseada em tempo)
		if jogo.Cliente != nil && jogo.ID != "" {
			seq := int(time.Now().UnixNano())
			mov := Movimento{
				ID:       jogo.ID,
				PosX:     nx,
				PosY:     ny,
				Sequence: seq,
			}
			var ack bool
			if err := jogo.Cliente.Call("Servidor.AtualizarMovimento", mov, &ack); err != nil {
				log.Println("Erro ao atualizar posição no servidor:", err)
			}
		}
	}
	return true
}
