// main.go - cliente: loop principal (sincroniza, desenha e processa input)
package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: go run . <playerID> [mapa.txt] [serverAddr]")
		return
	}
	id := os.Args[1]

	mapaFile := "mapa.txt"
	if len(os.Args) > 2 {
		mapaFile = os.Args[2]
	}

	endereco := "127.0.0.1:8932"
	if len(os.Args) > 3 {
		endereco = os.Args[3]
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Endereço do servidor (ENTER para %s): ", endereco)
		txt, _ := reader.ReadString('\n')
		txt = strings.TrimSpace(txt)
		if txt != "" {
			endereco = txt
		}
	}

	// preparar jogo local
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		fmt.Println("Erro carregar mapa:", err)
		return
	}

	// conectar servidor RPC
	client, err := rpc.Dial("tcp", endereco)
	if err != nil {
		fmt.Println("Erro conectar RPC:", err)
		return
	}
	jogo.Cliente = client
	jogo.ID = id

	var ok bool
	if err := client.Call("Servidor.RegistrarJogador", id, &ok); err != nil || !ok {
		fmt.Println("Erro registrar jogador no servidor")
		return
	}

	// iniciar interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// loop principal: buscar estado, desenhar, ler input, processar
	for !jogo.encerrar {
		// busca estado do servidor (timeout leve para não travar)
		var estado EstadoJogo
		_ = client.Call("Servidor.GetEstadoJogo", id, &estado)
		interfaceDesenharJogo(&jogo, estado)

		// ler input (bloqueia até tecla)
		ev := interfaceLerEventoTeclado()
		if !personagemExecutarAcao(ev, &jogo) {
			break
		}

		// pequeno delay para reduzir carga
		time.Sleep(25 * time.Millisecond)
	}

	// desconecta ao sair
	if jogo.Cliente != nil && jogo.ID != "" {
		var ack bool
		_ = jogo.Cliente.Call("Servidor.DesconectarJogador", jogo.ID, &ack)
	}
	interfaceLimparTela()
}
