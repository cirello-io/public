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
import TextField, { Input } from '@material/react-text-field';
import moment from 'moment'
import { Cell, Grid, Row } from '@material/react-layout-grid'
import { Headline6, Caption } from '@material/react-typography'
import { connect } from 'react-redux'

class HomePage extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      fuzzySearch: props.fuzzySearch,
      delete: null,
      addNewBookmark: false,
      newBookmark: {
        title: '',
        url: ''
      }
    }
    this.filterBy = this.filterBy.bind(this)
    this.deleteDialog = this.deleteDialog.bind(this)
    this.deleteAction = this.deleteAction.bind(this)
    this.addNewBookmarkAction = this.addNewBookmarkAction.bind(this)
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
    this.setState({ fuzzySearch: v }, () => {
      this.props.dispatch({ type: 'TRIGGER_FUZZY_SEARCH', fuzzySearch: v })
    })
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
    this.props.dispatch({ type: 'MARK_BOOKMARK_AS_READ', id })
  }

  addNewBookmarkAction() {
    this.setState({ addNewBookmark: null }, () => {
      this.props.dispatch({ type: 'ADD_BOOKMARK', newBookmark: { ...this.state.newBookmark } })
    })
  }

  render() {
    if (!this.dependenciesLoaded()) {
      return (<Grid></Grid>)
    }

    const listing = this.props.bookmarks.filteredBookmarks

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
            <DialogButton action='keep' isDefault>Keep</DialogButton>
            <DialogButton action='delete'>Delete</DialogButton>
          </DialogFooter>
        </Dialog>
        : null}

      {this.state.addNewBookmark
        ? <Dialog
          open
          onClose={(action) => {
            if (action === 'add') {
              this.addNewBookmarkAction()
            }
            this.setState({ addNewBookmark: false })
          }}>
          <DialogTitle>Add Bookmark</DialogTitle>
          <DialogContent>
            <div className='add-new-bookmark-url'>
              <TextField
                label='URL'
                onTrailingIconSelect={() => this.setState({
                  newBookmark: { ...this.state.newBookmark, url: '' },
                })}
                trailingIcon={<MaterialIcon role="button" icon="delete" />} >
                <Input
                  value={this.state.newBookmark.url}
                  onChange={(e) => this.setState({
                    newBookmark: { ...this.state.newBookmark, url: e.currentTarget.value },
                  })} />
              </TextField>
            </div>

            <div>
              <TextField
                label='Title'
                onTrailingIconSelect={() => this.setState({
                  newBookmark: { ...this.state.newBookmark, title: '' },
                })}
                trailingIcon={<MaterialIcon role="button" icon="delete" />} >
                <Input
                  value={this.state.newBookmark.title}
                  onChange={(e) => this.setState({
                    newBookmark: { ...this.state.newBookmark, title: e.currentTarget.value },
                  })} />
              </TextField>
            </div>
          </DialogContent>
          <DialogFooter>
            <DialogButton action='add' isDefault>add</DialogButton>
          </DialogFooter>
        </Dialog>
        : null}

      <Grid key={'homePageRoot'}>
        <Row key={'searchBarRow'} className='searchbar-row'>
          <Cell columns={3} key={'searchBarLeftPadding'} />
          <Cell columns={6} align="middle">
            <TextField
              fullWidth
              label='search'
              onTrailingIconSelect={() => this.filterBy('')}
              trailingIcon={<MaterialIcon role="button" icon="delete" />} >
              <Input
                value={this.state.fuzzySearch}
                onChange={(e) => this.filterBy(e.currentTarget.value)} />
            </TextField>
          </Cell>
          <Cell columns={3} key={'searchBarRightPadding'} />
        </Row>
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
        <MaterialIcon hasRipple icon='add' />
      } onClick={() => this.setState({
        addNewBookmark: true,
        newBookmark: {
          title: '',
          url: ''
        }
      })} />
    </div>
  }
}

function s2p(state) {
  return {
    bookmarks: state.bookmarks
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
