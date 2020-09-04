package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/parvez0/disk-watcher/responses"
	"github.com/parvez0/disk-watcher/watchers"
	"github.com/robfig/cron/v3"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

//if req.URL.Path != "/" {
//http.NotFound(w, req)
//return
//}

func RecoverFromPanic(message string, writer *http.ResponseWriter) {
	logger := watchers.NewLogger()
	if err := recover(); err != nil {
		logger.Error(message, "error -", err)
		resp := map[string]interface{}{"message": "encountered an error, while processing this request", "error": fmt.Sprint(err)}
		responses.ResponseWithFailedMessage(http.StatusInternalServerError, resp, writer)
	}
}

func HealthCheckHandler(writer http.ResponseWriter, req *http.Request) {
	logger := watchers.NewLogger()
	logger.Info("received a health check request from ", req.RemoteAddr)
	resp := responses.GenericResponse{Success: true, Data: map[string]string{"message": "master server is working !!!"}}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(resp)
}

func ScaleDiskStorageHandler(writer http.ResponseWriter, req *http.Request) {
	body := responses.StorageResizeReq{}
	defer RecoverFromPanic("failed to increase storage size for account "+body.Namespace, &writer)
	json.NewDecoder(req.Body).Decode(&body)
	if body.Namespace == "" {
		resp := map[string]string{"message": "required parameter namespace not provided"}
		responses.ResponseWithFailedMessage(http.StatusBadRequest, resp, &writer)
		return
	}
	newValue := watchers.IncreaseWhatsappDiskSize(body.Namespace)
	resp := responses.GenericResponse{Success: true, Data: map[string]string{"currentSize": strconv.FormatInt(newValue, 10)}}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(resp)
}

func main() {
	workerType := flag.String("worker", "disk-watcher", "type of worker, defaults to master server")
	dirToMonitor := flag.String("dir", "wamedia", "directory to monitor for storage consumption")

	flag.Parse()

	logger := watchers.NewLogger()

	switch *workerType {
	case "disk-watcher":
		logger.Info("starting disk watcher for dir -", *dirToMonitor)
		c := cron.New()
		c.AddFunc("*/5 * * * *", func() {
			watchers.CheckDiskStorage(dirToMonitor)
		})
		go c.Start()
		logger.Printf("cron output - %+v", c.Entries())
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		c.Stop()
	case "master":
		logger.Info("starting master server on port 5000")
		http.HandleFunc("/health-check", HealthCheckHandler)
		http.HandleFunc("/increase-storage-size", ScaleDiskStorageHandler)
		http.ListenAndServe(":5000", nil)
		logger.Info("server is listening on port 5000")
	}
}
