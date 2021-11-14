const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');
var taskRouter = require('./taskRouter');
var uploadRouter = require('./uploadRouter');
var datesetRouter = require('./datesetRouter')
var computationRouter = require('./computationRouter')

//判断是否符合条件进入后台中间件
async function permisson(ctx,next){
  if(ctx.session.username==undefined){
      //尚未登陆，返回至登录页
     await ctx.render('info/info',{
          title:"尚未登陆",
          content:"请重新登陆,即将进入登陆页",
          href:"/rl/login",
          hrefTxt:"登录页"
      })
  }else{
      //正常进入
      await next();
  }
}

//跳转到用户首页
router.get('/users',permisson,async (ctx, next) => {
  await ctx.render('users/index.ejs',{
    username:ctx.session.username
  })
})

//用户点击个人信息
router.get('/users/selfinfo',async (ctx,next) => {
  //首先获取登录用户的用户名 来进行查询该用户的所有信息
  let username = ctx.session.username;
  let sqlStr = 'select *from user where username=?';
  let result = await sqlQuery(sqlStr,[username]);
  let user = result[0];
  var options = {user};
  await ctx.render('users/selfinfo.ejs',options);
})

//修改个人信息
router.post('/users/selfinfo',async (ctx)=>{
  let username = ctx.session.username;
  let password = ctx.request.body.password;
  let email = ctx.request.body.email;
  let phone = ctx.request.body.phone;
  let sqlStr = 'update user set password=?,email=?,phone=? where username=?';
  let data = [password,email,phone,username];
  let result = await sqlQuery(sqlStr,data);
  console.log('个人信息更新成功');
  let temp = {};
  temp.success= 'ok';
  temp.data = '修改信息成功';
  ctx.body = JSON.stringify(temp);
  console.log(result);
})

//将任务路由挂载在这里
router.use('/users/task',taskRouter.routes(),taskRouter.allowedMethods());
router.use('/users/upload',uploadRouter.routes(),uploadRouter.allowedMethods());
router.use('/users/date',datesetRouter.routes(),datesetRouter.allowedMethods());
router.use('/users/computation',computationRouter.routes(),computationRouter.allowedMethods());

module.exports = router
