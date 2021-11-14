/*
将数据集id加入用户数据库id表方法合并到创建数据库
将在用户结构体中删除对应的数据集id整合到删除数据库中
将设置用户算力id整合到创建算力
将更新任务运行轮次整合到更新聚合后模型地址中，除了tname和模型地址外还需要输入任务轮次
查找聚合前和后模型地址，更新为返回字符串，即最新的地址
tcode从任务结构体中删除，变为一个map的公共字段存储训练代码地址 增加了对tcode的初始化，更新tcode函数和查询tcode
任务状态变为7个分别是未开始匹配数据集，未开始匹配算力，未开始上传书数据集,未开始上传模型训练代码，运行中，任务结束和任务失败用1 2 3 4 5 6 7 来进行标识
增加数据集选择完成函数，在选择完数据集后，会变更任务状态并发出事件让算力开始竞争
算力竞争函数中，在最后一个算力竞争者加入竞争集后，会自动触发初始化，并发送事件使得数据集提供者上传数据集，数据集提供者上传后，会根据map中各个值是否含有
判断是否添加完毕，发送事件让模型提供者上传训练代码，代码上传完毕后会发送事件，让任务开始。
修改修改数据集和算力信息函数，使得输入可以为空
*/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type User struct {
	U_name        string   `json : U_name`
	U_description string   `json : U_description`
	U_pbkey       string   `json : U_pbkey`
	U_email       string   `json : U_email`
	U_role        int      `json : U_role `
	U_publish     []string `json : U_publish`
	U_coins       int      `json : U_coins`
	U_accept      []string `json : U_accept`
	U_datalist    []string `json : U_datalist`
	U_comp        string   `json : U_comp`
	U_status      int      `json : U_status`
}
type Dataset struct {
	Doctype       string `json : Doctype`
	D_id          string `json : D_id`
	D_type        string `json : D_type`
	D_description string `json : D_description`
	D_ownid       string `json : D_ownid`
	D_coin        int    `json : D_coin`
}
type Comput struct {
	Doctype  string `json : Doctype`
	C_id     string `json : C_id`
	C_cpu    string `json : C_cpu`
	C_memory string `json : C_memory`
	C_coin   int    `json : C_coin`
}
type Task struct {
	T_name         string              `json : T_name`
	T_description  string              `json : T_description`
	T_epochs       int                 `json : T_epochs`
	T_curepoch     int                 `json : T_curepoch`
	T_nums         int                 `json : T_nums`
	T_dataset      string              `json : T_dataset`
	T_computation  string              `json : T_computation`
	T_datahash     map[string][]string `json : T_datahash`
	T_publisher    string              `json : T_publisher`
	T_acceptdata   []string            `json : T_accpetdata`
	T_acceptcompor []string            `json : T_acceptcompor`
	T_datacoins    int                 `json : T_datacoins`
	T_compcoins    int                 `json : T_compcoins`
	T_state        int                 `json : T_state`
}
type Smartcontract struct {
	contractapi.Contract
}

//创建一个用户
func (s *Smartcontract) CreateUser(ctx contractapi.TransactionContextInterface, uname string, upbkey string, uemail string, ucoins int, udescription string) error {
	upublish := []string{}
	uaccept := []string{}
	udatalist := []string{}

	user := User{
		U_name:        uname,
		U_pbkey:       upbkey,
		U_description: udescription,
		U_email:       uemail,
		U_role:        1,
		U_coins:       ucoins,
		U_accept:      uaccept,
		U_datalist:    udatalist,
		U_publish:     upublish,
		U_status:      0,
	}
	userjson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("create user failed")
	}

	err = ctx.GetStub().PutState(user.U_name, userjson)
	if err != nil {
		return fmt.Errorf("create user failed")

	}
	return nil
}

