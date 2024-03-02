package main

import (
	"context"
	"flag"
	"log"
	"net"

	desc "github.com/KozlovNikolai/chat-server/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	//"honnef.co/go/tools/config"
	"github.com/KozlovNikolai/chat-server/internal/config"
	"github.com/KozlovNikolai/chat-server/internal/config/env"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

//const grpcPort = 50051

type server struct {
	desc.UnimplementedChat_V1Server
}

// Create is receivig list of users, creating new chat and responsing ID of new chat
func (s *server) Create(_ context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Received Usernames: %#v\n", in.Usernames)
	return &desc.CreateResponse{Id: int64(len(in.Usernames))}, nil
}

// Delete is receiving ID of chat, deleting it, and nothing respons
func (s *server) Delete(_ context.Context, in *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Received Id: %v\n", in.Id)
	return &empty.Empty{}, nil
}

// SendMessage is receiving Sender, His text and Timestamp, and nothing response
func (s *server) SendMessage(_ context.Context, in *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Received From: %v\n", in.From)
	log.Printf("Received Text: %v\n", in.Text)
	log.Printf("Received Time: %v\n", in.Timestamp)
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	//ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	// pgConfig, err := env.NewPGConfig()
	// if err != nil {
	// 	log.Fatalf("failed to get pg config: %v", err)
	// }

	lis, err := net.Listen("tcp", grpcConfig.Address()) //открываем порт для прослушивания
	if err != nil {                                     //проверяем, что порт открылся, иначе
		log.Fatalf("failed to listen %v", err) //выводим в лог ошибку и закрываем программу
	}
	s := grpc.NewServer()                    //создаем  grpc сервер
	reflection.Register(s)                   //включаем возможность сервера рассказать о своих методах клиенту
	desc.RegisterChat_V1Server(s, &server{}) //регистрируем наши методы в grpc сервере
	log.Printf("server listerning at %v", lis.Addr())

	if err = s.Serve(lis); err != nil { //вызываем метод grpc сервера "s" - Serve, для начала работы на порту "lis"
		log.Fatalf("failed to serv %v", err)
	}

}
