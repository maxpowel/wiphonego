package controller

import (
	"net/http"
	"github.com/RichardKnop/machinery/v1"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RangelReale/osin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"github.com/maxpowel/dislet"
	"github.com/maxpowel/dislet/usermngr"
	"github.com/maxpowel/dislet/apirest"
	"github.com/maxpowel/wiphonego/protomodel"
	"github.com/maxpowel/wiphonego"
)

type CredentialsValidator struct {
	Username  string `validate:"required"`
	Password string `validate:"required"`
	Operator string `validate:"required"`
}


func Me(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request, user *usermngr.User) (proto.Message, error){

	var ts proto.Message = &protomodel.UseLoginResponse{
		Token: "333",
	}

	return ts, nil
	//return ts, apirest.StatusError{Code:401, Err: fmt.Errorf("LOL")}
}

func Index(w http.ResponseWriter, r *http.Request) {

	u := usermngr.NewUser()
	u.Username = "pepe"
	usermngr.PlainPassword(u, "123456")
	fmt.Fprintln(w, "Welcome!")
	fmt.Fprintln(w, u.Password)
	fmt.Fprintln(w, u.Salt)
}

func Login(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error){
	loginData := &protomodel.UserLogin{}
	err := apirest.GetBody(loginData, r)
	if err != nil {
		return nil, apirest.StatusError{Code:401, Err: err}
	}

	userManager := kernel.Container.MustGet("user_manager").(*usermngr.Manager)
	user, err := userManager.FindUser(loginData.Username)
	if err != nil {
		return nil, apirest.StatusError{Code:404, Err: err}
	}

	if usermngr.CheckPassword(user, loginData.Password) == nil {
		token, err := apirest.GenerateToken(kernel, user)
		if err != nil {
			return nil, apirest.StatusError{Code:500, Err: err}
		} else {
			ts := &protomodel.UseLoginResponse{
				Token: token,
			}

			/*data, err := proto.Marshal(&ts)
			if err != nil {
				return apirest.StatusError{Code:500, Err: err}
			}

			w.Write(data)*/
			return ts, nil
		}
		//w.Write()
	} else {
		return nil, apirest.StatusError{Code:401, Err: fmt.Errorf("Invalid password")}
	}

	return nil, nil
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo Index!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func Index2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("COSA","ALGO")
	fmt.Println(r.Header.Get("token"))
	w.WriteHeader(34)

	fmt.Fprintln(w, "Welcome!")

}


type AnonymousConsumptionValidator struct {
	DeviceId  string `validate:"required"`
	Credentials CredentialsValidator `validate:"required"`
	PhoneNumber  string `validate:"required"`
}

func GetAnonymousConsumption(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error) {
	requestData := &protomodel.AnonymousConsumptionRequest{}
	err := apirest.GetBody(requestData, r)
	if err != nil {
		return nil, apirest.StatusError{401, err}
	}
	//fmt.Println(requestData.DeviceId)
	//requestData.DeviceId = "lolazo"
	/*requestData.Credentials = &protomodel.Credentials{}
	requestData.Credentials.Operator = "masmovil"
	requestData.Credentials.Username = "222alvaro_gg@hotmail.com"
	requestData.Credentials.Password = "22MBAR4B1"
	requestData.PhoneNumber = "677077536"
	requestData.DeviceId = "5656"*/

	k := &AnonymousConsumptionValidator{}
	fmt.Println(requestData)
	if requestData.Credentials != nil {
		k.Credentials = CredentialsValidator{
			Username: requestData.Credentials.Username,
			Password: requestData.Credentials.Password,
			Operator: requestData.Credentials.Operator,
		}
	}
	_, err = apirest.Validate(requestData, k)

	if err != nil {
		return nil, apirest.StatusError{402, err}
	}

	response, err := apirest.SendTask(kernel, anonymousConsumptionSignature(requestData))
	if err != nil {
		return nil, apirest.StatusError{http.StatusInternalServerError, err}
	}

	return response, nil
}

func GetLastAnonymousConsumption(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	database := kernel.Container.MustGet("database").(*gorm.DB)
	consumptionData := wiphonego.UserDeviceConsumption{}

	database.Joins("left join user_devices on user_devices.id = user_device_consumptions.device_id").Where("user_devices.uuid = ?", deviceId).Last(&consumptionData)

	response := &protomodel.ConsumptionResponse{}
	response.CallConsumed = int32(consumptionData.CallConsumed)
	response.CallTotal = int32(consumptionData.CallTotal)
	response.InternetConsumed = consumptionData.InternetConsumed
	response.InternetTotal = consumptionData.InternetTotal
	response.PeriodStart = int32(consumptionData.PeriodStart.Unix())
	response.PeriodEnd = int32(consumptionData.PeriodEnd.Unix())
	response.UpdatedAt = int32(consumptionData.UpdatedAt.Unix())
	response.PhoneNumber = consumptionData.PhoneNumber

	if consumptionData.ID == 0 {
		return nil, apirest.StatusError{400, fmt.Errorf("Uuid not found")}
	}

	return response, nil
}

func consumptionSignature (username, password, operator string) (*tasks.Signature){
	return &tasks.Signature{
		Name: "consumption",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: username,
			},
			{
				Type:  "string",
				Value: password,
			},
			{
				Type:  "string",
				Value: operator,
			},
		},
	}
}

func anonymousConsumptionSignature (data *protomodel.AnonymousConsumptionRequest) (*tasks.Signature){
	return &tasks.Signature{
		Name: "anonymousConsumption",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: data.Credentials.Username,
			},
			{
				Type:  "string",
				Value: data.Credentials.Password,
			},
			{
				Type:  "string",
				Value: data.Credentials.Operator,
			},
			{
				Type:  "string",
				Value: data.PhoneNumber,
			},
			{
				Type:  "string",
				Value: data.DeviceId,
			},
		},
	}
}



