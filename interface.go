// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
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

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
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

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo, estado EstadoJogo) {
	interfaceLimparTela()

	// desenha mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	// desenha outros jogadores
	for _, p := range estado.Players {
		if p.ID == jogo.ID {
			continue
		}
		interfaceDesenharElemento(p.PosX, p.PosY, OtherPlayer)
	}

	lifeMsg := fmt.Sprintf("Vidas: %d/3", jogo.vidas)
	for i, c := range lifeMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorVerde, CorPadrao)
	}

	// mostra barra de status
	if jogo.StatusMsg != "" && time.Now().Before(jogo.StatusExp) {
		for i, c := range jogo.StatusMsg {
			termbox.SetCell(i, len(jogo.Mapa)+3, c, CorMagenta, CorPadrao)
		}
	} else {
		// limpa mensagem
		jogo.StatusMsg = ""
	}

	// instruções fixas
	instr := "Use WASD para mover, E para interagir, ESC para sair. "
	for i, c := range instr {
		termbox.SetCell(i, len(jogo.Mapa)+5, c, CorTexto, CorPadrao)
	}

	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// flush
	interfaceAtualizarTela()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}
