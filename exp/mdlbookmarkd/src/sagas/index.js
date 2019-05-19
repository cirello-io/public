import { put, takeLatest, takeLeading, all } from 'redux-saga/effects';
import config from '../config'

var cfg = config()

function* initialDataload() {
  const json = yield fetch(cfg.http + '/state', {
    credentials: 'same-origin'
  }).then(response => response.json())
  yield put({ type: 'INITIAL_LOAD_COMPLETE', bookmarks: json });
}

function* deleteBookmark(action) {
  yield fetch(cfg.http + '/deleteBookmark', {
    method: 'POST',
    body: JSON.stringify({ id: action.card.id }),
    credentials: 'same-origin'
  }).then(response => response.json()).catch((e) => {
    console.log('cannot delete bookmark:', e)
  })
}

function* linkWatcher() {
  yield takeLatest('INITIAL_LOAD_START', initialDataload)
  yield takeLeading('DELETE_BOOKMARK', deleteBookmark)
}

function* fuzzySearch(action) {
  yield put({ type: 'FUZZY_SEARCH', fuzzySearch: action.fuzzySearch });
}

function* fuzzySearchWatcher() {
  yield takeLatest('TRIGGER_FUZZY_SEARCH', fuzzySearch)
}

export default function* rootSaga() {
  yield all({
    link: linkWatcher(),
    fuzzySearch: fuzzySearchWatcher()
  })
}