import { put, takeLatest, all } from 'redux-saga/effects';
import config from '../config'

var cfg = config()

function* initialDataload() {
  const json = yield fetch(cfg.http + '/state', {
    credentials: 'same-origin'
  }).then(response => response.json())
  yield put({ type: 'INITIAL_LOAD_COMPLETE', bookmarks: json });
}

function* linkWatcher() {
  yield takeLatest('INITIAL_LOAD_START', initialDataload)
}

export default function* rootSaga() {
  yield all({
    link: linkWatcher()
  })
}