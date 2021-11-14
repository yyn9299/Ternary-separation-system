const ipfsAPI = require('ipfs-api');
const ipfs = ipfsAPI({host: 'localhost', port: '5001', protocol: 'http'});
const fs  = require('fs');

exports.add = (addPath) =>{
    return new Promise((resolve,reject)=>{
        try {
            let buffer = fs.readFileSync(addPath);
            ipfs.add(buffer, function (err, files) {
                if (err || typeof files == "undefined") {
                    reject(err);
                } else {
                    resolve(files[0].hash);
                }
            })
        }catch(ex) {
            reject(ex);
        }
    })
}
exports.get = (hash,getPath) =>{
    return new Promise((resolve,reject)=>{
        try{
            ipfs.get(hash,function (err,files) {
                if (err || typeof files == "undefined") {
                    reject(err);
                }else{
                    fs.writeFileSync(getPath,files[0].content);
                    resolve('下载成功！');                   
                }
            })
        }catch (ex){
            reject(ex);
        }
    });
}