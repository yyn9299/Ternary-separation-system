const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');

function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}

//跳转到发布任务列表页面
router.get('/pushTask',async (ctx)=>{
    console.log('任务发布页面')
    await ctx.render('users/tasks/pushTask');
})

//查询用户发布任务列表页面
router.get('/pushTasklist',async (ctx)=>{
    console.log('文件列表页面')
    let username = ctx.session.username;
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetUserPubTaskByUname',username);
    let pushTasklist = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
    var counts = pushTasklist.length;
    //分页展示 一页展示5条记录
    let page = ctx.query.page;
    page = page ? page : 1;
    let pushTasklist1 = pushTasklist.slice((page-1)*5,page*5);
    let options = {
        pushTasklist:pushTasklist1,
        counts:counts,
        page:page,
        //pages:pages,
    }
    //千万不要忘了异步 异步  await
    await ctx.render('users/tasks/pushTasklist',options);
})

//查询用户接受任务列表页面
router.get('/acceptTasklist',async (ctx)=>{
    console.log('接受任务列表页面')
    let username = ctx.session.username;
    // const gateway = await getContract(username); //这是封装的一个初始化获取gateway的函数,需要将用户名传入
    // const network = await gateway.getNetwork(channelName);
	// const contract = network.getContract(chaincodeName);
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetUserAcceptTaskByUname',username);
    let acceptTasklist = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
    var counts = acceptTasklist.length;
    //分页展示 一页展示5条记录
    let page = ctx.query.page;
    page = page ? page : 1;
    let acceptTasklist1 = acceptTasklist.slice((page-1)*5,page*5);
    let options = {
        acceptTasklist:acceptTasklist1,
        counts:counts,
        page:page,
        //pages:pages,
    }
    //千万不要忘了异步 异步  await
    await ctx.render('users/tasks/acceptTasklist',options);
})

//根据任务名获取任务信息
router.post('/GetTaskImByTname',async (ctx)=>{
    let username = ctx.session.username;
    let tname = ctx.request.body.tname;
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetTaskImByTname',tname);
    let tasklist = JSON.parse(prettyJSONString(result.toString()));
    console.log('获取任务信息成功'+task);
    var counts = tasklist.length;  //只会出现一条信息 也将其放入到页面中
    let page = ctx.query.page;
    page = page ? page : 1;
    //let tasklist1 = tasklist.slice((page-1)*5,page*5);
    let options = {
        tasklist:tasklist,
        counts:counts,
        page:page,
    }
    await ctx.render('users/tasks/acceptTasklist',options);
})

//根据任务名修改任务状态
router.post('/UpdateTaskStateByTname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username;
    let tname = ctx.request.body.tname;
    let tstate = parseInt(ctx.request.body.tstate);
    const contract = getcontract(username);
    await contract.submitTransaction('UpdateTaskStateByTname',tname,tstate);
    console.log('修改任务状态成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '修改任务状态成功';
    ctx.body = JSON.stringify(temp);
})

//根据任务状态和发布人查找任务信息
router.post('/GetTasksByStateAndPublish',async (ctx)=>{
    let username = ctx.session.username;
    let publisher = ctx.request.body.publisher;
    let tstate = parseInt(ctx.request.body.tstate);
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetTasksByStateAndPublish',publisher,tstate);
    let tasklist = JSON.parse(prettyJSONString(result.toString()));
    console.log('获取任务信息成功'+task);
    var counts = tasklist.length;  //只会出现一条信息 也将其放入到页面中
    let page = ctx.query.page;
    page = page ? page : 1;
    //let tasklist1 = tasklist.slice((page-1)*5,page*5);
    let options = {
        tasklist:tasklist,
        counts:counts,
        page:page,
    }
    await ctx.render('users/tasks/acceptTasklist',options);
})

//*******用户发布任务----需要添加监听器 监听发布任务事件******
router.post('/PublishTask',async (ctx)=>{ 
    let username = ctx.session.username;
    let tname = ctx.request.body.taskname;
    let tdescription = ctx.request.body.tdescription;s
    let epochs = parseInt(ctx.request.body.epochs);
    let tnums = parseInt(ctx.request.body.tnums);
    let tdataset = ctx.request.body.tdataset;
    let tpublisher = username;
    let tdatacoins = parseInt(ctx.request.body.tdatacoins);
    let tcompcoins = parseInt(ctx.request.body.tcompcoins);
    let tcomp = ctx.request.body.tcomp;
    const contract = getcontract(username);
    let result = await contract.submitTransaction('PublishTask',tname,tdescription,epochs,tnums,tdataset,
                                                                tpublisher,tdatacoins,tcompcoins,tcomp);

    if(`${result}` !== ''){
        let temp = {};
        temp.success= 'ok';
        temp.data = '发布任务成功';
        ctx.body = JSON.stringify(temp);
      }else{
        console.log('发布任务失败');
      }                                              
})

