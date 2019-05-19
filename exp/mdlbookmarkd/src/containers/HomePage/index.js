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

import './style.scss'
import Button from '@material/react-button'
import Card, { CardPrimaryContent, CardActions, CardActionButtons, CardActionIcons } from "@material/react-card";
import Dialog, { DialogTitle, DialogContent, DialogFooter, DialogButton } from '@material/react-dialog'
import Fab from '@material/react-fab'
import MaterialIcon from '@material/react-material-icon'
import React from 'react'
import moment from 'moment'
import { Cell, Grid, Row } from '@material/react-layout-grid'
import { Headline6, Caption } from '@material/react-typography'
import { connect } from 'react-redux'
import { folders } from '../../helpers/folders'

class HomePage extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      visible: -1,
      fuzzySearch: '',
      delete: null
    }
    this.filterBy = this.filterBy.bind(this)
    this.deleteDialog = this.deleteDialog.bind(this)
    this.deleteAction = this.deleteAction.bind(this)
  }

  componentDidMount() {
    if (!this.dependenciesLoaded()) {
      this.props.dispatch({ type: 'INITIAL_LOAD_START' })
    }
  }

  dependenciesLoaded() {
    return this.props.bookmarks && this.props.bookmarks.loaded
  }

  filterBy(v) {
    this.setState({ fuzzySearch: v.toLowerCase() })
  }

  deleteDialog(e, card) {
    e.preventDefault()
    this.setState({ delete: card })
  }

  deleteAction() {
    this.props.dispatch({ type: 'DELETE_BOOKMARK', card: { ...this.state.delete } })
    this.setState({ delete: null })
  }

  markAsRead(e, id) {
    e.preventDefault()
    this.props.markBookmarkAsRead(id)
  }

  render() {
    if (!this.dependenciesLoaded()) {
      return (<Grid></Grid>)
    }

    const listing = this.props.filteredBookmarks

    return <div>
      {this.state.delete !== null
        ? <Dialog
          open
          onClose={(action) => {
            switch (action) {
              case 'delete':
                return this.deleteAction()
              default:
                this.setState({ delete: null })
            }
          }}>
          <DialogTitle>Delete Bookmark</DialogTitle>
          <DialogContent>
            Delete "{this.state.delete.title.trim() !== '' ? this.state.delete.title.trim() : this.state.delete.url}" ?
          </DialogContent>
          <DialogFooter>
            <DialogButton action='keep'>Keep</DialogButton>
            <DialogButton action='delete' isDefault>Delete</DialogButton>
          </DialogFooter>
        </Dialog>
        : null}

      <Grid key={'homePageRoot'}>
        <Row key={'homePageRow0'}>
          {listing.map((v) => {
            return <Cell columns={4} key={'bookmarkCard-cell-' + v.id}>
              <BookmarkCard
                key={'bookmarkCard' + v.id}
                card={v}
                markAsRead={(e) => { this.markAsRead(e, v.id) }}
                delete={(e) => { this.deleteDialog(e, v) }} />
            </Cell>
          })}
        </Row>
      </Grid>
      <Fab key={'addLink'} className='addNewBookmark' icon={
        <MaterialIcon hasRipple icon='add' onClick={() => console.log('add link')} />
      } />
    </div>
  }
}

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

function s2p(state) {
  return {
    bookmarks: state.bookmarks ? state.bookmarks : { loaded: false },
    filteredBookmarks: state.bookmarks && state.bookmarks.loaded
      ? folders[state.bookmarks.selectedIndex].filter(state.bookmarks.bookmarks)
      : []
  }
}

export default connect(s2p, null)(HomePage)

function BookmarkCard(props) {
  const card = props.card
  const openLink = () => window.open(card.url, 'bookmark-window-' + card.id)
  return <Card className='link-card' key={card.id}>
    <CardPrimaryContent className='primary-content' onClick={openLink}>
      <Headline6 className='headline-6'>
        {card.inbox ? <div className='inbox'> <MaterialIcon hasRipple icon='inbox' /> &nbsp;</div> : <span />}
        {card.title.trim() !== '' ? card.title.trim() : card.url}
      </Headline6>
      <Caption>
        {card.host} - {moment(card.created_at).fromNow()}
        {card.last_status_code !== 200
          ? [
            ' - ',
            card.last_status_code === 0
              ? 'unknown HTTP status'
              : 'HTTP ' + card.last_status_code
          ]
          : ''}
      </Caption>
    </CardPrimaryContent>
    <CardActions>
      <CardActionButtons>
        <Button onClick={openLink}>Open</Button>
      </CardActionButtons>
      <CardActionIcons>
        {card.inbox
          ? <MaterialIcon
            hasRipple icon='visibility'
            onClick={props.markAsRead} />
          : <div />}
        <MaterialIcon
          hasRipple icon='remove'
          onClick={props.delete} />
      </CardActionIcons>
    </CardActions>
  </Card>
}
