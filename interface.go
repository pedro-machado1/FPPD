// interface.go - renderização com termbox e barra de status
package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type Cor = termbox.Attribute

const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorMagenta         = termbox.ColorMagenta
	CorVerde           = termbox.ColorGreen
	CorParede          = termbox.ColorBlack
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
)

type EventoTeclado struct {
	Tipo  string
	Tecla rune
}

func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

func interfaceFinalizar() { termbox.Close() }

func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

func interfaceDesenharJogo(jogo *Jogo, estado EstadoJogo) {
	interfaceLimparTela()

	// desenha mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceSetCell(x, y, elem)
		}
	}

	// desenha outros jogadores segundo o estado do servidor
	for _, p := range estado.Players {
		if p.ID == jogo.ID {
			// não desenha aqui próprio (será desenhado localmente)
			continue
		}
		interfaceSetCell(p.PosX, p.PosY, OtherPlayer)
	}

	// desenha inimigos/elementos já no mapa (mapa já tem isso)
	// mostra suas vidas
	lifeMsg := fmt.Sprintf("Vidas: %d/3", jogo.vidas)
	for i, c := range lifeMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorVerde, CorPadrao)
	}

	// mostra barra de status (se não expirou)
	if jogo.StatusMsg != "" && time.Now().Before(jogo.StatusExp) {
		for i, c := range jogo.StatusMsg {
			termbox.SetCell(i, len(jogo.Mapa)+3, c, CorMagenta, CorPadrao)
		}
	} else {
		// limpa mensagem se expirou
		jogo.StatusMsg = ""
	}

	// instruções fixas
	instr := "Use WASD para mover, E para interagir, ESC para sair. "
	for i, c := range instr {
		termbox.SetCell(i, len(jogo.Mapa)+5, c, CorTexto, CorPadrao)
	}

	// desenha o player local por último (garante visibilidade)
	interfaceSetCell(jogo.PosX, jogo.PosY, Personagem)

	// flush
	interfaceAtualizarTela()
}

func interfaceSetCell(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

func interfaceAtualizarTela() {
	termbox.Flush()
}
