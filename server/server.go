package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Informações jogador.
type EstadoPlayer struct {
	ID       string
	PosX     int
	PosY     int
	Sequence int
	Vidas    int
}

// estado global do jogo.
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

// Servidor central
type Servidor struct {
	mu     sync.Mutex
	estado EstadoJogo
}

// novo jogador.
func (s *Servidor) RegistrarJogador(id string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.estado.Players == nil {
		s.estado.Players = make(map[string]EstadoPlayer)
	}

	if _, existe := s.estado.Players[id]; !existe {
		s.estado.Players[id] = EstadoPlayer{
			ID:       id,
			PosX:     11,
			PosY:     6,
			Sequence: 0,
			Vidas:    3,
		}
		fmt.Printf("[%s] novo jogador registrado (%d jogadores agora)\n",
			id, len(s.estado.Players))
	}
	*reply = true
	return nil
}

// remove o jogador.
func (s *Servidor) DesconectarJogador(id string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.estado.Players != nil {
		delete(s.estado.Players, id)
		fmt.Printf("[%s] desconectado (%d jogadores restantes)\n",
			id, len(s.estado.Players))
	}
	*reply = true
	return nil
}

// envia o estado atual.
func (s *Servidor) GetEstadoJogo(id string, estado *EstadoJogo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	*estado = s.estado
	fmt.Printf("[%s] pediu estado (%d jogadores conectados)\n",
		id, len(s.estado.Players))
	return nil
}

// Atualiza a posição .
func (s *Servidor) AtualizarMovimento(mov Movimento, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.estado.Players[mov.ID]
	if !ok {
		*reply = false
		return fmt.Errorf("jogador não encontrado")
	}

	if mov.Sequence <= player.Sequence {
		fmt.Printf("[%s] seq inválida %d <= %d (ignorado)\n",
			mov.ID, mov.Sequence, player.Sequence)
		*reply = false
		return nil
	}

	// atualiza informações do jogador
	player.PosX = mov.PosX
	player.PosY = mov.PosY
	player.Sequence = mov.Sequence
	player.Vidas = mov.Vidas
	s.estado.Players[mov.ID] = player

	fmt.Printf("[%s] mov → (%d,%d) seq=%d vidas=%d [%s]\n",
		mov.ID, player.PosX, player.PosY, player.Sequence, player.Vidas,
		time.Now().Format("15:04:05"))

	*reply = true
	return nil
}

var srv = &Servidor{
	estado: EstadoJogo{
		Players: make(map[string]EstadoPlayer),
	},
}

// Função principal
func main() {
	if err := rpc.RegisterName("Servidor", srv); err != nil {
		log.Fatal("Falha ao registrar serviço:", err)
	}

	l, err := net.Listen("tcp", "0.0.0.0:8932")
	if err != nil {
		log.Fatal("Erro ao abrir porta 8932:", err)
	}
	log.Println("Servidor RPC do jogo em :8932")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
