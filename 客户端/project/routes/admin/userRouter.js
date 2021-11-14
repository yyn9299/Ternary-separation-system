const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');

/* GET users listing. */
//查询所有用户列表
router.get('/',async (ctx,next) =>{

  let page = ctx.query.page;
  console.log('前端获取的当前页码:'+page);
  page = page ? page : 1;
  let sqlStr = "select id,username,password,email,phone,status,rolename "+
                "from user,role where user.roleid=role.roleid limit ?,5";
  let result = await sqlQuery(sqlStr,[(parseInt(page)-1)*5])
  console.log(result);
  var sql = 'select count(*) as usersnum from user';//查询总数
  let resucnts = await sqlQuery(sql);
  var counts = resucnts[0].usersnum;
  var pages = Math.ceil(counts/5);
  let options = {
      userlist:Array.from(result),
      counts:counts,
      page:page,
      pages:pages,
  }
  //千万不要忘了异步 异步  await
  await ctx.render('admin/userList',options);

});

//删除单个用户
router.post('/delUserByID',async (ctx,next) =>{
  console.log(ctx.request.body)
  let id = ctx.request.body.id;
  let sqlStr = "delete from user where id=?";
  let result = await sqlQuery(sqlStr,[id]);
  console.log('删除成功！'+id);
  let temp = {};
  temp.success= 'ok';
  temp.data = '删除用户成功';
  ctx.body = JSON.stringify(temp);
  console.log(result);
  //ctx.body = '删除用户!';
})

//一次性删除多个用户--批量删除
router.post('/delAllUser',async (ctx,next)=>{
  console.log(ctx.request.body)
  let dellist = ctx.request.body.dellist;
  console.log(dellist)
  dellist.forEach(async (item,i)=>{
      let sqlStr = "delete from user where id = ?";
      await sqlQuery(sqlStr,item);
  })
  let temp = {};
  temp.success= 'ok';
  temp.data = '删除用户成功';
  ctx.body = JSON.stringify(temp);
  
})

//停用用户
router.post('/stopUser',async (ctx,next) =>{
  console.log(ctx.request.body)
  let id = ctx.request.body.id;
  let sqlStr = "update user set status=0 where id=?";  //将状态设置为0表示停用
  let result = await sqlQuery(sqlStr,[id]);
  console.log('停用用户成功！'+id);
  let temp = {};
  temp.success= 'ok';
  temp.data = '停用用户成功';
  ctx.body = JSON.stringify(temp);
  console.log(result);
})

//开启用户
router.post('/openUser',async (ctx,next) =>{
  console.log(ctx.request.body)
  let id = ctx.request.body.id;
  let sqlStr = "update user set status=1 where id=?";  //将状态设置为0表示停用
  let result = await sqlQuery(sqlStr,[id]);
  console.log('开启用户成功！'+id);
  let temp = {};
  temp.success= 'ok';
  temp.data = '开启用户成功';
  ctx.body = JSON.stringify(temp);
  console.log(result);
})
module.exports = router;
