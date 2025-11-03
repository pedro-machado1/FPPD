// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"net/rpc"
	"os"
	"time"
)

type EstadoPlayer struct {
	ID       string
	PosX     int
	PosY     int
	Sequence int
	Vidas    int
}

type EstadoJogo struct {
	Players map[string]EstadoPlayer
}

type Movimento struct {
	ID       string
	PosX     int
	PosY     int
	Sequence int
	Vidas    int
}

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa      [][]Elemento
	PosX      int
	PosY      int
	encerrar  bool
	StatusMsg string
	StatusExp time.Time // quando a mensagem deve sumir
	vidas     int
	seq       int
	Cliente   *rpc.Client
	ID        string
}

// Elementos visuais do jogo
var (
	Personagem  = Elemento{'☺', CorVerde, CorPadrao, false}
	OtherPlayer = Elemento{'☻', CorMagenta, CorPadrao, false}
	Parede      = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao   = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio       = Elemento{' ', CorPadrao, CorPadrao, false}
	Inimigo     = Elemento{'☠', CorVermelho, CorPadrao, true}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{
		Mapa:     [][]Elemento{},
		PosX:     1,
		PosY:     1,
		vidas:    3,
		encerrar: false,
		seq:      1,
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y
				e = Vazio
			case Inimigo.simbolo:
				e = Inimigo
			default:
				e = Vazio
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

func jogoVerificarColisao(jogo *Jogo, x, y int) bool {
	if jogo.Mapa[y][x].simbolo == Inimigo.simbolo {
		jogo.vidas--
		if jogo.vidas <= 0 {
			jogo.encerrar = true
		}
		return true
	}
	return false
}
