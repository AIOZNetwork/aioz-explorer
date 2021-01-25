import React, { useState } from "react";
import { Link } from 'react-router-dom';
import "./header.scss";
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import Form from 'react-bootstrap/Form'
import FormControl from 'react-bootstrap/FormControl'
import Button from 'react-bootstrap/Button'
import { withRouter } from 'react-router-dom';
import { ReactComponent as IcSearch } from './../../../assets/svg/search.svg';
import { ReactComponent as Logo } from './../../../assets/svg/logo.svg';
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Container from 'react-bootstrap/Container'

export default withRouter((props) => {
  const { location } = props;

  const [searchValue, setSearchValue] = useState('');
  const [open, setOpen] = useState();

  const onTxtChange = (evt) => {
    setSearchValue(evt.target.value);
  };

  const getRouterPathBySearchInput = (searchInput) => {
    const inputValue = searchInput.trim();

    if (/^[0-9]*$/.test(inputValue)) {
      // is Block
      return `/blocks/${inputValue}`;
    }

    if (/^(aioz|aiozvaloper)(-)[a-zA-Z0-9]{39}$/.test(inputValue)) {
      // is Wallet
      return `/address/${inputValue}`;
    }

    if (/^(0x)?[0-9a-fA-F]{64}$/.test(inputValue)) {
      // is Transaction
      return `/transactions/${inputValue}`;
    }
    return `/search-not-found/`;
  }

  const onSubmitHandler = (evt) => {
    evt.preventDefault();
    if (!searchValue) {
      return
    }
    props.history.push(getRouterPathBySearchInput(searchValue));
    setSearchValue("");
  };
  return <>
    <header>
      <Container>
        <Navbar variant="dark" expand="lg" onToggle={() => setOpen(true)}
          expanded={open}>
          <Navbar.Brand to="/" as={Link} className='text-decoration-none text-white'>
            <span className='pb-1 d-inline-block position-relative'>
              <Logo width='100' />
              <span className='identical text-uppercase d-none'>mainnet</span>
            </span>
            <span className='d-md-inline-block align-middle ml-4 text-uppercase d-none'>Blockchain Explorer</span>
          </Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse >
            <Nav className="ml-auto custom-nav" activeKey={location.pathname} onSelect={() => setOpen(false)}>
              <Nav.Item>
                <Nav.Link as={Link} to="/blocks" eventKey="/blocks">Blocks</Nav.Link>
              </Nav.Item>
              <Nav.Item>
                <Nav.Link as={Link} to="/transactions" eventKey="/transactions">Transactions</Nav.Link>
              </Nav.Item>
              <Nav.Item>
                <Nav.Link as={Link} to="/stakes" eventKey="/stakes">Stakes</Nav.Link>
              </Nav.Item>
              <Nav.Item>
                <Nav.Link as={Link} to="/stats" eventKey="/stats">Nodes Stats</Nav.Link>
              </Nav.Item>
            </Nav>

          </Navbar.Collapse>
        </Navbar>
      </Container>
    </header>
    {
      location && location.pathname === '/stats'
        ? null
        : <Container>
          <Row>
            <Col md={{ span: 8, offset: 2 }}>
              <Form inline onSubmit={onSubmitHandler} className='my-3 my-sm-5 position-relative search-block'>
                <Form.Group controlId="searchText" className='w-100 mb-0'>
                  <FormControl value={searchValue} type="text" placeholder="Search for wallet address, block or transaction" className=" py-4 mr-2 w-100 bg-white border-0 pr-5 text-truncate" onChange={onTxtChange} size='lg' />
                </Form.Group>

                <Button type='submit' size='lg' variant="outline-dark" className=' position-absolute mr-2 border-0 bg-white h-100' style={{ right: 0 }}>
                  <IcSearch style={{ fill: '#444444' }} />
                </Button>
              </Form>
            </Col>
          </Row>
        </Container>
    }

  </>
})
