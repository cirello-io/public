// Copyright 2019 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { folders } from '../helpers/folders'

const initialState = {
  loaded: false,
  bookmarks: [],
  filteredBookmarks: [],
  bookmark: null,
  folder: folders[0],
  selectedIndex: 0,
  fuzzySearch: ''
}

function reducer(state = initialState, action) {
  switch (action.type) {
    case 'SELECT_BOOKMARK_FOLDER':
      return {
        ...state,
        selectedIndex: action.selectedIndex,
        folder: folders[action.selectedIndex],
        filteredBookmarks: folders[action.selectedIndex].filter(state.bookmarks)
      }
    case 'INITIAL_LOAD_COMPLETE': {
      let ret = {
        ...state,
        loaded: true,
        bookmarks: (
          state.loaded
            ? []
            : state.bookmarks.slice()
        ).filter((v) => v.id !== state.bookmark.id).concat(action.bookmarks)
      }
      ret.filteredBookmarks = folders[state.selectedIndex].filter(ret.bookmarks)
      return ret
    }
    case 'FUZZY_SEARCH': {
      if (action.fuzzySearch === '') {
        return {
          ...state,
          filteredBookmarks: state.folder.filter(state.bookmarks)
        }
      }
      var fuzzySearch = action.fuzzySearch.toLowerCase()
      return {
        ...state,
        filteredBookmarks: state.folder.filter(state.bookmarks).filter((v) => {
          return fuzzyMatch(v.url.toLowerCase(), fuzzySearch) ||
            fuzzyMatch(v.title.toLowerCase(), fuzzySearch)
        })
      }
    }
    case 'BOOKMARK_ADDED': {
      let bookmarks = state.bookmarks.slice()
      bookmarks.unshift(action.bookmark)
      return {
        ...state,
        bookmark: action.bookmark,
        bookmarks
      }
    }
    case 'BOOKMARK_UPDATED': {
      var bookmarks = state.bookmarks.slice().map((v) => {
        if (v.id === action.bookmark.id) {
          return action.bookmark
        }
        return v
      })
      return {
        ...state,
        bookmark: action.bookmark,
        bookmarks
      }
    }
    case 'BOOKMARK_DELETED': {
      return {
        ...state,
        bookmarks: state.bookmarks.slice().filter((v) => v.id !== action.id),
        filteredBookmarks: state.filteredBookmarks.slice().filter((v) => v.id !== action.id)
      }
    }
    case 'BOOKMARK':
      return {
        ...state,
        bookmark: action.bookmark
      }
    default:
      return state
  }
}

export default reducer

// distilled from https://gist.github.com/mdwheele/7171422
function fuzzyMatch(haystack, needle) {
  var caret = 0
  for (var i = 0; i < needle.length; i++) {
    var c = needle[i]
    if (c === ' ') {
      continue
    }
    caret = haystack.indexOf(c, caret)
    if (caret === -1) {
      return false
    }
    caret++
  }
  return true
}