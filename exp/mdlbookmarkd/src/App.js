import React from 'react';
import MaterialIcon from '@material/react-material-icon';
import TopAppBar, { TopAppBarFixedAdjust, TopAppBarIcon, TopAppBarRow, TopAppBarSection, TopAppBarTitle } from '@material/react-top-app-bar';
import Card, {
  CardPrimaryContent,
  CardMedia,
  CardActions,
  CardActionButtons,
  CardActionIcons
} from "@material/react-card";
import Button from '@material/react-button';
import { Headline6, Caption } from '@material/react-typography';
import IconButton from '@material/react-icon-button';
import { Cell, Grid, Row } from '@material/react-layout-grid';
import Tab from '@material/react-tab';
import TabBar from '@material/react-tab-bar';

import './App.scss';
import '@material/react-button/index.scss';
import '@material/react-card/index.scss';
import '@material/react-icon-button/index.scss';
import '@material/react-layout-grid/index.scss';
import '@material/react-material-icon/index.scss';
import '@material/react-tab-bar/index.scss';
import '@material/react-tab-indicator/index.scss';
import '@material/react-tab-scroller/index.scss';
import '@material/react-tab/index.scss';
import '@material/react-top-app-bar/index.scss';
import '@material/react-typography/index.scss';

class App extends React.Component {
  state = { activeIndex: 0 };
  handleActiveIndexUpdate = (activeIndex) => this.setState({ activeIndex });

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

    return (
      <div>
        <TopAppBar fixed>
          <TopAppBarRow>
            <TopAppBarSection align='start'>
              <TopAppBarIcon navIcon tabIndex={0}>
                <MaterialIcon hasRipple icon='home' onClick={() => console.log('click')} />
              </TopAppBarIcon>
              <TopAppBarTitle>Bookmarks Manager</TopAppBarTitle>
            </TopAppBarSection>
            <TopAppBarSection align='end' role='toolbar'>
              <TopAppBarIcon actionItem tabIndex={0}>
                <MaterialIcon
                  aria-label="add new bookmark"
                  hasRipple
                  icon='add'
                  onClick={() => console.log('print')}
                />
              </TopAppBarIcon>
            </TopAppBarSection>
          </TopAppBarRow>
        </TopAppBar>
        <TopAppBarFixedAdjust>
          <TabBar
            activeIndex={this.state.activeIndex}
            handleActiveIndexUpdate={this.handleActiveIndexUpdate}
          >
            <Tab>
              <span className='mdc-tab__text-label'>All</span>
            </Tab>
            <Tab>
              <span className='mdc-tab__text-label'>Repeated</span>
            </Tab>
          </TabBar>

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
          </Grid>
        </TopAppBarFixedAdjust>
      </div >
    );
  }
}

export default App;
