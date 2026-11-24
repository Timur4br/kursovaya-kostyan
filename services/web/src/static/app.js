const $ = (id) => document.getElementById(id);
const CAMERA_PAGE_SIZE = 8;
const EVENT_PAGE_SIZE = 15;
const state = {
  cameraPage: 0, cameraCount: 0, cameras: [], selected: null,
  eventPage: 0, eventCount: 0, filters: {}, cameraRequest: null,
  eventRequest: null, cameraBusy: false, eventBusy: false, saving: false,
  objectNames: new Map(),
};
const dateFormat = new Intl.DateTimeFormat('ru-RU', { day: '2-digit', month: 'short', year: 'numeric' });
const timeFormat = new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
const percentFormat = new Intl.NumberFormat('ru-RU', { style: 'percent', maximumFractionDigits: 1 });

function node(tag, className, text) {
  const element = document.createElement(tag);
  if (className) element.className = className;
  if (text !== undefined) element.textContent = text;
  return element;
}
function notice(id, text = '') {
  $(id).textContent = text;
  $(id).hidden = !text;
}
function empty(container, title, description) {
  const wrapper = node('div', 'empty-state');
  wrapper.append(node('h3', '', title));
  if (description) wrapper.append(node('p', '', description));
  container.replaceChildren(wrapper);
}
function cameraName(camera) { return camera.comment?.trim() || `Камера № ${camera.id}`; }
function endpoint(entity, action, params = {}) {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== '' && value !== null && value !== undefined) query.set(key, String(value));
  });
  return `/api/${entity}/${action}?${query}`;
}
async function api(path, options = {}) {
  let response;
  try { response = await fetch(path, options); }
  catch (error) {
    if (error.name === 'AbortError') throw error;
    throw new Error('Нет соединения с сервером.');
  }
  if (!response.ok) {
    if (response.status === 400) throw new Error('Проверьте введённые значения.');
    if (response.status === 404) throw new Error('Запись не найдена. Обновите список.');
    if (response.status === 502 || response.status === 503 || response.status === 504) {
      throw new Error('Сервер недоступен. Попробуйте позже.');
    }
    throw new Error('Не удалось выполнить запрос. Попробуйте ещё раз.');
  }
  return response.json();
}
function updateCameraPager() {
  const pages = Math.max(1, Math.ceil(state.cameraCount / CAMERA_PAGE_SIZE));
  $('camera-prev').disabled = state.cameraBusy || state.cameraPage === 0;
  $('camera-next').disabled = state.cameraBusy || state.cameraPage + 1 >= pages;
  $('camera-page').textContent = state.cameraCount ? `${state.cameraPage + 1} / ${pages}` : 'Нет камер';
}
function updateEventPager() {
  $('event-prev').disabled = state.eventBusy || !state.selected || state.eventPage === 0;
  $('event-next').disabled = state.eventBusy || !state.selected || (state.eventPage + 1) * EVENT_PAGE_SIZE >= state.eventCount;
  if (state.selected && !state.eventBusy) {
    const start = state.eventCount ? state.eventPage * EVENT_PAGE_SIZE + 1 : 0;
    const end = Math.min((state.eventPage + 1) * EVENT_PAGE_SIZE, state.eventCount);
    $('event-page').textContent = `${start}–${end} из ${state.eventCount} · Местное время`;
  }
}
function renderCameras() {
  const list = $('camera-list');
  list.replaceChildren();
  if (!state.cameras.length) {
    empty(list, 'Камер пока нет', $('camera-status').value ? 'Попробуйте выбрать все камеры.' : 'Добавьте первую камеру.');
    return;
  }
  for (const camera of state.cameras) {
    const selected = state.selected?.id === camera.id;
    const card = node('button', `camera-card${selected ? ' selected' : ''}`);
    card.type = 'button';
    card.setAttribute('aria-pressed', String(selected));
    const info = node('span', 'camera-info');
    info.append(node('span', 'camera-name', cameraName(camera)),
      node('span', 'camera-address', `${camera.ip}:${camera.port}`),
      node('span', `camera-state${camera.is_enabled ? ' enabled' : ''}`, camera.is_enabled ? 'Включена' : 'Отключена'));
    card.append(info);
    card.addEventListener('click', () => selectCamera(camera));
    list.append(card);
  }
}
function selectCamera(camera) {
  state.selected = camera || null;
  state.eventPage = 0;
  renderCameras();
  $('apply-filter').disabled = !camera;
  $('reset-filter').disabled = !camera;
  if (!camera) {
    state.eventRequest?.abort();
    state.eventBusy = false;
    state.eventCount = 0;
    $('events-content').setAttribute('aria-busy', 'false');
    $('events-heading').textContent = 'События';
    $('selected-camera-detail').textContent = 'Выберите камеру';
    $('event-page').textContent = 'Местное время';
    notice('event-error');
    empty($('events-content'), 'Выберите камеру', 'Здесь будут её события.');
    updateEventPager();
    return;
  }
  $('events-heading').textContent = cameraName(camera);
  $('selected-camera-detail').textContent = `${camera.ip}:${camera.port}`;
  loadEvents();
}
async function loadCameras({ selectID = null, refreshEvents = false } = {}) {
  state.cameraRequest?.abort();
  const controller = new AbortController();
  state.cameraRequest = controller;
  const options = { signal: controller.signal };
  state.cameraBusy = true;
  notice('camera-error');
  $('camera-list').setAttribute('aria-busy', 'true');
  $('camera-list').replaceChildren(node('p', 'loading', 'Загрузка…'));
  updateCameraPager();
  try {
    const filters = { is_enabled: $('camera-status').value };
    const [count, cameras] = await Promise.all([
      api(endpoint('camera', 'count', filters), options),
      api(endpoint('camera', 'list', { ...filters, limit: CAMERA_PAGE_SIZE, offset: state.cameraPage * CAMERA_PAGE_SIZE }), options),
    ]);
    if (controller.signal.aborted) return;
    const lastPage = Math.max(0, Math.ceil(count / CAMERA_PAGE_SIZE) - 1);
    if (state.cameraPage > lastPage) {
      state.cameraPage = lastPage;
      return loadCameras({ selectID, refreshEvents });
    }
    state.cameraCount = count;
    state.cameras = cameras;
    $('camera-list-caption').textContent = String(count);
    const previousID = state.selected?.id;
    const selected = cameras.find((camera) => camera.id === (selectID ?? previousID)) || cameras[0];
    if (selected?.id !== previousID || selectID !== null || !selected) selectCamera(selected);
    else {
      state.selected = selected;
      renderCameras();
      if (refreshEvents) selectCamera(selected);
    }
  } catch (error) {
    if (controller.signal.aborted) return;
    notice('camera-error', error.message);
    empty($('camera-list'), 'Не удалось загрузить камеры', 'Попробуйте обновить список.');
  } finally {
    if (state.cameraRequest === controller) {
      state.cameraBusy = false;
      $('camera-list').setAttribute('aria-busy', 'false');
      updateCameraPager();
    }
  }
}
async function loadObjectNames(events, signal) {
  const ids = [...new Set(events.map((event) => event.object_id))].filter((id) => !state.objectNames.has(id));
  const results = await Promise.allSettled(ids.map((id) => api(`/api/object/${id}`, { signal })));
  results.forEach((result, index) => {
    if (result.status === 'fulfilled') state.objectNames.set(ids[index], result.value.type);
  });
}
function renderEvents(events) {
  const container = $('events-content');
  if (!events.length) {
    empty(container, 'Событий пока нет', Object.keys(state.filters).length ? 'За этот период ничего нет.' : '');
    return;
  }
  const scroll = node('div', 'table-scroll');
  const table = node('table');
  const caption = node('caption', 'sr-only', `События: ${cameraName(state.selected)}`);
  const head = node('thead');
  const heading = node('tr');
  ['Когда произошло', 'Событие', 'Объект', 'Уверенность'].forEach((label) => {
    const cell = node('th', '', label);
    cell.scope = 'col';
    heading.append(cell);
  });
  head.append(heading);
  const body = node('tbody');
  for (const event of events) {
    const row = node('tr');
    const dateCell = node('td');
    const date = new Date(event.created_at);
    if (Number.isNaN(date.getTime())) dateCell.textContent = 'Дата не указана';
    else {
      const time = node('time', 'event-time', timeFormat.format(date));
      time.dateTime = date.toISOString();
      time.title = date.toLocaleString('ru-RU');
      dateCell.append(time, node('span', 'event-date', dateFormat.format(date)));
    }
    const type = node('td');
    type.append(node('span', 'event-type', event.type));
    const object = node('td', 'object-name', state.objectNames.get(event.object_id) || `Объект № ${event.object_id}`);
    const confidence = node('td', 'confidence', Number.isFinite(event.confidence) ? percentFormat.format(event.confidence) : '—');
    row.append(dateCell, type, object, confidence);
    body.append(row);
  }
  table.append(caption, head, body);
  scroll.append(table);
  container.replaceChildren(scroll);
}
async function loadEvents() {
  if (!state.selected) return;
  state.eventRequest?.abort();
  const controller = new AbortController();
  state.eventRequest = controller;
  state.eventBusy = true;
  notice('event-error');
  $('events-content').setAttribute('aria-busy', 'true');
  $('events-content').replaceChildren(node('p', 'loading', 'Загрузка…'));
  updateEventPager();
  try {
    const filters = { ...state.filters, camera_id: state.selected.id };
    const options = { signal: controller.signal };
    const [count, events] = await Promise.all([
      api(endpoint('event', 'count', filters), options),
      api(endpoint('event', 'list', { ...filters, limit: EVENT_PAGE_SIZE, offset: state.eventPage * EVENT_PAGE_SIZE }), options),
    ]);
    if (controller.signal.aborted) return;
    const lastPage = Math.max(0, Math.ceil(count / EVENT_PAGE_SIZE) - 1);
    if (state.eventPage > lastPage) {
      state.eventPage = lastPage;
      return loadEvents();
    }
    await loadObjectNames(events, controller.signal);
    if (controller.signal.aborted) return;
    state.eventCount = count;
    renderEvents(events);
  } catch (error) {
    if (controller.signal.aborted) return;
    state.eventCount = 0;
    notice('event-error', error.message);
    empty($('events-content'), 'Не удалось загрузить события', 'Попробуйте ещё раз.');
  } finally {
    if (state.eventRequest === controller) {
      state.eventBusy = false;
      $('events-content').setAttribute('aria-busy', 'false');
      updateEventPager();
    }
  }
}
$('refresh').addEventListener('click', () => loadCameras({ refreshEvents: true }));
$('camera-status').addEventListener('change', () => { state.cameraPage = 0; loadCameras(); });
$('camera-prev').addEventListener('click', () => { if (state.cameraPage > 0) { state.cameraPage--; loadCameras(); } });
$('camera-next').addEventListener('click', () => { state.cameraPage++; loadCameras(); });
$('event-prev').addEventListener('click', () => { if (state.eventPage > 0) { state.eventPage--; loadEvents(); } });
$('event-next').addEventListener('click', () => { state.eventPage++; loadEvents(); });
$('event-filter').addEventListener('submit', (event) => {
  event.preventDefault();
  const from = $('created-from').value;
  const to = $('created-to').value;
  if (from && to && new Date(from) > new Date(to)) {
    notice('event-error', 'Начало периода должно быть раньше его окончания.');
    return;
  }
  state.filters = {};
  if (from) state.filters.created_from = new Date(from).toISOString();
  if (to) state.filters.created_to = new Date(to).toISOString();
  state.eventPage = 0;
  loadEvents();
});
$('reset-filter').addEventListener('click', () => {
  $('event-filter').reset();
  state.filters = {};
  state.eventPage = 0;
  loadEvents();
});
$('camera-dialog').addEventListener('cancel', (event) => { if (state.saving) event.preventDefault(); });
$('open-camera-dialog').addEventListener('click', () => { notice('form-error'); $('camera-dialog').showModal(); });
for (const id of ['close-dialog', 'cancel-dialog']) $(id).addEventListener('click', () => $('camera-dialog').close());
$('camera-form').addEventListener('submit', async (event) => {
  event.preventDefault();
  if (state.saving) return;
  const form = $('camera-form');
  const data = new FormData(form);
  const camera = {
    comment: data.get('comment').trim(), ip: data.get('ip').trim(),
    port: Number(data.get('port')), stream_url: data.get('stream_url').trim(),
    is_enabled: data.has('is_enabled'),
  };
  if (!camera.ip) { notice('form-error', 'Укажите IP-адрес или хост камеры.'); return; }
  state.saving = true;
  $('close-dialog').disabled = true;
  $('cancel-dialog').disabled = true;
  $('save-camera').disabled = true;
  $('save-camera').textContent = 'Сохраняем…';
  notice('form-error');
  try {
    const id = await api('/api/camera', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(camera) });
    $('camera-dialog').close();
    form.reset();
    $('toast').textContent = 'Камера добавлена';
    $('toast').hidden = false;
    clearTimeout(state.toastTimeout);
    state.toastTimeout = setTimeout(() => { $('toast').hidden = true; }, 4500);
    $('camera-status').value = '';
    try {
      const total = await api('/api/camera/count');
      state.cameraPage = Math.max(0, Math.ceil(total / CAMERA_PAGE_SIZE) - 1);
      await loadCameras({ selectID: id });
    } catch {
      notice('camera-error', 'Камера добавлена, но список не обновился. Нажмите кнопку обновления.');
    }
  } catch (error) { notice('form-error', error.message); }
  finally {
    state.saving = false;
    $('close-dialog').disabled = false;
    $('cancel-dialog').disabled = false;
    $('save-camera').disabled = false;
    $('save-camera').textContent = 'Добавить камеру';
  }
});
loadCameras();