//添加任务tname到用户接受任务列表 -也就是用户点击接受任务所执行
router.post('/AddUserAcceptTaskByUname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username;
    let tname = ctx.request.body.tname;
    const contract = getcontract(username);
    await contract.submitTransaction('AddUserAcceptTaskByUname',username,tname);
    console.log('接受任务成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '接收任务成功';
    ctx.body = JSON.stringify(temp);
})

//更新聚合前模型存放地址
router.post('/UpdateModelhash',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username;
    let tname = ctx.request.body.tname;
    let modelhash = ctx.request.body.modelhash;
    const contract = getcontract(username);
    await contract.submitTransaction('UpdateModelhash',tname,modelhash,username);
    console.log('更新聚合前模型成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '更新聚合前模型成功';
    ctx.body = JSON.stringify(temp);
})

//更新聚合后模型存放地址
router.post('/UpdateAggModelhash',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username;
    let tname = ctx.request.body.tname;
    let aggmodelhash = ctx.request.body.aggmodelhash;
    const contract = getcontract(username);
    await contract.submitTransaction('UpdateAggModelhash',tname,aggmodelhash,username);
    console.log('更新聚合后模型成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '更新聚合后模型成功';
    ctx.body = JSON.stringify(temp);
})

//根据任务名，算力者uname找到对应聚合前最新模型存储地址
router.post('/GetModelhashByTnameAndUname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username; //当前登录用户的name
    let tname = ctx.request.body.tname;
    let uname = ctx.request.body.uname;  //算力者的name
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetModelhashByTnameAndUname',tname,uname);
    let modelhash = JSON.parse(prettyJSONString(result.toString()));
    //此处可以将modelhash返回给页面**  待处理
    console.log('查找聚合前模型成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '查找聚合前模型成功';
    ctx.body = JSON.stringify(temp);
})

//根据任务名，算力者uname找到对应聚合前最新模型存储地址
router.post('/GetAggModelhashByTnameAndUname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username; //当前登录用户的name
    let tname = ctx.request.body.tname;
    let uname = ctx.request.body.uname;  //算力者的name
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetAggModelhashByTnameAndUname',tname,uname);
    let aggmodelhash = JSON.parse(prettyJSONString(result.toString()));
    //此处可以将aggmodelhash返回给页面**  待处理
    console.log('查找聚合后模型成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '查找聚合后模型成功';
    ctx.body = JSON.stringify(temp);
})

//根据任务名，算力者uname找到模型训练代码存储地址
router.post('/GetCodehashByTnameAndUname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username; //当前登录用户的name
    let tname = ctx.request.body.tname;
    let uname = ctx.request.body.uname;  //算力者的name
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetCodehashByTnameAndUname',tname,uname);
    let codehash = JSON.parse(prettyJSONString(result.toString()));
    //此处可以将aggmodelhash返回给页面**  待处理
    console.log('查找模型训练代码成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '查找模型训练代码成功';
    ctx.body = JSON.stringify(temp);
})

//查询竞争算力用户
router.post('/GetCompTenderByTname',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username; //当前登录用户的name
    let tname = ctx.request.body.tname;
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('GetCompTenderByTname',tname);
    let cplist = JSON.parse(prettyJSONString(result.toString()));
    //此处可以将cplist返回给页面**  待处理
    console.log('查询竞争算力用户成功！'+tname);
    let temp = {};
    temp.success= 'ok';
    temp.data = '查询竞争算力用户成功';
    ctx.body = JSON.stringify(temp);
})

//删除任务
router.post('/deltaskByID',async (ctx,next) =>{
    console.log(ctx.request.body)
    let taskid = ctx.request.body.taskid;
    let sqlStr = "delete from task where id=?";
    let result = await sqlQuery(sqlStr,[taskid]);
    console.log('删除成功！'+taskid);
    let temp = {};
    temp.success= 'ok';
    temp.data = '删除任务成功';
    ctx.body = JSON.stringify(temp);
    console.log(result);
})
module.exports = router