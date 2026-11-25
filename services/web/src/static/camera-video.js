import { cameraSourceKind, savedDevice } from './camera-source.mjs';
const $ = (id) => document.getElementById(id);
let selected = null;
let stream = null;
let generation = 0;
let imageTimer = null;

function mediaError(error) {
  if (error.name === 'NotAllowedError' || error.name === 'SecurityError') return 'Нет доступа к веб-камере. Разрешите его в настройках браузера.';
  if (error.name === 'NotFoundError') return 'Веб-камера не найдена. Подключите её к компьютеру.';
  if (error.name === 'NotReadableError') return 'Камера занята другой программой или недоступна.';
  if (error.name === 'OverconstrainedError') return 'Выбранная веб-камера недоступна. Выберите другую в настройках камеры.';
  return 'Не удалось включить камеру. Попробуйте ещё раз.';
}
function cameraAccessSupported() {
  return window.isSecureContext && navigator.mediaDevices?.getUserMedia;
}
function message(text) { $('live-status').textContent = text; }
function canPlay() { return selected?.is_enabled && Boolean(selected.stream_url); }
export function stopLivePreview() {
  generation++;
  clearTimeout(imageTimer);
  imageTimer = null;
  stream?.getTracks().forEach((track) => track.stop());
  stream = null;
  const video = $('local-video');
  video.pause(); video.srcObject = null; video.hidden = true;
  const image = $('network-video');
  image.onload = null; image.onerror = null; image.removeAttribute('src'); image.hidden = true;
  $('live-placeholder').hidden = false;
  $('start-live').disabled = !canPlay();
  $('stop-live').disabled = true;
}
export function renderLivePreview(camera) {
  stopLivePreview();
  selected = camera || null;
  $('live-panel').hidden = !camera;
  $('start-live').disabled = !canPlay();
  if (!camera) return;
  if (!camera.is_enabled) message('Камера отключена. Включите её в настройках.');
  else if (cameraSourceKind(camera) === 'webcam') message('Камера этого компьютера. Видео не отправляется на сервер.');
  else if (!camera.stream_url) message('Добавьте адрес видеопотока в настройках камеры.');
  else message('Просмотр по запросу. Без звука и записи.');
}
async function startLivePreview() {
  if (!canPlay()) return;
  stopLivePreview();
  const current = generation;
  const camera = selected;
  $('start-live').disabled = true;
  $('stop-live').disabled = false;
  message('Подключаем камеру…');
  if (cameraSourceKind(camera) === 'webcam') {
    if (!cameraAccessSupported()) {
      stopLivePreview();
      message('Для веб-камеры откройте панель на localhost или через HTTPS в Chrome или Edge.');
      return;
    }
    try {
      const deviceId = savedDevice(camera.id);
      const media = await navigator.mediaDevices.getUserMedia({ video: deviceId ? { deviceId: { exact: deviceId } } : true, audio: false });
      if (current !== generation) { media.getTracks().forEach((track) => track.stop()); return; }
      stream = media;
      const video = $('local-video');
      video.srcObject = media; video.hidden = false;
      $('live-placeholder').hidden = true;
      await video.play();
      if (current !== generation) return;
      message('Веб-камера включена. Без звука и записи.');
      for (const track of media.getVideoTracks()) track.addEventListener('ended', () => {
        if (stream === media) { stopLivePreview(); message('Веб-камера отключилась.'); }
      }, { once: true });
    } catch (error) {
      if (current !== generation) return;
      stopLivePreview(); message(mediaError(error));
    }
    return;
  }
  const image = $('network-video');
  image.onload = () => {
    if (current !== generation) return;
    clearTimeout(imageTimer);
    image.hidden = false;
    $('live-placeholder').hidden = true;
    message('Поток подключён. Без звука и записи.');
  };
  image.onerror = () => {
    if (current !== generation) return;
    stopLivePreview();
    message('Не удалось открыть поток. Проверьте адрес, логин и доступность камеры.');
  };
  imageTimer = setTimeout(() => {
    if (current !== generation || image.naturalWidth) return;
    stopLivePreview(); message('Камера не ответила. Проверьте адрес и подключение.');
  }, 18000);
  image.src = `/live/${camera.id}?session=${Date.now()}`;
}
export async function populateWebcams(select) {
  if (!cameraAccessSupported()) throw new Error('Веб-камера доступна на localhost или HTTPS в Chrome или Edge.');
  let media;
  try { media = await navigator.mediaDevices.getUserMedia({ video: true, audio: false }); }
  catch (error) { throw new Error(mediaError(error)); }
  try {
    const devices = (await navigator.mediaDevices.enumerateDevices()).filter((device) => device.kind === 'videoinput');
    if (!devices.length) throw new Error('Веб-камеры не найдены.');
    const previous = select.value;
    select.replaceChildren(new Option('Камера по умолчанию', ''));
    devices.forEach((device, index) => select.add(new Option(device.label || `Веб-камера ${index + 1}`, device.deviceId)));
    if (devices.some((device) => device.deviceId === previous)) select.value = previous;
  } finally { media.getTracks().forEach((track) => track.stop()); }
}
$('start-live').addEventListener('click', startLivePreview);
$('stop-live').addEventListener('click', () => { stopLivePreview(); message('Просмотр остановлен.'); });
window.addEventListener('pagehide', stopLivePreview);
document.addEventListener('visibilitychange', () => {
  if (document.hidden) { stopLivePreview(); if (selected) message('Просмотр остановлен.'); }
});