//查询所有用户
func (s *Smartcontract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	querystring := `{"selector":{"U_role":1}}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(querystring)
	if err != nil {
		return nil, fmt.Errorf("get all users failed")
	}
	defer resultsIterator.Close()
	var users []*User
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var user User
		err = json.Unmarshal(queryResult.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}


//根据用户uname查询用户信息
func (s *Smartcontract) GetUserImByUname(ctx contractapi.TransactionContextInterface, uname string) (*User, error) {
	userjson, err := ctx.GetStub().GetState(uname)
	if err != nil {
		return nil, fmt.Errorf("get user failed")
	}
	if userjson == nil {
		return nil, fmt.Errorf("this user doesn't exist")
	}
	var user *User
	user = new(User)
	err = json.Unmarshal(userjson, user)
	if err != nil {

		return nil, fmt.Errorf("get user failed")
	}

	return user, nil

}

//根据用户uname修改用户状态
func (s *Smartcontract) UpdateUserstatusByUname(ctx contractapi.TransactionContextInterface, uname string, status int) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("Update user status failed because : %s", err)

	}
	user.U_status = status
	userjson, err1 := json.Marshal(user)
	if err1 != nil {

		return fmt.Errorf("Update user status failed")
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {

		return fmt.Errorf("Update user status failed")
	}
	return nil

}

//根据用户uname修改用户拥有代币
func (s *Smartcontract) UpdateUsercoinByUname(ctx contractapi.TransactionContextInterface, uname string, ucoins int) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("Update user coins failed because : %s", err)

	}
	user.U_coins = ucoins
	userjson, err1 := json.Marshal(user)
	if err1 != nil {

		return fmt.Errorf("Update user coins failed")
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {
		return fmt.Errorf("Update user coins failed")
	}
	return nil

}



//添加任务tname到用户发布任务列表
func (s *Smartcontract) AddUserPubTaskByUname(ctx contractapi.TransactionContextInterface, uname string, tname string) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("add user published task failed because %s", err)
	}

	user.U_publish = append(user.U_publish, tname)
	userjson, err1 := json.Marshal(user)
	if err1 != nil {
		return fmt.Errorf("add user published task failed")
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {

		return fmt.Errorf("add user published task failed")
	}
	return nil

}


//添加任务tname到用户接受任务列表
func (s *Smartcontract) AddUserAcceptTaskByUname(ctx contractapi.TransactionContextInterface, uname string, tname string) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("add user accepted task failed because %s", err)
	}

	user.U_accept = append(user.U_accept, tname)
	userjson, err1 := json.Marshal(user)
	if err1 != nil {
		return fmt.Errorf("add user accepted task failed")
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {

		return fmt.Errorf("add user accepted task failed")
	}
	return nil

}

//did应该是系统生成的

这样写有问题

func (s *Smartcontract) CreateData(ctx contractapi.TransactionContextInterface, did string, dtype string, description string, uname string, coin int) error {
	dataset := Dataset{
		Doctype:       "data",
		D_id:          did,
		D_type:        dtype,
		D_description: description,
		D_ownid:       uname,
		D_coin:        coin,
	}
	datasetjson, err := json.Marshal(dataset)
	if err != nil {
		return fmt.Errorf("create dataset failed")
	}

	err1 := s.AddUserDatalistByUname(ctx, uname, did)
	if err1 != nil {
		return fmt.Errorf("create dataset failed because :%s", err1)
	}

	err = ctx.GetStub().PutState(dataset.D_id, datasetjson)
	if err != nil {
		return fmt.Errorf("create dataset failed")

	}

	return nil
}

//添加新数据集id到用户数据库id列表
func (s *Smartcontract) AddUserDatalistByUname(ctx contractapi.TransactionContextInterface, uname string, dataid string) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("add user datalist failed because %s", err)
	}

	user.U_datalist = append(user.U_datalist,dataid)

	userjson, err1 := json.Marshal(user)
	if err1 != nil {
		return fmt.Errorf("add user datalist failed")
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {

		return fmt.Errorf("add user datalist failed")
	}
	return nil

}

//根据数据集id删除对应的数据集
func (s *Smartcontract) DelDatasetByDid(ctx contractapi.TransactionContextInterface, did string) error {
	datajson, err := ctx.GetStub().GetState(did)
	if err != nil {
		fmt.Errorf("del dataset failed")
	}
	if datajson==nil {
		fmt.Errorf("dataset is not exist")

	}
	err1 := ctx.GetStub().DelState(did)
	if err1 != nil {
		fmt.Errorf("del dataset failed")
	}
	dataset, err2 := s.SelectDataImByDid(ctx, did)
	if err2 != nil {
		fmt.Errorf("del dataset failed because %s", err2)
	}
	err3 := s.DelUserDatalistByUname(ctx, dataset.D_ownid, did)
	if err3 != nil {
		fmt.Errorf("del dataset failed because %s", err3)
	}
	return nil
}

//根据传输数据集id在用户拥有的数据集id中删除对应的数据集id
func (s *Smartcontract) DelUserDatalistByUname(ctx contractapi.TransactionContextInterface, uname string, dataid string) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("del user datalist failed because %s", err)
	}

	for index,val := range user.U_datalist {
		if val == dataid {
			user.U_datalist=append(user.U_datalist[:index],user.U_datalist[index+1:]...)
		}
	}
	userjson, err1 := json.Marshal(user)
	if err1 != nil {
		return fmt.Errorf("del user datalist failed")
	}
	err1=ctx.GetStub().PutState(uname, userjson)
	if err1!=nil{
		return fmt.Errorf("del user datalist failed")
	}

	return nil

}

//创建一个算力
cid应该也是系统生成

func (s *Smartcontract) CreateComp(ctx contractapi.TransactionContextInterface, cid string, cpu string, cmemory string, coin int, uname string) error {
	comp := Comput{
		Doctype:  "computation",
		C_id:     cid,
		C_cpu:    cpu,
		C_memory: cmemory,
		C_coin:   coin,
	}
	compjson, err := json.Marshal(comp)
	if err != nil {
		return fmt.Errorf("create comput failed")
	}

	err = ctx.GetStub().PutState(comp.C_id, compjson)
	if err != nil {
		return fmt.Errorf("failed to put compt imformation to worldstate")

	}
	err1 := s.SetUserCompByUname(ctx, uname, cid)
	if err1 != nil {
		return fmt.Errorf("failed to create Comp  because : %s", err1)

	}
	return nil
}
//设置用户算力id
func (s *Smartcontract) SetUserCompByUname(ctx contractapi.TransactionContextInterface, uname string, comp string) error {
	user, err := s.GetUserImByUname(ctx, uname)
	if err != nil {
		return fmt.Errorf("Set  UserComp failed because %s", err)
	}

	user.U_comp = comp
	userjson, err1 := json.Marshal(user)
	if err1 != nil {
		return fmt.Errorf("Set userComp is failed because : %s", err1)
	}
	err1 = ctx.GetStub().PutState(uname, userjson)
	if err1 != nil {

		return fmt.Errorf("Set userComp is failed because : %s", err1)
	}
	return nil

}

////查询所有黑名单用户uname
//func (s *Smartcontract) GetAllBlacklist(ctx contractapi.TransactionContextInterface) ([]string, error) {
//	userblackjson, err := ctx.GetStub().GetState("UserblackUname")
//	userblack := make([]string, 1)
//	err1 := json.Unmarshal(userblackjson, &userblack)
//
//	if err1 != nil {
//		return nil, err1
//	}
//	if err != nil {
//		return nil, fmt.Errorf("Get blacklist is failed because %s", err)
//
//	}
//	return userblack, nil
//
//}
//
////设置用户进入黑名单
//func (s *Smartcontract) SetUserBan(ctx contractapi.TransactionContextInterface, uname string) error {
//	userblack, err := s.GetAllBlacklist(ctx)
//	if err != nil {
//		return fmt.Errorf("get Blacklist is failed because %s", err)
//	}
//	userblack = append(userblack, uname)
//	userblackjson, err1 := json.Marshal(&userblack)
//	if err1 != nil {
//
//		return fmt.Errorf("Set user ban is failed because : %s", err1)
//	}
//	err2 := ctx.GetStub().PutState("UserblackUname", userblackjson)
//	if err2 != nil {
//
//		return fmt.Errorf("Set user ban is failed because : %s", err2)
//	}
//	return nil
//
//}



//根据数据集id选取对应的数据集
func (s *Smartcontract) SelectDataImByDid(ctx contractapi.TransactionContextInterface, did string) (*Dataset, error) {
	datasetjson, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from world state")
	}
	if datasetjson == nil {
		return nil, fmt.Errorf("nobody is in this chain block")
	}
	var dataset *Dataset
	dataset = new(Dataset)
	err = json.Unmarshal(datasetjson, dataset)
	if err != nil {

		return nil, fmt.Errorf("Dataset Query is failed beacuse:%s", err)
	}

	return dataset, nil

}





//修改数据集类型和描述
func (s *Smartcontract) UpdateDatasetByDid(ctx contractapi.TransactionContextInterface, did string, dtype string, description string, dcoin int) error {
	dataset, err := s.SelectDataImByDid(ctx, did)


	if err != nil {
		fmt.Errorf("Update Dataset is failed because %s", err)
	}
	if dataset== nil {
		fmt.Errorf("Update Dataset is failed because data doesn't exist")
	}

	if description != "" {
		dataset.D_description = description
	}
	if dtype != "" {
		dataset.D_type = dtype
	}

	dataset.D_coin = dcoin

	datasetjson, err1 := json.Marshal(dataset)
	if err1 != nil {
		fmt.Errorf("Update Data is failed because %s", err1)
	}
	err1 = ctx.GetStub().PutState(did, datasetjson)
	if err1 != nil {
		fmt.Errorf("Update Dataset is failed beacuse %s", err1)
	}
	return nil

}

//查询所有数据集
func (s *Smartcontract) GetAllDatasets(ctx contractapi.TransactionContextInterface) ([]*Dataset, error) {
	querystring := `{"selector":{"Doctype":"data"}}`
	return getQueryDataResultForQueryString(ctx, querystring)

}

//根据数据集类型查找数据集
func (s *Smartcontract) GetDatasetsByDtype(ctx contractapi.TransactionContextInterface, dtype string) ([]*Dataset, error) {
	querystring := fmt.Sprintf(`{"selector":{"D_type":"%s"}}`, dtype)
	return getQueryDataResultForQueryString(ctx, querystring)

}
func constructQueryDataResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Dataset, error) {
	var datas []*Dataset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var data Dataset
		err = json.Unmarshal(queryResult.Value, &data)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, nil
}
func getQueryDataResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Dataset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryDataResponseFromIterator(resultsIterator)
}



//根据算力id查询算力数据
func (s *Smartcontract) SelectCompImByCid(ctx contractapi.TransactionContextInterface, cid string) (*Comput, error) {
	compjson, err := ctx.GetStub().GetState(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to get state from world state")
	}
	if compjson == nil {
		return nil, fmt.Errorf("comp %s is in this chain block", cid)
	}
	var comp *Comput
	comp = new(Comput)
	err = json.Unmarshal(compjson, comp)
	if err != nil {

		return nil, fmt.Errorf("Comp Query is failed beacuse:%s", err)
	}

	return comp, nil

}



//修改算力内容
func (s *Smartcontract) UpdateCompByCid(ctx contractapi.TransactionContextInterface, cid string, cpu string, cmemory string, coin int) error {
	comp, err := s.SelectCompImByCid(ctx, cid)

	if err != nil {
		fmt.Errorf("Update Comp is failed because %s", err)
	}
	if comp== nil {
		fmt.Errorf("Update Comp is failed because computation doesn't exist")
	}

	if cpu != "" {
		comp.C_cpu = cpu
	}
	if cmemory != "" {
		comp.C_memory = cmemory
	}

	comp.C_coin = coin

	compjson, err1 := json.Marshal(comp)
	if err1 != nil {
		fmt.Errorf("Update comp is failed because %s", err1)
	}
	err1 = ctx.GetStub().PutState(cid, compjson)
	if err1 != nil {
		fmt.Errorf("Update Compset is failed beacuse %s", err1)
	}
	return nil

}

//查询所有算力信息
func (s *Smartcontract) GetAllComps(ctx contractapi.TransactionContextInterface) ([]*Comput, error) {
	querystring := `{"selector":{"Doctype":"comput"}}`
	return getQueryCompResultForQueryString(ctx, querystring)

}

func constructQueryCompResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Comput, error) {
	var comps []*Comput
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var comp Comput
		err = json.Unmarshal(queryResult.Value, &comp)
		if err != nil {
			return nil, err
		}
		comps = append(comps, &comp)
	}

	return comps, nil
}
func getQueryCompResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Comput, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryCompResponseFromIterator(resultsIterator)
}

//根据任务名获取任务信息
func (s *Smartcontract) GetTaskImByTname(ctx contractapi.TransactionContextInterface, tname string) (*Task, error) {
	taskjson, err := ctx.GetStub().GetState(tname)
	if err != nil {
		return nil, fmt.Errorf("Fail to get data from world state")
	}
	if taskjson == nil {
		return nil, fmt.Errorf("Task %s isn't in this chain block", tname)
	}
	var task *Task
	task = new(Task)
	err1 := json.Unmarshal(taskjson, task)
	if err1 != nil {

		return nil, fmt.Errorf("Task Query is failed beacuse:%s", err1)
	}

	return task, nil

}

//根据任务名修改任务状态
func (s *Smartcontract) UpdateTaskStateByTname(ctx contractapi.TransactionContextInterface, tname string, tstate int) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Update Task State is failed because : %s", err)

	}
	task.T_state = tstate
	taskjson, err1 := json.Marshal(task)
	if err1 != nil {

		return fmt.Errorf("Update Task state is failed because : %s", err1)
	}
	err1 = ctx.GetStub().PutState(tname, taskjson)
	if err1 != nil {

		return fmt.Errorf("update task state is failed because : %s", err1)
	}
	return nil

}

////根据任务状态和发布人查找任务
//func (s *Smartcontract) GetTasksByStateAndPublish(ctx contractapi.TransactionContextInterface, publisher string, tstate int) ([]*Task, error) {
//	querystring := fmt.Sprintf(`{"selector":{"T_state":%d,"T_publisher":"%s"}}`, tstate, publisher)
//	return getQueryResultTaskForQueryString(ctx, querystring)
//
//}
//func constructQueryResponseTaskFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Task, error) {
//	var tasks []*Task
//	for resultsIterator.HasNext() {
//		queryResult, err := resultsIterator.Next()
//		if err != nil {
//			return nil, err
//		}
//		var task Task
//		err = json.Unmarshal(queryResult.Value, &task)
//		if err != nil {
//			return nil, err
//		}
//		tasks = append(tasks, &task)
//	}
//
//	return tasks, nil
//}
//func getQueryResultTaskForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Task, error) {
//	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
//	if err != nil {
//		return nil, err
//	}
//	defer resultsIterator.Close()
//
//	return constructQueryResponseTaskFromIterator(resultsIterator)
//}

//更新聚合前模型存放地址
func (s *Smartcontract) UpdateModelhash(ctx contractapi.TransactionContextInterface, tname string, modelhash string, uname string) error {
	modelhashjson, err := ctx.GetStub().GetState(tname + "modelhash")
	if err != nil {
		return fmt.Errorf("Update modelhash failed because %s", err)
	}
	modelhash1 := make(map[string][]string)
	err1 := json.Unmarshal(modelhashjson, &modelhash1)
	if modelhash1[uname][0] == "" {

		modelhash1[uname][0] = modelhash

	} else {
		modelhash1[uname] = append(modelhash1[uname], modelhash)
	}
	newmodelhashjson, err2 := json.Marshal(modelhash1)
	if err1 != nil {
		return fmt.Errorf("Update modelhash is failed because : %s", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Update modelhash is failed because : %s", err2)
	}
	err2 = ctx.GetStub().PutState(tname+"modelhash", newmodelhashjson)
	if err2 != nil {

		return fmt.Errorf("update modelhash is failed because : %s", err2)
	}
	return nil

}

//更新聚合后的模型
func (s *Smartcontract) UpdateAggModelhash(ctx contractapi.TransactionContextInterface, tname string, aggmodelhash string, uname string, tcurepo int) error {
	aggmodelhashjson, err := ctx.GetStub().GetState(tname + "aggmodelhash")
	if err != nil {
		return fmt.Errorf("Update aggmodelhash failed because %s", err)
	}
	aggmodelhash1 := make(map[string][]string)
	err1 := json.Unmarshal(aggmodelhashjson, &aggmodelhash1)
	if aggmodelhash1[uname][0] == "" {
		aggmodelhash1[uname][0] = aggmodelhash

	} else {
		aggmodelhash1[uname] = append(aggmodelhash1[uname], aggmodelhash)
	}
	newmodelhashjson, err2 := json.Marshal(aggmodelhash1)
	if err1 != nil {
		return fmt.Errorf("Update aggmodelhash is failed because : %s", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Update aggmodelhash is failed because : %s", err2)
	}
	err2 = ctx.GetStub().PutState(tname+"aggmodelhash", newmodelhashjson)
	if err2 != nil {

		return fmt.Errorf("update aggmodelhash is failed because : %s", err2)
	}
	err3 := s.UpdateTaskCurepochByTname(ctx, tname, tcurepo)
	if err3 != nil {

		return fmt.Errorf("update aggmodelhash is failed because : %s", err3)
	}
	return nil

}

//判断任务是否存在
func (s *Smartcontract) TaskExist(ctx contractapi.TransactionContextInterface, tname string) (bool, error) {
	taskjson, err := ctx.GetStub().GetState(tname)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state")
	}
	return taskjson != nil, nil
}

//更新任务运行轮次
func (s *Smartcontract) UpdateTaskCurepochByTname(ctx contractapi.TransactionContextInterface, tname string, tcurepo int) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Update Task curepoch is failed because : %s", err)

	}
	task.T_curepoch = tcurepo
	taskjson, err1 := json.Marshal(task)
	if err1 != nil {

		return fmt.Errorf("Update Task curepoch is failed because : %s", err1)
	}
	err1 = ctx.GetStub().PutState(tname, taskjson)
	if err1 != nil {

		return fmt.Errorf("update task curepoch is failed because : %s", err1)
	}
	return nil

}

//发布任务，需要填入模型需求以及对聚合前和聚合后模型存放地址，竞标算力者字段和数据集竞标者字段，模型训练代码进行初始化
func (s *Smartcontract) PublishTask(ctx contractapi.TransactionContextInterface, tname string, tdescription string, epochs int, tnums int,
	tdataset string, tpublisher string, tdatacoins int, tcompcoins int, tcomp string) error {
	tdatahash := make(map[string][]string)
	tacceptcom := []string{}
	tmodelhash := make(map[string][]string)
	taggmodelhash := make(map[string][]string)

	tdatatender := []string{}
	tacceptdata := []string{}
	tcode := make(map[string][]string)

	task := Task{
		T_name:         tname,
		T_description:  tdescription,
		T_epochs:       epochs,
		T_curepoch:     -1,
		T_nums:         tnums,
		T_dataset:      tdataset,
		T_computation:  tcomp,
		T_datahash:     tdatahash,
		T_publisher:    tpublisher,
		T_acceptdata:   tacceptdata,
		T_acceptcompor: tacceptcom,
		T_datacoins:    tdatacoins,
		T_compcoins:    tcompcoins,
		T_state:        1,
	}
	taskjson, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("create task failed")
	}

	err = ctx.GetStub().PutState(task.T_name, taskjson)
	if err != nil {
		return fmt.Errorf("failed to put task imformation to worldstate")

	}
	modelhashjson, err1 := json.Marshal(tmodelhash)
	if err1 != nil {
		return fmt.Errorf("initilize modelhash failed")
	}

	err1 = ctx.GetStub().PutState(task.T_name+"modelhash", modelhashjson)
	if err1 != nil {
		return fmt.Errorf("failed to initilize modelhash to worldstate")

	}
	aggmodelhashjson, err5 := json.Marshal(taggmodelhash)
	if err5 != nil {
		return fmt.Errorf("initilize aggmodelhash failed")
	}

	err5 = ctx.GetStub().PutState(task.T_name+"aggmodelhash", aggmodelhashjson)
	if err5 != nil {
		return fmt.Errorf("failed to initilize aggmodelhash to worldstate")

	}
	tdatatenderjson, err2 := json.Marshal(tdatatender)
	if err2 != nil {
		return fmt.Errorf("initilize datatender failed")
	}

	err2 = ctx.GetStub().PutState(task.T_name+"dataers", tdatatenderjson)
	if err2 != nil {
		return fmt.Errorf("failed to initilize datatender to worldstate")

	}

	codejson, err5 := json.Marshal(tcode)
	if err5 != nil {
		return fmt.Errorf("initilize tcode failed")
	}

	err5 = ctx.GetStub().PutState(task.T_name+"code", codejson)
	if err5 != nil {
		return fmt.Errorf("failed to initilize code to worldstate")

	}
	err4 := ctx.GetStub().SetEvent("find", taskjson)
	if err4 != nil {
		return fmt.Errorf("fail to setEvent find")
	}

	return nil
}

//根据任务名，模型方找到对应聚合前模型存储地址
func (s *Smartcontract) GetModelhashByTnameAndUname(ctx contractapi.TransactionContextInterface, tname string) ([]string, error) {
	modeljson, err := ctx.GetStub().GetState(tname + "modelhash")

	if err != nil {
		return nil, fmt.Errorf("Fail to get modelhash from world state")
	}
	if modeljson == nil {
		return nil, fmt.Errorf("Task%s Modelhash is empty", tname)
	}
	taskjson, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return nil,fmt.Errorf("Get task failed")
	}

	modelhash := make(map[string][]string)
	err1 := json.Unmarshal(modeljson, &modelhash)
	if err1 != nil {

		return nil, fmt.Errorf("Task Modelhash Query is failed beacuse:%s", err1)
	}
	curmodelhash:=[]string{}
	for _, value := range modelhash {
		curmodelhash=append(curmodelhash,value[taskjson.T_epochs])
	}
	return curmodelhash, nil

}

//根据任务名，算力者uname找到对应聚合后模型存储地址
func (s *Smartcontract) GetAggModelhashByTnameAndUname(ctx contractapi.TransactionContextInterface, tname string, uname string) (string, error) {
	aggmodeljson, err := ctx.GetStub().GetState(tname + "aggmodelhash")
	if err != nil {
		return "", fmt.Errorf("Fail to get aggmodelhash from world state")
	}
	if aggmodeljson == nil {
		return "", fmt.Errorf("Task%s Aggmodelhash is empty", tname)
	}
	aggmodelhash := make(map[string][]string)
	err1 := json.Unmarshal(aggmodeljson, &aggmodelhash)
	if err1 != nil {

		return "", fmt.Errorf("Task Aggmodelhash Query is failed beacuse:%s", err1)
	}
	taskjson, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return "",fmt.Errorf("Get task failed")
	}

	return aggmodelhash[uname][taskjson.T_curepoch], nil

}

//根据任务名，算力者uname找到模型训练代码存储地址
func (s *Smartcontract) GetCodehashByTnameAndUname(ctx contractapi.TransactionContextInterface, tname string, uname string) (string, error) {
	codejson, err := ctx.GetStub().GetState(tname + "code")
	if err != nil {
		return "", fmt.Errorf("Fail to get code from world state")
	}
	if codejson == nil {
		return "", fmt.Errorf("Task%s code is empty", tname)
	}
	code := make(map[string][]string)
	err1 := json.Unmarshal(codejson, &code)
	if err1 != nil {

		return "", fmt.Errorf("Task codehash Query is failed beacuse:%s", err1)
	}
	taskjson, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return "",fmt.Errorf("Get task failed")
	}
	return code[uname][taskjson.T_curepoch], nil

}



//添加竞争算力用户
func (s *Smartcontract) Addcomputation(ctx contractapi.TransactionContextInterface, uname string, tname string) error {
	taskjson, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Get task failed")
	}
	if len(taskjson.T_acceptcompor)<taskjson.T_nums{
		taskjson.T_acceptcompor=append(taskjson.T_acceptcompor,uname)
	}else {

	}
}

//查找数据集竞标列表
func (s *Smartcontract) GetDataTenderByTname(ctx contractapi.TransactionContextInterface, tname string) ([]string, error) {
	datajson, err := ctx.GetStub().GetState(tname + "dataers")
	if err != nil {
		return nil, fmt.Errorf("Fail to get Datatender from world state")
	}
	if datajson == nil {
		return nil, nil
	}
	datatender := make([]string, 1)

	err1 := json.Unmarshal(datajson, &datatender)
	if err1 != nil {

		return nil, fmt.Errorf("Task datatender Query is failed beacuse:%s", err1)
	}

	return datatender, nil

}

//添加数据集id到数据集竞标列表内
func (s *Smartcontract) AddDatatender(ctx contractapi.TransactionContextInterface, did string, tname string) error {
	datatenderjson, err := ctx.GetStub().GetState(tname + "dataers")
	if err != nil {
		return fmt.Errorf("Add Datatender failed because %s", err)
	}
	data := make([]string, 1)
	err1 := json.Unmarshal(datatenderjson, &data)
	if data[0] == "" {
		data[0] = did

	} else {
		data = append(data, did)
	}
	datajson, err2 := json.Marshal(data)
	if err1 != nil {
		return fmt.Errorf("Add Datatender is failed because : %s", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Add Datatender is failed because : %s", err2)
	}
	err2 = ctx.GetStub().PutState(tname+"dataers", datajson)
	if err2 != nil {

		return fmt.Errorf("Add Datatender is failed because : %s", err2)
	}
	return nil

}

//补充任务结构体的选择的数据集信息，将任务结构体填入数据集id
func (s *Smartcontract) SupplyTaskImByTnameAndDid(ctx contractapi.TransactionContextInterface, tname string, did string) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Supply Task  is failed because : %s", err)

	}

	if task.T_acceptdata[0] == "" {
		task.T_acceptdata[0] = did

	} else {
		task.T_acceptdata = append(task.T_acceptdata, did)
	}

	taskjson, err3 := json.Marshal(task)
	if err3 != nil {

		return fmt.Errorf("Supply Task  is failed because : %s", err3)
	}
	err1 := ctx.GetStub().PutState(tname, taskjson)
	if err1 != nil {

		return fmt.Errorf("Supply task  is failed because : %s", err1)
	}
	return nil

}

//数据集选择完成
func (s *Smartcontract) DatasetChoseFinish(ctx contractapi.TransactionContextInterface, tname string) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf(" Finish Dataset Chose  is failed because : %s", err)

	}
	task.T_state = 2
	taskjson, err1 := json.Marshal(task)
	if err1 != nil {

		return fmt.Errorf("Finish Dataset Chose  is failed because : %s", err1)
	}

	err1 = ctx.GetStub().PutState(tname, taskjson)
	if err1 != nil {

		return fmt.Errorf("Finish Dataset Chose   is failed because : %s", err1)
	}
	err4 := ctx.GetStub().SetEvent("comptend", taskjson)
	if err4 != nil {
		return fmt.Errorf("fail to setEvent start")
	}
	return nil

}

//上传数据集哈希地址
func (s *Smartcontract) UpdateDatahashByTname(ctx contractapi.TransactionContextInterface, dhash string, tname string, uname string) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Update Datahash failed because %s", err)
	}
	datahash := task.T_datahash

	if datahash[uname][0] == "" {
		datahash[uname][0] = dhash

	} else {
		datahash[uname] = append(datahash[uname], dhash)
	}

	flag := 0
	for _, v := range datahash {
		if v[0] != "" {
			flag++
		}
	}
	if flag == task.T_nums {
		task.T_datahash = datahash
		task.T_state = 4
		taskjson, err2 := json.Marshal(task)
		err4 := ctx.GetStub().SetEvent("upload code", taskjson)
		if err4 != nil {
			return fmt.Errorf("fail to setEvent upload code")
		}

		if err2 != nil {
			return fmt.Errorf("Update Datahash is failed because : %s", err2)
		}
		err2 = ctx.GetStub().PutState(tname, taskjson)
		if err2 != nil {

			return fmt.Errorf("Update Datahash is failed because : %s", err2)
		}

	} else {
		task.T_datahash = datahash
		taskjson, err2 := json.Marshal(task)

		if err2 != nil {
			return fmt.Errorf("Update Datahash is failed because : %s", err2)
		}
		err2 = ctx.GetStub().PutState(tname, taskjson)
		if err2 != nil {

			return fmt.Errorf("Update Datahash is failed because : %s", err2)
		}
	}

	return nil

}

//上传训练代码哈希地址
func (s *Smartcontract) UpdateCodehashByTname(ctx contractapi.TransactionContextInterface, codehash string, tname string, uname string) error {
	task, err5 := s.GetTaskImByTname(ctx, tname)
	if err5 != nil {
		return fmt.Errorf("Update code  is failed because : %s", err5)

	}
	codehashjson, err := ctx.GetStub().GetState(tname + "code")
	if err != nil {
		return fmt.Errorf("Update codehash failed because %s", err)
	}
	codehash1 := make(map[string][]string)
	err1 := json.Unmarshal(codehashjson, &codehash1)
	if codehash1[uname][0] == "" {

		codehash1[uname][0] = codehash

	} else {
		codehash1[uname] = append(codehash1[uname], codehash)
	}
	newcodehashjson, err2 := json.Marshal(codehash1)
	if err1 != nil {
		return fmt.Errorf("Update codehash is failed because : %s", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Update codehash is failed because : %s", err2)
	}
	err2 = ctx.GetStub().PutState(tname+"codehash", newcodehashjson)
	if err2 != nil {

		return fmt.Errorf("update codehash is failed because : %s", err2)
	}
	flag := 0
	for _, v := range codehash1 {
		if v[0] != "" {
			flag++
		}
	}
	if flag == task.T_nums {
		err4 := s.StartTask(ctx, tname)
		if err4 != nil {
			return fmt.Errorf("fail to start task")
		}

	}

	return nil

}

//任务开始函数，调用本函数前，任务结构体内部应该全部填好
func (s *Smartcontract) StartTask(ctx contractapi.TransactionContextInterface, tname string) error {
	task, err := s.GetTaskImByTname(ctx, tname)
	if err != nil {
		return fmt.Errorf("Start Task  is failed because : %s", err)

	}
	task.T_state = 5
	taskjson, err1 := json.Marshal(task)
	if err1 != nil {

		return fmt.Errorf("Start Task  is failed because : %s", err1)
	}

	err1 = ctx.GetStub().PutState(tname, taskjson)
	if err1 != nil {

		return fmt.Errorf("Start task  is failed because : %s", err1)
	}
	err4 := ctx.GetStub().SetEvent("start", taskjson)
	if err4 != nil {
		return fmt.Errorf("fail to setEvent start")
	}
	return nil

}


func main1()  {
	fmt.Println(1)

}