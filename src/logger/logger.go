package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger

	lokiURL = "http://loki:3100/loki/api/v1/push"
	appName = "my-go-app"

	bufferMutex sync.Mutex
	logStreams  = make(map[string]*LokiStream)
)

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

type lokiWriter struct {
	level  string
	writer io.Writer
}

func (w *lokiWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	enqueueLog(w.level, msg)
	return w.writer.Write(p)
}

func Init() {
	// create log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de log: %v", err)
	}

	multi := io.MultiWriter(os.Stdout, file)

	Info = log.New(&lokiWriter{level: "info", writer: multi}, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(&lokiWriter{level: "warn", writer: multi}, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(&lokiWriter{level: "error", writer: multi}, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(&lokiWriter{level: "debug", writer: multi}, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)

	// goroutine que envia logs a cada 5s
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			flushLogs()
		}
	}()
}

func enqueueLog(level, msg string) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	level = strings.ToUpper(level)

	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	stream, ok := logStreams[level]
	if !ok {
		stream = &LokiStream{
			Stream: map[string]string{
				"app":   appName,
				"level": level,
			},
			Values: [][]string{},
		}
		logStreams[level] = stream
	}

	stream.Values = append(stream.Values, []string{ts, msg})
}

func flushLogs() {
	bufferMutex.Lock()
	if len(logStreams) == 0 {
		bufferMutex.Unlock()
		return
	}

	// copy streams
	streams := []LokiStream{}
	for _, s := range logStreams {
		streams = append(streams, *s)
	}
	logStreams = make(map[string]*LokiStream) // limpa buffer
	bufferMutex.Unlock()

	data, err := json.Marshal(LokiPayload{Streams: streams})
	if err != nil {
		log.Println("Falha ao serializar Loki payload:", err)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", lokiURL, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Falha ao criar requisição para Loki:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Falha ao enviar logs para Loki:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("Logs enviados com sucesso para o Loki!")
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erro ao enviar logs para Loki. Status: %d, Body: %s\n", resp.StatusCode, string(body))
	}
}

func InfoLog(msg string) {
	Info.Println(msg)
}

func WarnLog(msg string) {
	Warn.Println(msg)
}

func ErrorLog(msg string) {
	Error.Println(msg)
}

func DebugLog(msg string) {
	Debug.Println(msg)
}
