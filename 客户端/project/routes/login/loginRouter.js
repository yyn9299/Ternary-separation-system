const router = require('koa-router')()
var crypto = require('crypto');
var sqlQuery = require('../../module/lcMysql');

function jiami(str){
    let salt = "fjdsoigijasoigjasdiodgjasdiogjoasid"
    let obj = crypto.createHash('md5')
    str = salt+str;
    obj.update(str)
    return obj.digest('hex')
}

router.prefix('/rl')
router.get('/login', async (ctx, next) => {
  await ctx.render('login/login.ejs');
})

router.get('/register', async (ctx, next) => {
    await ctx.render('login/register.ejs');
})

//处理post表单提交的数据
//登录处理 
router.post('/login',async (ctx,next) =>{
    let username = ctx.request.body.username;
    let password = ctx.request.body.password;
    // let pawmd5 = jiami(password);
    // console.log(pawmd5);
    let sqlStr = 'select * from user where username=? and password=?';
    let result = await sqlQuery(sqlStr,[username,password]);
    if(result.length == 0){
        await ctx.render('info/info',{
            title:"登陆失败",
            content:"用户或密码错误",
            href:"/rl/login",
            hrefTxt:"登陆页"
        })
    }else if(result[0].roleid == 1){
        ctx.session.username = username;
        await ctx.render('info/info',{
            title:"登陆成功",
            content:"立即跳转至后台页面",
            href:"/admin",
            hrefTxt:"后台"
        })
    }else if(result[0].roleid == 2){
        ctx.session.username = username;
        await ctx.render('info/info',{
            title:"登陆成功",
            content:"普通用户登录成功，立即跳转至用户页面",
            href:"/users",
            hrefTxt:"用户页面"
        })
    }
    //console.log(username);
})
//注册处理
router.post('/register',async (ctx)=>{
    //获取username和密码
    let username = ctx.request.body.username;
    let password = ctx.request.body.password;
    let roleid = ctx.request.body.roleid;
    console.log(password);
    console.log(roleid);
    //判断用户是否存在,如果没有用户才进行插入
    let sqlStr = "select * from user where username = ?";
    let result = await sqlQuery(sqlStr,[username]);
    // let pawmd5 = jiami(password);
    // console.log(pawmd5);
    if(result.length!=0){
        //告知此用户名已存在，请直接登陆或者找寻密码
        await ctx.render('info/info',{
            title:"注册失败",
            content:"用户已存在",
            href:"/rl/register",
            hrefTxt:"注册页"
        })
    }else{
        //告知注册成功
        sqlStr = "insert into user(id,username,password,roleid) values (default,?,?,?)"
        await sqlQuery(sqlStr,[username,password,roleid])
        await ctx.render('info/info',{
            title:"注册成功",
            content:"注册成功，即将进入登陆页",
            href:"/rl/login",
            hrefTxt:"登录页"
        })
    }
})

//退出登录
router.get('/loginout',async (ctx,next) =>{
    ctx.session.username = undefined;
    await ctx.render('info/info',{
        title:"登出成功",
        content:"退出登录成功，即将进入登陆页",
        href:"/rl/login",
        hrefTxt:"登录页"
    })

})
module.exports = router
