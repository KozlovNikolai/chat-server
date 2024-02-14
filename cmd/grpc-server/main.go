package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/KozlovNikolai/chat-server/pkg/chat_v1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051 //порт, который слушает grpc сервер

type server struct { //структура которая через алиас "desc" имплементирует API (все методы) grpc сервера
	desc.UnimplementedChat_V1Server
}

// первый entrypoint (точка входа)
func (s *server) Create(_ context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Received Usernames: %#v\n", in.Usernames)          //выводим, что приняли
	return &desc.CreateResponse{Id: int64(len(in.Usernames))}, nil //возвращаем количество элементов слайса (что бы нескучно было)
}

// второй entrypoint (точка входа)
func (s *server) Delete(_ context.Context, in *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Received Id: %v\n", in.Id) //выводим, что приняли
	return &empty.Empty{}, nil             //ничего не возвращаем
}

// третий entrypoint (точка входа)
func (s *server) SendMessage(_ context.Context, in *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Received From: %v\n", in.From) //выводим, что приняли
	log.Printf("Received Text: %v\n", in.Text)
	log.Printf("Received Time: %v\n", in.Timestamp)
	return &empty.Empty{}, nil //ничего не возвращаем
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort)) //открываем порт для прослушивания
	if err != nil {                                             //проверяем, что порт открылся, иначе
		log.Fatalf("failed to listen %v", err) //выводим в лог ошибку и закрываем программу
	}
	s := grpc.NewServer()                    //создаем grpc сервер
	reflection.Register(s)                   //включаем возможность сервера рассказать о своих методах клиенту
	desc.RegisterChat_V1Server(s, &server{}) //регистрируем наши методы в grpc сервере
	log.Printf("server listerning at %v", lis.Addr())

	if err = s.Serve(lis); err != nil { //вызываем метод grpc сервера "s" - Serve, для начала работы на порту "lis"
		log.Fatalf("failed to serv %v", err)
	}

}
