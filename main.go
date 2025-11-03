// main.go - cliente: loop principal
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

	// iniciar interface
	interfaceIniciar()
	defer interfaceFinalizar()

	go func() {
		for !jogo.encerrar {
			var estado EstadoJogo
			err := client.Call("Servidor.GetEstadoJogo", id, &estado)
			if err == nil {
				interfaceDesenharJogo(&jogo, estado)
			}
			time.Sleep(100 * time.Millisecond) // 100ms
		}
	}()

	// loop principal
	for !jogo.encerrar {
		ev := interfaceLerEventoTeclado()
		if !personagemExecutarAcao(ev, &jogo) {
			break
		}
	}

	// desconecta ao sair
	if jogo.Cliente != nil && jogo.ID != "" {
		var ack bool
		_ = jogo.Cliente.Call("Servidor.DesconectarJogador", jogo.ID, &ack)
	}
	interfaceLimparTela()
}