func Bootstrap(k *dislet.Kernel) {
	var baz dislet.OnKernelReady = func(k *dislet.Kernel){
		router := k.Container.MustGet("api").(*mux.Router)
		router.HandleFunc("/", Index)
		router.Handle("/login", apirest.Handler{k, Login})
		//router.Handle("/me", apirest.Handler{k, Me})
		router.Handle("/me", apirest.SecureHandler{k, Me})
		router.HandleFunc("/todos", TodoIndex)
		router.HandleFunc("/todos/{todoId}", TodoShow)
		router.Methods("PUT").Path("/este").Name("este").HandlerFunc(Index2)

		router.Handle("/tarea", apirest.Handler{k, GetIndex})
		router.Handle("/anonymousConsumption", apirest.Handler{k, GetAnonymousConsumption}).Methods("POST")
		router.Handle("/anonymousConsumption/{deviceId}", apirest.Handler{k, GetLastAnonymousConsumption})
		router.Handle("/consumption", apirest.Handler{k, GetConsumption})
		router.Handle("/consumption/{taskUid}", apirest.Handler{k, GetTaskState})
		router.Handle("/task/{taskUid}", apirest.Handler{k, GetTaskState})

	}

	k.Subscribe(baz)

}


func GetIndex(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error) {

	fmt.Println("EL OTRO HILO")

	// Enviar tarea
	task0 := tasks.Signature{
		Name: "add",
		Args: []tasks.Arg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 1,
			},
		},
	}

	fmt.Println("Enviando task...")
	server := kernel.Container.MustGet("machinery").(*machinery.Server)
	asyncResult, err := server.SendTask(&task0)

	if err != nil {
		// We return a status error here, which conveniently wraps the error
		// returned from our DB queries. We can clearly define which errors
		// are worth raising a HTTP 500 over vs. which might just be a HTTP
		// 404, 403 or 401 (as appropriate). It's also clear where our
		// handler should stop processing by returning early.
		return nil, apirest.StatusError{500, err}
	}

	return &protomodel.UseLoginResponse{Token: asyncResult.Signature.UUID}, nil
	//w.Write([]byte(asyncResult.Signature.UUID))

	/*results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		fmt.Println("Getting task result failed with error: %s", err.Error())
	}
	fmt.Printf(
		"%v + %v = %v\n",
		asyncResult.Signature.Args[0].Value,
		asyncResult.Signature.Args[1].Value,
		results[0].Interface(),
	)*/


	//return nil
}



func CheckToken(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) error {
	server := kernel.Container.MustGet("oauth").(*osin.Server)
	database := kernel.Container.MustGet("database").(*gorm.DB)

	resp := server.NewResponse()
	defer resp.Close()
	fmt.Println("CERO")
	var err error
	if ar := server.HandleAccessRequest(resp, r); ar != nil {
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		user := usermngr.User{}
		database.Where("username = ?", username).First(&user)
		err = usermngr.CheckPassword(&user, password)
		ar.Authorized = err == nil
		if ar.Authorized {
			ar.UserData = user.ID
		}
		server.FinishAccessRequest(resp, r, ar)
	}
	osin.OutputJSON(resp, w, r)

	/*if err != nil {
		// We return a status error here, which conveniently wraps the error
		// returned from our DB queries. We can clearly define which errors
		// are worth raising a HTTP 500 over vs. which might just be a HTTP
		// 404, 403 or 401 (as appropriate). It's also clear where our
		// handler should stop processing by returning early.
		return StatusError{404, fmt.Errorf("User not found")}
	}*/


	return nil
}




func GetConsumption(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error) {
	//server := kernel.Container.MustGet("oauth").(*osin.Server)
	//database := kernel.Container.MustGet("database").(*gorm.DB)

	//"alvaro_gg@hotmail.com"
	//"MBAR4B1"

	buf, err := ioutil.ReadAll(r.Body)
	requestData := &protomodel.Credentials{}
	err = proto.Unmarshal(buf, requestData)
	if err != nil {
		return nil, apirest.StatusError{400, err}
	}

	//requestData.Password = "pepe"
	//v := &CredentialsValidator{Username: requestData.Username, Password:requestData.Password}

	_, err = apirest.Validate(requestData, &CredentialsValidator{})
	if err != nil {
		return nil, apirest.StatusError{400, err}
	}

	/*username := r.Form.Get("username")
	password := r.Form.Get("password")
	username = "alvaro_gg@hotmail.com"
	password = "MBAR4B1"

	fmt.Println(username)*/
	// Enviar tarea
	task0 := tasks.Signature{
		Name: "consumption",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: requestData.Username,
			},
			{
				Type:  "string",
				Value: requestData.Password,
			},
		},
	}

	response, err := apirest.SendTask(kernel, &task0)
	if err != nil {
		return nil, apirest.StatusError{http.StatusInternalServerError, err}
	}

	return response, nil
}

func GetTaskState(kernel *dislet.Kernel, w http.ResponseWriter, r *http.Request) (proto.Message, error) {
	server := kernel.Container.MustGet("machinery").(*machinery.Server)
	//api := kernel.Container.MustGet("api").(*mux.Router)
	vars := mux.Vars(r)
	taskUid := vars["taskUid"]


	task, err := server.GetBackend().GetState(taskUid)


	if err != nil {
		return nil, apirest.StatusError{http.StatusNotFound, fmt.Errorf("Task not found")}
	}

	data := apirest.TaskResponseHandler(task)

	return data, nil
}