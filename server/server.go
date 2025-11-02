package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// Tipos compartilhados pelo servidor com os clientes
type EstadoPlayer struct {
	ID         string
	PosX, PosY int
	Sequence   int
}

type EstadoJogo struct {
	Players     map[string]EstadoPlayer
	BotaoAtivo  bool
	PortalAtivo bool
}

type Movimento struct {
	ID       string
	PosX     int
	PosY     int
	Sequence int
}

// Serviço RPC do jogo
type Servidor struct {
	mu     sync.Mutex
	estado EstadoJogo
}

func (s *Servidor) RegistrarJogador(id string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.estado.Players == nil {
		s.estado.Players = make(map[string]EstadoPlayer)
	}
	if _, existe := s.estado.Players[id]; !existe {
		s.estado.Players[id] = EstadoPlayer{
			ID:       id,
			PosX:     1,
			PosY:     1,
			Sequence: 0,
		}
	}
	*reply = true
	return nil
}

func (s *Servidor) DesconectarJogador(id string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.estado.Players != nil {
		delete(s.estado.Players, id)
	}
	*reply = true
	return nil
}

func (s *Servidor) GetEstadoJogo(_ string, estado *EstadoJogo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	*estado = s.estado
	return nil
}

func (s *Servidor) AtualizarEstadoLogico(estadoLocal EstadoJogo, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.estado.BotaoAtivo = estadoLocal.BotaoAtivo
	s.estado.PortalAtivo = estadoLocal.PortalAtivo
	*reply = true
	return nil
}

func (s *Servidor) AtualizarMovimento(mov Movimento, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.estado.Players[mov.ID]
	if !ok {
		*reply = false
		return fmt.Errorf("jogador não encontrado")
	}
	// aplica somente movimentos mais recentes
	if mov.Sequence > player.Sequence {
		player.PosX = mov.PosX
		player.PosY = mov.PosY
		player.Sequence = mov.Sequence
		s.estado.Players[mov.ID] = player
	}
	*reply = true
	return nil
}

func main() {
	srv := new(Servidor)
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
