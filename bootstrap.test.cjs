const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const crypto = require('node:crypto');
const vm = require('node:vm');
const directory = fs.mkdtempSync(path.join(os.tmpdir(), 'chuncord-bootstrap-'));
const target = path.join(directory, 'chuncord.asar');
const candidate = Buffer.from('verified new archive');
const receipt = {version:1,size:candidate.length,sha256:crypto.createHash('sha256').update(candidate).digest('hex'),revision:'a'.repeat(40)};
const source = fs.readFileSync(path.join(__dirname,'bootstrap.cjs'),'utf8').replace('__CHUNCORD_PATH__',JSON.stringify(target));
let loaded;
function run(filesystem=fs) { vm.runInNewContext(source,{require(name){if(name==='original-fs')return filesystem;if(name==='crypto')return crypto;assert.equal(name,target);loaded=fs.readFileSync(target,'utf8');}}); }
try {
 fs.writeFileSync(target,'original');run();assert.equal(loaded,'original');
 fs.writeFileSync(target+'.pending',candidate);fs.writeFileSync(target+'.pending.json',JSON.stringify({...receipt,sha256:'0'.repeat(64)}));run();assert.equal(loaded,'original');
 fs.writeFileSync(target+'.pending.json',JSON.stringify(receipt));
 run({...fs,renameSync(){throw Object.assign(new Error('Windows file is in use'),{code:'EPERM'});}});
 assert.equal(loaded,'original');assert(fs.existsSync(target+'.pending'));
 run();assert.equal(loaded,candidate.toString(),fs.existsSync(target+'.update-error.log')?fs.readFileSync(target+'.update-error.log','utf8'):'no error log');assert.equal(fs.readFileSync(target+'.previous','utf8'),'original');assert(!fs.existsSync(target+'.pending.json'));
 console.log('Bootstrap verified commit, corrupt rejection, locked target and old archive preservation PASS');
} finally { for(const name of fs.readdirSync(directory))fs.unlinkSync(path.join(directory,name));fs.rmdirSync(directory); }
