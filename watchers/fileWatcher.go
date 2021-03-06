package watchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/parvez0/disk-watcher/responses"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DfCMD struct {
	Filesystem string
	Size       string
	Used       string
	Available  string
	UsePer     int
	MountedOn  string
}

var logger = NewLogger()
var cache = make(map[string]int)
var watcherMasterUrl = os.Getenv("WATCHER_MASTER")

func CheckError(err interface{}, msg string) {
	if err != nil {
		logger.Panic(msg, err)
	}
}

func Greet() {
	logger.Info("Starting the watcher...")
}

func getDiskSpaceExceedMessage(accountName string, dirUsage *DfCMD, newValue string, message string) string {
	return fmt.Sprintf(
		`The whatsapp disk for account <b>%s</b> has exceeded <b>%d%%</b> of the total disk %s %s
	<table>
		<tr>
			<td>Mounted On: </td>
			<td>%s</td>
		</tr>
		<tr>
			<td>Size:</td>
			<td>%s</td>
		</tr>
		<tr>
			<td>Used: </td>
			<td>%s</td>
		</tr>
		<tr>
			<td>Used%%: </td>
			<td>%d%%</td>
		</tr>
		<tr>
			<td>Available: </td>
			<td>%s</td>
		</tr>
		<tr>
			<td>Increasing: </td>
			<td>%s</td>
		</tr>
	</table>`,
		accountName, dirUsage.UsePer, dirUsage.Size,
		message, dirUsage.MountedOn, dirUsage.Size, dirUsage.Used, dirUsage.UsePer, dirUsage.Available, newValue,
	)
}

func trim(values []string) []string {
	var formatted []string
	for _, v := range values {
		if v != "" {
			formatted = append(formatted, v)
		}
	}
	return formatted
}

func GetDirectoryUsage(diskToMonitor string, namespace string) *DfCMD {
	defer func() {
		if err := recover(); err != nil {
			logger.Info("recovered from panic disk storage for account ", namespace)
		}
	}()
	cmd := ExecPod(namespace)
	lines := strings.Split(cmd, "\n")
	var actualDisk string
	for _, v := range lines {
		if index := strings.Index(v, diskToMonitor); index != -1 {
			actualDisk = v
			break
		}
	}
	values := trim(strings.Split(actualDisk, " "))
	if len(values) < 4 {
		errMessage := "failed to get usage for directory " + diskToMonitor + " please verify if the directory exits"
		logger.Error(errMessage)
		panic(errMessage)
	}
	cs, err := strconv.Atoi(strings.Replace(values[4], "%", "", -1))
	CheckError(err, "failed to convert usage percent to int")
	res := DfCMD{
		Filesystem: values[0],
		Size:       values[1],
		Used:       values[2],
		Available:  values[3],
		MountedOn:  values[len(values)-1],
		UsePer:     cs,
	}
	logger.Debug("fetched directory usage - ", res)
	return &res
}

func IncreaseStorageSpace(namespace string) string {
	logger := NewLogger()
	reqBody := responses.StorageResizeReq{
		Namespace: namespace,
	}
	defer func() {
		if e := recover(); e != nil {
			logger.Error("failed to increase storage size, recovering from panic - ", e)
			return
		}
	}()
	buf, _ := json.Marshal(reqBody)
	resp, err := http.Post(watcherMasterUrl+"/increase-storage-size", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logger.Error("failed to scale, got api error - ", err)
		return "0"
	}
	defer resp.Body.Close()
	data := responses.GenericResponse{}
	buf, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &data)
	if err != nil {
		logger.Error("failed to unmarshall the response - ", err)
		return "0"
	}
	if data.Success == false {
		logger.Error("api returned failed response - ", data)
		return "0"
	}
	value := data.Data.(map[string]interface{})
	return value["currentSize"].(string)
}

func ProcessDiskUsageOutput(dirUsage *DfCMD, accountName string) {
	switch {
	case os.Getenv("GO_ENV") == "development":
		logger.Printf("whatsapp usage has exceeds 70%% sending a mail to notify - %+v", dirUsage)
		SendMail("whatsapp-disk", getDiskSpaceExceedMessage(accountName, dirUsage, "0", ""), nil)
	case dirUsage.UsePer < 70:
		logger.Printf("disk check passed current metrics - %+v", dirUsage)
	case dirUsage.UsePer > 70 && dirUsage.UsePer < 80:
		logger.Printf("whatsapp usage has exceeds 70%% sending a mail to notify - %+v", dirUsage)
		if cache[accountName] == 0 {
			SendMail("whatsapp-disk", getDiskSpaceExceedMessage(accountName, dirUsage, "0", ""), nil)
			cache[accountName] = 1
		}
	case dirUsage.UsePer > 80:
		logger.Warnf("whatsapp usage has exceeds 80%% increasing 100GB  - %+v", dirUsage)
		newValue := IncreaseWhatsappDiskSize(accountName)
		value := strconv.FormatInt(newValue, 10)
		message := ", increasing the volume to " + value + "GB"
		SendMail("whatsapp-disk", getDiskSpaceExceedMessage(accountName, dirUsage, value, message), nil)
		cache[accountName] = 0
	}
}

func CheckDiskStorage(dirToMonitor *string) {
	namespaces, err := NamespaceList()
	if err != nil {
		logrus.Panic("failed to list namespaces -", err)
	}
	for _, ns := range namespaces.Items {
		if strings.HasPrefix(ns.Name, "wa-") {
			logrus.Infof("checking the storage of account - %+v", ns.Name)
			usage := GetDirectoryUsage(*dirToMonitor, ns.Name)
			if usage != nil {
				ProcessDiskUsageOutput(usage, ns.Name)
			} else {
				SendMail("whatsapp-disk", "Failed to check the available disk space for account "+ns.Namespace, nil)
			}
		}
	}
}
