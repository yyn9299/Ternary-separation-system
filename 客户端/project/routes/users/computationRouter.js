const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');

//跳转到添加算力页面
router.get('/addComputation',async (ctx)=>{
    console.log('添加算力页面')
    await ctx.render('users/computation/addComputation');
})

//查询所有算力列表
router.get('/GetAllComps',async (ctx)=>{
    let uname = ctx.session.username;
    const contract = getcontract(uname);
    let result = await contract.evaluateTransaction('GetAllComps');
    let computationlist = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
    var counts = computationlist.length;
    //分页展示 一页展示5条记录
    let page = ctx.query.page;
    page = page ? page : 1;
    let computationlist1 = computationlist.slice((page-1)*5,page*5);
    let options = {
        computationlist:computationlist1,
        counts:counts,
        page:page,
    }
    await ctx.render('users/computation/computationList',options);
})

//添加算力信息
router.post('/addComputation',async (ctx)=>{
    let uname = ctx.session.username;
    let cid = ctx.request.body.cid;
    let cpu = ctx.request.body.cpu;
    let cmemory = ctx.request.body.cmemory;
    let coin = parseInt(ctx.request.body.coin);
    const contract = getcontract(uname);
    let result = await contract.submitTransaction('CreateData',cid,cpu,cmemory,coin,uname);
    if(`${result}` !== ''){
      let temp = {};
      temp.success= 'ok';
      temp.data = '添加算力成功';
      ctx.body = JSON.stringify(temp);
    }else{
      console.log('添加算力失败');
    }
})

//根据id查询算力
router.post('/SelectCompImByCid',async (ctx,next) =>{
    console.log(ctx.request.body)
    let cid = ctx.request.body.cid;
    let username = ctx.session.username;
    const contract = getcontract(username);
    let result = await contract.evaluateTransaction('SelectCompImByCid',cid);
    let computation = JSON.parse(prettyJSONString(result.toString()));  //解析成数组
    console.log('查询算力成功！'+computation);
    let temp = {};
    temp.success= 'ok';
    temp.data = '查询算力成功';
    ctx.body = JSON.stringify(temp);
  })
//修改算力信息
router.post('/UpdateCompByCid',async (ctx,next) =>{
    console.log(ctx.request.body)
    let username = ctx.session.username;
    let cid = ctx.request.body.cid;
    let cpu = ctx.request.body.cpu;
    let cmemory = ctx.request.body.cmemory;
    let coin = parseInt(ctx.request.body.coin);
    const contract = getcontract(username);
    await contract.submitTransaction('UpdateCompByCid',cid,cpu,cmemory,coin);
    console.log('修改成功！'+cid);
    let temp = {};
    temp.success= 'ok';
    temp.data = '修改算力成功';
    ctx.body = JSON.stringify(temp);
  })

module.exports = router