const Koa = require('koa')
const app = new Koa()
const views = require('koa-views')
const json = require('koa-json')
const onerror = require('koa-onerror')
const bodyparser = require('koa-bodyparser')
const logger = require('koa-logger')
const session = require('koa-session');  //引入session
//引入路由模块
const indexRouter = require('./routes/index')
const usersRouter = require('./routes/users/usersRouter')
//登录模块
var loginRouter = require('./routes/login/loginRouter')
//后台模块
var adminRouter = require('./routes/admin/adminRouter')
//客户端用户模块
const koaBody = require('koa-body');
app.use(koaBody({
    multipart: true,
    formidable: {
        maxFileSize: 100*200*1024*1024    // 设置上传文件大小最大限制，默认2M
    }
}));


// error handler
onerror(app)

// 处理表单数据
app.use(bodyparser({
  enableTypes:['json', 'form', 'text']
}))
app.use(json())
app.use(logger())
app.use(require('koa-static')(__dirname + '/public'))

app.use(views(__dirname + '/views', {
  extension: 'ejs'
}))

//配置session
app.keys = ['some secret hurr'];
const CONFIG = {
    key: 'koa.sess',  //默认
    maxAge: 864000, //cookie过期时间 需要设置
    httpOnly: true, //true表示只有服务器端可以获取到cookie
    signed: true,  //默认
    rolling: false, 
    renew: true, //需要设置修改
    sameSite: null, 
  };
app.use(session(CONFIG, app));

// logger
app.use(async (ctx, next) => {
  const start = new Date()
  await next()
  const ms = new Date() - start
  console.log(`${ctx.method} ${ctx.url} - ${ms}ms`)
})

// routes
app.use(indexRouter.routes(),indexRouter.allowedMethods());
app.use(usersRouter.routes(),usersRouter.allowedMethods());
app.use(adminRouter.routes(),adminRouter.allowedMethods());
app.use(loginRouter.routes(),loginRouter.allowedMethods());


// error-handling
app.on('error', (err, ctx) => {
  console.error('server error', err, ctx)
});

module.exports = app
