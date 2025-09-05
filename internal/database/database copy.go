/*

   Em um sistema em produção, geralmente é preferível manter uma conexão persistente
   com o banco de dados, em vez de abrir e fechar conexões para cada consulta. Abrir
   uma conexão com um banco de dados é uma operação relativamente cara em termos de
   tempo e recursos. Portanto, se você estiver fazendo muitas consultas, o custo de
   abrir e fechar a conexão repetidamente pode se somar.

   No entanto, manter uma conexão aberta indefinidamente também tem suas desvantagens.
   Por exemplo, se o seu aplicativo mantém muitas conexões abertas simultaneamente,
   isso pode sobrecarregar o banco de dados. Além disso, se a conexão for interrompida
   por algum motivo (por exemplo, se o banco de dados cair ou a rede falhar), seu
   aplicativo precisará ser capaz de lidar com isso.

   Uma abordagem comum é usar um pool de conexões. Um pool de conexões mantém um número
   de conexões abertas e reutiliza-as conforme necessário. Quando uma consulta é feita,
   uma conexão é retirada do pool, usada e depois retornada ao pool. Isso oferece um bom
   equilíbrio entre eficiência (porque você não precisa abrir uma nova conexão para cada
   consulta) e uso de recursos (porque você limita o número de conexões abertas).

   Neste código, SetMaxOpenConns(25) define o número máximo de conexões abertas para
   o banco de dados, SetMaxIdleConns(25) define o número máximo de conexões ociosas
   que podem existir simultaneamente e SetConnMaxLifetime(5 * time.Minute) define a
   duração máxima que uma conexão pode ser reutilizada.
*/

//SQL, err = sql.Open("pgx", "postgres://docker:docker@localhost:5432/app")

package database

/*import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	//_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

func loadDBConnectionString() (string, error) {
    err := godotenv.Load()
	if err != nil {
        return "", fmt.Errorf(".env file not found - Could not load environment variables")
    }

    raw := os.Getenv("DATABASE_URL")
    if raw == "" {
        return "", fmt.Errorf("DATABASE_URL not set")
    }

    // Expande variáveis internas ($PGUSER etc.) se houver
    dsn := os.ExpandEnv(raw)

    return dsn, nil
}

func newConnection() (*sql.DB, error) {
	dsn, err := loadDBConnectionString()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("Unable to connect to database: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

	err = db.PingContext(ctx)
    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    return db, nil
}

func Disconnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}

func Connection() (*sql.DB, error){
	db, err := newConnection()
    if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
    }

    fmt.Println("Successfully connected!")

	return db, nil
}

func GetDBInfo() error {
    db, err := Connection()
    if err != nil {
        return err
    }
    defer Disconnect(db)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var (
        serverVersion    string
        maxConnections   int
        openedConnections int
    )

    err = db.QueryRowContext(
        ctx,
        `SELECT 
            current_setting('server_version') AS server_version, 
            current_setting('max_connections')::int AS max_connections, 
            (SELECT COUNT(*)::int FROM pg_stat_activity WHERE datname = $1) AS opened_connections;`,
        os.Getenv("PGDATABASE"), // valor para $1
    ).Scan(&serverVersion, &maxConnections, &openedConnections)

    if err != nil {
        return fmt.Errorf("failed to get database info: %v", err)
    }

    fmt.Printf("Versão: %s \nMáx conexões: %d \nConexões abertas: %d\n",
        serverVersion, maxConnections, openedConnections)

    return nil
}*/

