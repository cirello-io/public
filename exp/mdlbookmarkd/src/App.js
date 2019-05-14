import React from 'react';
import MaterialIcon from '@material/react-material-icon';
import TopAppBar, { TopAppBarFixedAdjust, TopAppBarIcon, TopAppBarRow, TopAppBarSection, TopAppBarTitle } from '@material/react-top-app-bar';
import Card, {
  CardPrimaryContent,
  CardActions,
  CardActionButtons,
  CardActionIcons
} from "@material/react-card";
import Button from '@material/react-button';
import { Headline6, Caption } from '@material/react-typography';
import { Cell, Grid, Row } from '@material/react-layout-grid';
import Fab from '@material/react-fab';
import Drawer, {
  DrawerContent,
  DrawerAppContent,
} from '@material/react-drawer';
import List, { ListItem, ListItemText, ListItemGraphic } from '@material/react-list';

import './App.scss';
import '@material/react-button/index.scss';
import '@material/react-card/index.scss';
import '@material/react-drawer/index.scss';
import '@material/react-fab/index.scss';
import '@material/react-layout-grid/index.scss';
import '@material/react-list/index.scss';
import '@material/react-material-icon/index.scss';
import '@material/react-top-app-bar/index.scss';
import '@material/react-typography/index.scss';

class App extends React.Component {
  state = { open: false };

  render() {
    const card = <Card style={{ margin: 5 }}>
      <CardPrimaryContent style={{ paddingLeft: 17, paddingRight: 10 }}>
        <Headline6 style={{
          marginTop: 15, marginBottom: 5,
          overflow: "hidden", textOverflow: "ellipsis",
          whiteSpace: "nowrap"
        }}>The McDonaldization of UXaaaa</Headline6>
        <Caption>uxdesign.cc - 12 days ago</Caption>
      </CardPrimaryContent>
      <CardActions>
        <CardActionButtons>
          <Button>Open</Button>
        </CardActionButtons>
        <CardActionIcons>
          <MaterialIcon
            hasRipple icon='visibility'
            onClick={() => console.log('click')} />
          <MaterialIcon
            hasRipple icon='remove'
            onClick={() => console.log('click')} />
        </CardActionIcons>
      </CardActions>
    </Card>

    const drawer = <Drawer dismissible open={this.state.open}>
      <DrawerContent>
        <List singleSelection selectedIndex={this.state.selectedIndex}>
          <ListItem>
            <ListItemGraphic graphic={<MaterialIcon icon='all_inbox' />} />
            <ListItemText primaryText='Pending' />
          </ListItem>
          <ListItem>
            <ListItemGraphic graphic={<MaterialIcon icon='bookmarks' />} />
            <ListItemText primaryText='Bookmarks' />
          </ListItem>
          <ListItem>
            <ListItemGraphic graphic={<MaterialIcon icon='compare_arrows' />} />
            <ListItemText primaryText='Duplicated' />
          </ListItem>
        </List>
      </DrawerContent>
    </Drawer>

    return (
      <div className='drawer-container'>
        <TopAppBar fixed>
          <TopAppBarRow>
            <TopAppBarSection align='start'>
              <TopAppBarIcon navIcon tabIndex={0}>
                <MaterialIcon hasRipple icon='menu' onClick={
                  () => this.setState({ open: !this.state.open })
                } />
              </TopAppBarIcon>
              <TopAppBarTitle>Bookmarks Manager</TopAppBarTitle>
            </TopAppBarSection>
            <TopAppBarSection align='end' role='toolbar'>
              <TopAppBarIcon actionItem tabIndex={0}>
                <MaterialIcon
                  aria-label="search"
                  hasRipple
                  icon='search'
                  onClick={() => console.log('print')}
                />
              </TopAppBarIcon>
            </TopAppBarSection>
          </TopAppBarRow>
        </TopAppBar>
        <TopAppBarFixedAdjust className='top-app-bar-fix-adjust'>
          {drawer}
          <DrawerAppContent className='drawer-app-content'>
            <Grid>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
              <Row>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
                <Cell columns={4}> {card} </Cell>
              </Row>
            </Grid>
          </DrawerAppContent>
        </TopAppBarFixedAdjust>
        <Fab className='addNewBookmark' icon={
          <MaterialIcon hasRipple icon='add' onClick={() => console.log('add')} />
        } />
      </div>
    );
  }
}

export default App;
