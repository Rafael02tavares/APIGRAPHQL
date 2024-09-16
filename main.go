package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// Estrutura para armazenar usuários temporariamente (simulando um banco de dados na memória)
var users = map[graphql.ID]*UserResolver{}
var nextID int = 1

// Resolver principal
type Resolver struct{}

// Resolver para a query `hello`
func (r *Resolver) Hello() string {
	return "Hello, GraphQL!"
}

// Resolver para a query `user`
func (r *Resolver) User(args struct{ ID graphql.ID }) *UserResolver {
	if user, ok := users[args.ID]; ok {
		return user
	}
	return nil
}

// Resolver para a mutation `createUser`
func (r *Resolver) CreateUser(args struct {
	Name  string
	Email string
}) *UserResolver {
	id := graphql.ID(string(nextID))
	user := &UserResolver{id: id, name: args.Name, email: args.Email}
	users[id] = user
	nextID++
	return user
}

// Estrutura para resolver o tipo `User`
type UserResolver struct {
	id    graphql.ID
	name  string
	email string
}

func (u *UserResolver) ID() graphql.ID {
	return u.id
}

func (u *UserResolver) Name() string {
	return u.name
}

func (u *UserResolver) Email() string {
	return u.email
}

func main() {
	// Ler o esquema GraphQL do arquivo
	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		log.Fatal("Erro ao ler o arquivo schema.graphql: ", err)
	}

	// Criar o esquema GraphQL
	schema := graphql.MustParseSchema(string(schemaFile), &Resolver{})

	// Configurar o servidor HTTP
	http.Handle("/graphql", &relay.Handler{Schema: schema})

	log.Println("Servidor rodando em http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
