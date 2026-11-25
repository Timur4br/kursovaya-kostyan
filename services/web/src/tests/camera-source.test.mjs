import test from 'node:test';
import assert from 'node:assert/strict';
import {buildCameraPayload, editableSource, cameraAddress, WEBCAM_SOURCE, parseNetworkURL} from '../static/camera-source.mjs';
const base = {kind:'network',comment:' Вход ',enabled:true,url:'rtsp://192.168.1.10:554/live',ip:'',port:'',login:'',password:''};
test('network source derives the host without confusing RTSP and HTTP ports',()=>{
  const camera=buildCameraPayload(base);
  assert.equal(camera.ip,'192.168.1.10'); assert.equal(camera.port,80); assert.equal(camera.comment,'Вход');
  assert.equal(camera.stream_url,base.url);
  assert.equal(buildCameraPayload({...base,url:'http://camera.local:8080/video'}).port,8080);
  assert.equal(buildCameraPayload({...base,url:'https://camera.local/video'}).port,443);
  assert.equal(parseNetworkURL('camera.local:554/live').protocol,'rtsp:');
});
test('credentials round trip without remaining in the visible stream address',()=>{
  const camera=buildCameraPayload({...base,login:'user@camera',password:'a:b@c /#'});
  const fields=editableSource(camera);
  assert.equal(fields.login,'user@camera'); assert.equal(fields.password,'a:b@c /#');
  assert.equal(fields.url,base.url); assert.equal(fields.kind,'network');
  assert.equal(buildCameraPayload({...base,url:fields.url,login:fields.login,password:fields.password}).stream_url,camera.stream_url);
});
test('webcam uses a browser-local marker and no device identifier in API data',()=>{
  const camera=buildCameraPayload({...base,kind:'webcam',enabled:false});
  assert.equal(camera.stream_url,WEBCAM_SOURCE); assert.equal(camera.ip,'localhost'); assert.equal(camera.port,80);
  assert.equal(camera.is_enabled,false); assert.equal(editableSource(camera).kind,'webcam');
  assert.equal(cameraAddress(camera),'Веб-камера этого устройства');
});
test('rejects unsupported sources and invalid ports',()=>{
  for(const url of ['', 'file:///etc/passwd','exec:command','ftp://camera/video','https://camera/video#fragment']) {
    assert.throws(()=>buildCameraPayload({...base,url}));
  }
  for(const port of ['0','65536','not-a-port','1.5']) assert.throws(()=>buildCameraPayload({...base,port}));
  assert.throws(()=>buildCameraPayload({...base,password:'secret',login:''}));
});
