import { put, takeLatest, all } from 'redux-saga/effects';

function* addLink(action) {
  yield put({ type: "LINK_ADDED", link: action.url });
}

function* linkWatcher() {
  yield takeLatest('ADD_LINK', addLink)
}

export default function* rootSaga() {
  yield all({
    link:linkWatcher()
  })
}
