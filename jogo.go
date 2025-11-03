// jogo.go - carregamento de mapa e utilitários do cliente
package main

import (
	"bufio"
	"net/rpc"
	"os"
	"time"
)

// EstadoPlayer/EstadoJogo/Movimento (exportadas para RPC se necessário)
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

// Elemento do mapa
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}

// Jogo local
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

var (
	Personagem  = Elemento{'☺', CorVerde, CorPadrao, false}
	OtherPlayer = Elemento{'☻', CorMagenta, CorPadrao, false}
	Parede      = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao   = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio       = Elemento{' ', CorPadrao, CorPadrao, false}
	Inimigo     = Elemento{'☠', CorVermelho, CorPadrao, true}
)

func jogoNovo() Jogo {
	return Jogo{
		Mapa:     [][]Elemento{},
		PosX:     1,
		PosY:     1,
		vidas:    3,
		encerrar: false,
		seq:      1,
	}
}

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

func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}
	if jogo.Mapa[y][x].tangivel {
		return false
	}
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
