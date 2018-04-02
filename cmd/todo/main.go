package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"unsafe"

	pb "github.com/gertcuykens/grpc"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
)

const U64 = int(unsafe.Sizeof(uint64(0)))

func Listen(server pb.TodoServer) {
	srv := grpc.NewServer()
	pb.RegisterTodoServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("could not listen to :8080 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() pb.TodoClient {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return pb.NewTodoClient(conn)
}

type Todo struct{}

func (Todo) Add(ctx context.Context, t *pb.Task) (*pb.Void, error) {
	return &pb.Void{}, add(t)
}

func (Todo) List(ctx context.Context, v *pb.Void) (*pb.TaskList, error) {
	return list()
}

//go:generate protoc -I ../../ ../../grpc.proto --go_out=plugins=grpc:../../
//go:generate mockgen -destination ../../mock_grpc/grpc.go github.com/gertcuykens/grpc TodoServer,TodoClient
func main() {
	go Listen(&Todo{})
	client := Client()

	var tasks *pb.TaskList
	for true {
		fmt.Print("cmd: ")
		reader := bufio.NewReader(os.Stdin)
		switch cmd, err := reader.ReadString('\n'); {
		case strings.HasPrefix(cmd, "list"):
			tasks, err = client.List(context.Background(), &pb.Void{})
			echo(tasks)
		case strings.HasPrefix(cmd, "add "):
			task := pb.Task{
				Text: cmd[4 : len(cmd)-1],
				Done: false,
			}
			_, err = client.Add(context.Background(), &task)
			tasks, err = client.List(context.Background(), &pb.Void{})
			echo(tasks)
		case strings.HasPrefix(cmd, "exit"):
			os.Exit(0)
		default:
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprint(os.Stderr, "unknown subcommand use (list/add/exit)\n")
			}
		}
	}

}

func add(task *pb.Task) error {
	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile("data", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s", err)
	}

	if err := binary.Write(f, binary.LittleEndian, uint64(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %s", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %s", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s", err)
	}
	return nil
}

func list() (*pb.TaskList, error) {
	b, err := ioutil.ReadFile("data")
	if err != nil {
		return nil, fmt.Errorf("could not read: %s", err)
	}

	var tasks pb.TaskList
	for {
		if len(b) == 0 {
			break
		} else if len(b) < U64 {
			return nil, fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l uint64
		if err := binary.Read(bytes.NewReader(b[:U64]), binary.LittleEndian, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[U64:]

		var task pb.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]

		tasks.Tasks = append(tasks.Tasks, &task)
	}
	return &tasks, err
}

func echo(l *pb.TaskList) {
	if l == nil {
		fmt.Println("🤬")
		return
	}
	for _, t := range l.Tasks {
		if t.Done {
			fmt.Printf("👍 %s\n", t.Text)
		} else {
			fmt.Printf("😱 %s\n", t.Text)
		}
	}
}
