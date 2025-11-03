// personagem.go - entrada do jogador e envio de movimento ao servidor
package main

import (
	"fmt"
	"log"
	"time"
)

// personagemExecutarAcao processa teclas: sair, interagir, mover.
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		if jogo.Cliente != nil && jogo.ID != "" {
			var ack bool
			_ = jogo.Cliente.Call("Servidor.DesconectarJogador", jogo.ID, &ack)
		}
		return false

	case "interagir":
		jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
		jogo.StatusExp = time.Now().Add(2 * time.Second)
		return true

	case "mover":
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
		default:
			return true
		}

		nx, ny := jogo.PosX+dx, jogo.PosY+dy

		// colisão com inimigo: dano, sem mover
		if jogo.Mapa[ny][nx].simbolo == Inimigo.simbolo {
			jogoVerificarColisao(jogo, nx, ny)
			return true
		}

		// valida movimento
		if !jogoPodeMoverPara(jogo, nx, ny) {
			return true
		}

		// atualiza localmente (garante desenho imediato)
		jogo.PosX = nx
		jogo.PosY = ny

		// prepara movimento com sequence incremental
		jogo.seq++
		mov := Movimento{
			ID:       jogo.ID,
			PosX:     jogo.PosX,
			PosY:     jogo.PosY,
			Sequence: jogo.seq,
			Vidas:    jogo.vidas,
		}

		// tenta enviar até 3 vezes se servidor pedir retransmissão
		var ack bool
		for tentativa := 0; tentativa < 3; tentativa++ {
			if err := jogo.Cliente.Call("Servidor.AtualizarMovimento", mov, &ack); err != nil {
				log.Println("erro RPC:", err)
			}
			if ack {
				// aceito pelo servidor
				break
			}
			// servidor pediu retransmissão (ack==false) ou erro: re-tenta
			time.Sleep(50 * time.Millisecond)
		}
		return true
	}
	return true
}
