export const WEBCAM_SOURCE = 'webcam://local';
export function cameraSourceKind(camera) { return camera?.stream_url === WEBCAM_SOURCE ? 'webcam' : 'network'; }
export function cameraAddress(camera) { return cameraSourceKind(camera) === 'webcam' ? 'Веб-камера этого устройства' : `${camera.ip}:${camera.port}`; }
export function parseNetworkURL(value) {
  const raw = value.trim();
  if (!raw) throw new Error('Укажите адрес потока.');
  let source;
  try { source = new URL(/^[a-z][a-z0-9+.-]*:\/\//i.test(raw) ? raw : `rtsp://${raw}`); }
  catch { throw new Error('Не удалось разобрать адрес потока.'); }
  if (!['rtsp:', 'rtsps:', 'http:', 'https:'].includes(source.protocol) || !source.hostname || source.hash) {
    throw new Error('Нужен адрес RTSP или HTTP(S).');
  }
  return source;
}
export function buildCameraPayload(fields) {
  const camera = { comment: fields.comment.trim(), is_enabled: Boolean(fields.enabled) };
  if (fields.kind === 'webcam') return { ...camera, ip: 'localhost', port: 80, stream_url: WEBCAM_SOURCE };
  const source = parseNetworkURL(fields.url);
  if (fields.login || fields.password) {
    if (!fields.login) throw new Error('Укажите логин камеры.');
    source.username = fields.login;
    source.password = fields.password;
  }
  const defaultPort = source.protocol === 'https:' ? 443 : 80;
  const port = fields.port === '' ? (['http:', 'https:'].includes(source.protocol) && source.port ? Number(source.port) : defaultPort) : Number(fields.port);
  if (!Number.isInteger(port) || port < 1 || port > 65535) throw new Error('HTTP-порт должен быть от 1 до 65535.');
  return { ...camera, ip: fields.ip.trim() || source.hostname, port, stream_url: source.toString() };
}
export function editableSource(camera) {
  const fields = { kind: cameraSourceKind(camera), url: camera?.stream_url || '', login: '', password: '' };
  if (fields.kind === 'network' && fields.url) {
    try {
      const source = parseNetworkURL(fields.url);
      fields.login = decodeURIComponent(source.username);
      fields.password = decodeURIComponent(source.password);
      source.username = ''; source.password = '';
      fields.url = source.toString();
    } catch { /* Keep an older address editable even if it cannot be played. */ }
  }
  return fields;
}
export function savedDevice(id) {
  try { return localStorage.getItem(`camera-device-${id}`) || ''; } catch { return ''; }
}
export function rememberDevice(id, device) {
  try {
    if (device) localStorage.setItem(`camera-device-${id}`, device);
    else localStorage.removeItem(`camera-device-${id}`);
  } catch { /* Playback can use the default camera without persistent storage. */ }
}
