const router = require('koa-router')();
var sqlQuery = require('../../module/lcMysql');
var fs = require('fs')
var  path = require('path')
//var ipfs = require('../../module/ipfs')
const ipfsAPI = require('ipfs-api');
const ipfs = ipfsAPI({host: 'localhost', port: '5001', protocol: 'http'});

//跳转到上传数据集页面
router.get('/uploadfile',async (ctx)=>{
    console.log('上传数据页面')
    await ctx.render('users/file/uploadfile.ejs');
})

router.post('/uploadfile', async (ctx, next) => {
    // 上传单个文件
    const file = ctx.request.files.file; // 获取上传文件
    const filename = ctx.request.body.filename;

    // 用异步方法读取流  将其转为buff 把buff传到ipfs中去
    const buff = fs.readFileSync(file.path);
    //ipfs.add(buff,)
    var hash = '';
    ipfs.add(buff,async (err,result)=>{
        if(err) throw err;
        console.log(result);
        hash = result[0].hash;
        console.log(hash);
    })
    var sqlStr = 'insert into file values(default,?,?)';
    var result = await sqlQuery(sqlStr,[filename,hash]);
    console.log(result);

    return ctx.body = "上传成功！";
  });

module.exports = router