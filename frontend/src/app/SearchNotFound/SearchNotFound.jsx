import React from "react";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
export default function () {

  return <>
    <Row className='align-items-center' style={{ minHeight: '700px' }}>
      <Col xs={12} >
        <h1 style={{ textAlign: 'center' }}>No Result !!!</h1>
      </Col>
    </Row>
  </>;
}
