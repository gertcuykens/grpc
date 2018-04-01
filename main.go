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

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"
)

const U64 = int(unsafe.Sizeof(uint64(0)))

type taskServer struct{}

func Listen(server *taskServer) {
	srv := grpc.NewServer()
	RegisterTasksServer(srv, server)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("could not listen to :8080 %s", err)
	}
	log.Fatal(srv.Serve(l))
}

func Client() TasksClient {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %s\n", err)
		os.Exit(1)
	}
	return NewTasksClient(conn)
}

func (taskServer) Add(ctx context.Context, t *Task) (*Void, error) {
	return nil, add(t)
}

func (taskServer) List(ctx context.Context, l *TaskList) (*Void, error) {
	return nil, list(l)
}

//go:generate protoc -I . task.proto --go_out=plugins=grpc:.
func main() {
	server := taskServer{}
	go Listen(&server)
	client := Client()

	var l TaskList
	for true {
		fmt.Print("cmd: ")
		reader := bufio.NewReader(os.Stdin)
		switch cmd, err := reader.ReadString('\n'); {
		case strings.HasPrefix(cmd, "list"):
			// err = list(&l)
			_, err = client.List(context.Background(), &l)
			print(&l)
		case strings.HasPrefix(cmd, "add "):
			task := Task{
				Text: cmd[4 : len(cmd)-1],
				Done: false,
			}
			// err = add(&task)
			// err = list(&l)
			_, err = client.Add(context.Background(), &task)
			_, err = client.List(context.Background(), &l)
			print(&l)
		case strings.HasPrefix(cmd, "exit"):
			os.Exit(0)
		default:
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Fprintf(os.Stderr, "unknown subcommand %s (list/add/exit)", cmd)
			}
		}
	}

}

func add(task *Task) error {
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

func list(tasks *TaskList) error {
	b, err := ioutil.ReadFile("data")
	if err != nil {
		return fmt.Errorf("could not read: %s", err)
	}
	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < U64 {
			return fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l uint64
		if err := binary.Read(bytes.NewReader(b[:U64]), binary.LittleEndian, &l); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[U64:]

		var task Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]

		tasks.Tasks = append(tasks.Tasks, &task)
	}
	return err
}

func print(l *TaskList) {
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