//INJEÇÃO DE DEPENDENCIAS
/*Entendi o que você quer dizer — se cada *model* ou função precisar receber `*sql.DB` como parâmetro, o código começa a ficar “poluído” com esse parâmetro em todo lugar.  
Desenvolvedores experientes evitam isso usando **injeção de dependência** e **structs que encapsulam a conexão**. Assim, você não precisa ficar passando `*sql.DB` manualmente em cada chamada.

---

## 🔹 Padrão usado por sêniores: Repositório/Service com conexão embutida

A ideia é criar um tipo que **guarda a conexão** e expõe métodos para trabalhar com o banco.  
Assim, a conexão é criada uma vez e usada em todos os métodos sem precisar passar como argumento.

---

### Exemplo

```go
package models

import "database/sql"

type UserModel struct {
    DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
    return &UserModel{DB: db}
}

func (m *UserModel) CountUsers() (int, error) {
    var count int
    err := m.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    return count, err
}
```

---

### Uso no `main.go`

```go
package main

import (
    "fmt"
    "log"
    "meuprojeto/database"
    "meuprojeto/models"
)

func main() {
    db, err := database.NewConnection()
    if err != nil {
        log.Fatalf("Erro ao conectar: %v", err)
    }
    defer db.Close()

    userModel := models.NewUserModel(db)

    count, err := userModel.CountUsers()
    if err != nil {
        log.Fatalf("Erro ao contar usuários: %v", err)
    }

    fmt.Println("Total de usuários:", count)
}
```

---

## 🔹 Vantagens dessa abordagem

- **Sem poluição de parâmetros**: você não precisa passar `db` em cada função.
- **Organização**: cada *model* ou *repository* cuida de uma tabela ou conjunto de queries.
- **Testabilidade**: fácil trocar `*sql.DB` por um mock em testes.
- **Reuso**: a mesma instância do model pode ser usada em vários lugares.

---

💡 **Resumo**:  
Você ainda reaproveita a conexão (o `*sql.DB` é criado uma vez só), mas encapsula dentro de structs específicas para cada domínio do seu sistema.  
Isso é o que a maioria dos sêniores faz para evitar passar a conexão manualmente por todo o código.

---

Se você quiser, posso te montar **um blueprint completo** com:
- Inicialização única do banco
- Models organizados por domínio
- Injeção de dependência limpa  
Assim seu projeto já fica pronto para crescer sem bagunça. Quer que eu monte?

*/

/*func GetDBInfo(db *sql.DB) (map[string]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    info := make(map[string]string)

    // Versão do PostgreSQL
    var version string
    if err := db.QueryRowContext(ctx, "SELECT version()").Scan(&version); err != nil {
        return nil, err
    }
    info["version"] = version

    // Nome do banco atual
    var dbName string
    if err := db.QueryRowContext(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
        return nil, err
    }
    info["database"] = dbName

    // Usuário conectado
    var user string
    if err := db.QueryRowContext(ctx, "SELECT current_user").Scan(&user); err != nil {
        return nil, err
    }
    info["user"] = user

    return info, nil
} */

/*func createTable() {
	_, err := SQL.Exec("CREATE TABLE IF NOT EXISTS \"user\" (id SERIAL PRIMARY KEY, email VARCHAR(100) UNIQUE NOT NULL, password VARCHAR(100) NOT NULL)")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}*/


/*func NewConnection() (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    return db, nil
}*/


/*package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	//_ "github.com/jackc/pgx/v4/stdlib"
)*/

/*var SQL *sql.DB

func connectDatabase(driver, dataSourceName string) error {
	var err error
	SQL, err = sql.Open(driver, dataSourceName)
	if err != nil {
		return err
	}

	// Configurar o pool de conexões
	SQL.SetMaxOpenConns(25)
	SQL.SetMaxIdleConns(25)
	SQL.SetConnMaxLifetime(5 * time.Minute)

	// Verificar a conexão
	err = SQL.Ping()
	if err != nil {
		return err
	}

	return nil
}*/

/*func Disconnect() {
	err := SQL.Close()
	if err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}*/



/*func Connection() {

	URLConnection, err := loadDBConnectionString()
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := newConnection()
    if err != nil {
        log.Fatalf("Falha ao inicializar banco: %v", err)
    }
    defer db.Close()

    fmt.Println("Successfully connected!")

	//fmt.Println(URLConnection)
	err = connectDatabase("pgx", URLConnection)
	//err = connectDatabase("pgx", "postgres://docker:docker@localhost:5432/app")
	//err := connectDatabase("pgx", "postgres://postgres:postgres@localhost:5432/app")
	//err := connectDatabase("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	createTable()

	fmt.Println("Successfully connected!")
}*/
