package task

import "fmt"
import "github.com/maxpowel/wiphonego/fetcher/masmovil"
import (
	"github.com/maxpowel/wiphonego"
	"github.com/maxpowel/dislet"
	"github.com/RichardKnop/machinery/v1"
	"log"
	"github.com/jinzhu/gorm"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/maxpowel/wiphonego/protomodel"
	"encoding/base64"


	"github.com/maxpowel/dislet/apirest"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var kernel *dislet.Kernel

func GetConsumptionTask(phoneLineId uint) (string, map[string]string, error){
	db := kernel.Container.MustGet("database").(*gorm.DB)

	phoneLine := &wiphonego.PhoneLine{}
	if db.Preload("Device").Preload("Credentials.Operator").Preload("Credentials").First(phoneLine, phoneLineId).RecordNotFound() {
		params := make(map[string]string)
		params["line"] = string(phoneLineId)
		return "Line {line} not found", params, nil
	}

	credentials := &phoneLine.Credentials
	if credentials.Operator.InternalName == "masmovil" {
		mv := masmovil.NewFetcher(credentials, kernel)
		consumption, err := mv.GetInternetConsumption(phoneLine.PhoneNumber)
		if err != nil {
			return "", nil, err
		}
		consumption.PhoneLineID = phoneLine.ID
		db.Create(&consumption)


		commandData, _ := proto.Marshal(&protomodel.NotificationConsumptionCommand{
			Consumption: []*protomodel.ConsumptionResponse{
				{
					PhoneNumber: phoneLine.PhoneNumber,
					PeriodEnd: int32(consumption.PeriodEnd.Unix()),
					PeriodStart: int32(consumption.PeriodStart.Unix()),
					UpdatedAt: int32(consumption.UpdatedAt.Unix()),
					CallConsumed: int32(consumption.CallConsumed),
					CallTotal: int32(consumption.CallTotal),
					InternetTotal: consumption.InternetTotal,
					InternetConsumed: consumption.InternetConsumed,
				},
			},
		})

		jsonValue, _ := json.Marshal(Notification{
			Notifications: []NotificationEntry{
				{
					Tokens: []string{phoneLine.Device.Uuid},
					Platform:2,
					Data: map[string]interface{} {
						"type": protomodel.NotificationCommandType_CONSUMPTION,
						"_": base64.StdEncoding.EncodeToString(commandData),
					},
				},
			},
		})
		fmt.Println(string(jsonValue))

		//resp, err := http.NewRequest("POST", "http://localhost:8088/api/push", bytes.NewBuffer(jsonValue))
		req, _ := http.NewRequest("POST", "http://localhost:8088/api/push", bytes.NewBuffer(jsonValue))
		client := &http.Client{}
		client.Do(req)


		params := make(map[string]string)
		params["id"] = fmt.Sprint(consumption.ID)
		return "Consumption created with id {id}", params, nil
	}




	params := make(map[string]string)
	params["operator"] = phoneLine.Credentials.Operator.InternalName
	return "Operator {operator} is not available", params, nil

}

type NotificationEntry struct {
	Tokens []string `json:"tokens"`
	Platform int `json:"platform"`
	Data map[string]interface{} `json:"data"`
}

type Notification struct {
	Notifications []NotificationEntry `json:"notifications"`
}
/*func GetAnonymousConsumptionTask(username, password, operator, phoneNumber, deviceId string) (string, map[string]string, error){

	db := kernel.Container.MustGet("database").(*gorm.DB)

	if operator == "masmovil" {
		mv := masmovil.NewFetcher(&wiphonego.Credentials{Username: username, Password: password}, kernel)
		c, err := mv.GetInternetConsumption(phoneNumber)
		c.PeriodEnd = now.EndOfMonth()
		c.PeriodStart = now.BeginningOfMonth()
		c.PhoneNumber = phoneNumber

		if err != nil {
			params := make(map[string]string)
			return "Invalid credentials", params, nil
		} else {
			device := wiphonego.UserDevice{}
			db.Where("Uuid = ?", deviceId).FirstOrCreate(&device, wiphonego.UserDevice{Uuid: deviceId})
			//device.Uuid = device
			c.Device = device
			db.Create(&c)
			return "", nil, err
		}
	}
	params := make(map[string]string)
	params["operator"] = operator
	return "Operator {operator} is not available", params, nil

}*/
func ConsumptionSignature (phoneLineId uint) (*tasks.Signature){
	return &tasks.Signature{
		Name: "consumption",
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: phoneLineId,
			},
		},
	}
}

func AnonymousCredentialsTask(username, password, operator, deviceId string) (string, map[string]string, error){

	db := kernel.Container.MustGet("database").(*gorm.DB)
	if operator == "masmovil" {
		lineOperator := &wiphonego.Operator{}
		db.Where(&wiphonego.Operator{InternalName: operator}).First(lineOperator)

		credentials := &wiphonego.Credentials{Username: username, Password: password, Operator: *lineOperator}

		mv := masmovil.NewFetcher(credentials, kernel)
		lines, err := mv.GetLines()
		if err != nil {
			params := make(map[string]string)
			return err.Error(), params, nil
		} else {
			fmt.Println("INSERTAR LNEAS")
			fmt.Println(lines)
			device := wiphonego.Device{}
			db.Where("Uuid = ?", deviceId).FirstOrCreate(&device, wiphonego.Device{Uuid: deviceId})
			db.Create(credentials)
			//device.Uuid = device


			//
			for _, element := range lines {
				//line := wiphonego.PhoneLine{Device: device, PhoneNumber: element,ID: uuid.NewV4().String()}
				line := wiphonego.PhoneLine{Device: device, PhoneNumber: element, Credentials: *credentials}
				fmt.Println(line)
				db.Create(&line)
				apirest.SendTask(kernel, ConsumptionSignature(line.ID))
			}
			params := make(map[string]string)
			// Pack the lines into the response
			protoLines, _ := proto.Marshal(&protomodel.AnonymousCredentialsResponse{CredentialsId: int32(credentials.ID), PhoneNumbers: lines})
			params["lines"] = base64.StdEncoding.EncodeToString(protoLines)
			return "", params, err
		}
	}
	params := make(map[string]string)
	params["operator"] = operator
	return "Operator {operator} is not available", params, nil

}

func Bootstrap(k *dislet.Kernel) {
	kernel = k
	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		server := k.Container.MustGet("machinery").(*machinery.Server)
		// Register tasks

		err := server.RegisterTasks(map[string]interface{}{
			"consumption": GetConsumptionTask,
			//"anonymousConsumption": GetAnonymousConsumptionTask,
			"anonymousCredentials": AnonymousCredentialsTask,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	k.Subscribe(baz)

}
