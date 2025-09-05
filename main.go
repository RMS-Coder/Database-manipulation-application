package main

import (
	//"fmt"
	"net/http"
	//"time"

	"github.com/RodrigoMS/app/internal/database"
)

func main() {

	database.OpenConnection()
	database.CheckConnection()
	database.GetDB().GetDBInfo()
	/*go func() {
		for {
			// Primeiro verifica se já existe conexão ativa
			if err := database.CheckConnection(); err != nil {
				fmt.Println("\033[31mBanco offline. Tentando reconectar...\033[0m", err)

				time.Sleep(10 * time.Second)
				continue
				
				fmt.Println("\033[32mReconexão bem-sucedida!\033[0m")
			} else {
				fmt.Println("\033[32mConexão com o banco está ativa.\033[0m")
			}

			// Aguarda antes de verificar novamente
			time.Sleep(10 * time.Second)
		}
	}()*/

	// Toda vez que for usar o banco, chama checkConnection()
    /*if err := database.CheckConnection(); err != nil {
        fmt.Println("\033[ Connection error:\033[0m", err)
        return
    }*/

    defer database.CloseConnection()

	routes()

	http.ListenAndServe(":8080", nil)
}
