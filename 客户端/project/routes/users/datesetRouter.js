const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');
var fs = require('fs')
var  path = require('path')
var ipfsFile = require('../../module/ipfs');
const ipfsAPI = require('ipfs-api');
const ipfs = ipfsAPI({host: 'localhost', port: '5001', protocol: 'http'});

//跳转到上传数据集页面
router.get('/uploadfile',async (ctx)=>{
    console.log('上传数据页面')
    await ctx.render('users/date/uploadfile.ejs');
})
//跳转到创建新数据集页面
router.get('/addDateset',async (ctx)=>{
  console.log('添加数据集页面')
  await ctx.render('users/date/addDateset.ejs');
})

//添加一个数据集
router.post('/addDateset',async (ctx)=>{
    let uname = ctx.session.username;
    let did = ctx.request.body.did;
    let dtype = ctx.request.body.dtype;
    let ddescription = ctx.request.body.ddescription;
    let coin = parseInt(ctx.request.body.coin);
    const contract = getcontract(uname);
    let result = await contract.submitTransaction('CreateData',did,dtype,ddescription,uname,coin);
    if(`${result}` !== ''){
      let temp = {};
      temp.success= 'ok';
      temp.data = '添加数据集成功';
      ctx.body = JSON.stringify(temp);
    }else{
      console.log('添加数据集失败');
    }
})

//查询所有的数据集
router.get('/GetAllDatasets',async (ctx)=>{
  let uname = ctx.session.username;
  const contract = getcontract(uname);
  let result = await contract.evaluateTransaction('GetAllDatasets');
  let datesetlist = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
  var counts = datesetlist.length;
  //分页展示 一页展示5条记录
  let page = ctx.query.page;
  page = page ? page : 1;
  let datesetlist1 = datesetlist.slice((page-1)*5,page*5);
  let options = {
      datesetlist:datesetlist1,
      counts:counts,
      page:page,
  }
  //千万不要忘了异步 异步  await
  await ctx.render('users/date/datesetlist',options);
})

//删除数据集
router.post('/DelDatasetByDid',async (ctx,next) =>{
  console.log(ctx.request.body)
  let did = ctx.request.body.did;
  let username = ctx.session.username;
  const contract = getcontract(username);
  await contract.submitTransaction('DelDatasetByDid',did);
  console.log('删除成功！'+did);
  let temp = {};
  temp.success= 'ok';
  temp.data = '删除数据集成功';
  ctx.body = JSON.stringify(temp);
})

//修改数据集信息
router.post('/UpdateDatasetByDid',async (ctx,next) =>{
  console.log(ctx.request.body)
  let did = ctx.request.body.did;
  let username = ctx.session.username;
  let dtype = ctx.request.body.dtype;
  let ddescription = ctx.request.body.ddescription;
  let coin = parseInt(ctx.request.body.coin);
  const contract = getcontract(username);
  await contract.submitTransaction('UpdateDatasetByDid',did,dtype,ddescription,coin);
  console.log('修改成功！'+did);
  let temp = {};
  temp.success= 'ok';
  temp.data = '修改数据集成功';
  ctx.body = JSON.stringify(temp);
})

//根据数据集类型查找数据集
router.post('/GetDatasetsByDtype',async (ctx)=>{
  let uname = ctx.session.username;
  const contract = getcontract(uname);
  let dtype = ctx.request.body.dtype;
  let result = await contract.evaluateTransaction('GetDatasetsByDtype',dtype);
  let datesetlist = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
  var counts = datesetlist.length;
  //分页展示 一页展示5条记录
  let page = ctx.query.page;
  page = page ? page : 1;
  let datesetlist1 = datesetlist.slice((page-1)*5,page*5);
  let options = {
      datesetlist:datesetlist1,
      counts:counts,
      page:page,
  }
  //千万不要忘了异步 异步  await
  await ctx.render('users/date/datesetlist',options);
})

//根据数据集id选取数据集
router.post('/SelectDataImByDid',async (ctx,next) =>{
  console.log(ctx.request.body)
  let did = ctx.request.body.did;
  let username = ctx.session.username;
  const contract = getcontract(username);
  let result = await contract.evaluateTransaction('SelectDataImByDid',did);
  let dateset = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
  console.log('选取成功！'+dateset);
  let temp = {};
  temp.success= 'ok';
  temp.data = '选数据集成功';
  ctx.body = JSON.stringify(temp);
})


// 上传单个文件
router.post('/uploadfile', async (ctx, next) => {
    const file = ctx.request.files.file; // 获取上传文件
    //console.log(file);
    var name = file.name;
    var result =/\.[^\.]+/.exec(name);

    console.log(result);
    const filename = ctx.request.body.filename;

    //用异步方法读取流  将其转为buff 把buff传到ipfs中去
    const buff = fs.readFileSync(file.path);
    var sqlStr = 'insert into file values(default,?,?,?)';

    ipfs.add(buff,async (err,result)=>{
        if(err) throw err;
        console.log(result);
        var hash = result[0].hash;
        var size = (result[0].size) / (1024*1024);
        size = size.toFixed(2);
        // console.log
        var result = await sqlQuery(sqlStr,[name,hash,size]);
        console.log(result);
    })
    return ctx.body = "上传成功！" + result;
  });


//下载数据集
router.post('/downloadfile',async (ctx,next) =>{
    console.log('进入下载方法')
    var filehash = ctx.request.body.filehash;
    var filename = ctx.request.body.filename;
    console.log(filehash+';'+filename);
    var getPath = "./download/"+filename; //这个传入的路径要是相对于module/ipfs.js的相对路径
    await ipfsFile.get(filehash,getPath).then((mes)=>{
      console.log(mes);
    })
    let temp = {};
    temp.success= 'ok';
    temp.data = '下载数据成功';
    ctx.body = JSON.stringify(temp);
  })
module.exports = router