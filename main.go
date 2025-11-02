// main.go - Loop principal do jogo
package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	estadoAtual EstadoJogo
)

func main() {
	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento (playerID e opcionalmente mapa)
	if len(os.Args) < 2 {
		fmt.Println("Use: jogo.exe <playerID> [mapa.txt] [endereco_servidor]")
		return
	}
	id := os.Args[1]

	mapaFile := "mapa.txt"
	if len(os.Args) > 2 {
		mapaFile = os.Args[2]
	}

	// Lê endereço do servidor ANTES de iniciar o termbox
	var endereco string
	if len(os.Args) > 3 {
		endereco = os.Args[3]
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Digite o endereco do servidor (ex: 127.0.0.1:8932): ")
		txt, _ := reader.ReadString('\n')
		endereco = strings.TrimSpace(txt)
	}

	// Inicializa o jogo local (mapa no cliente)
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Conecta ao servidor
	client, err := rpc.Dial("tcp", endereco)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor RPC:", err)
		return
	}
	jogo.Cliente = client
	jogo.ID = id

	var ok bool
	if err := client.Call("Servidor.RegistrarJogador", id, &ok); err != nil || !ok {
		fmt.Println("Erro ao registrar jogador no servidor")
		return
	}

	// Agora inicia a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Goroutine para polling de estado global dos players
	go func() {
		tk := time.NewTicker(100 * time.Millisecond)
		defer tk.Stop()
		for range tk.C {
			var estado EstadoJogo
			if err := client.Call("Servidor.GetEstadoJogo", id, &estado); err != nil {
				continue
			}
			mu.Lock()
			estadoAtual = estado
			mu.Unlock()
			jogoAtualizarEstadoMultiplayer(&jogo, estado)
		}
	}()

	// Loop de renderização
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			estado := estadoAtual
			mu.Unlock()
			interfaceDesenharJogo(&jogo, estado)
		}
	}()

	// Loop principal de entrada
	go jogo.Run()
	for !jogo.encerrar {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar || jogo.encerrar {
			break
		}
	}

	defer interfaceLimparTela()
}
