import React, { useState, useEffect } from "react";
import Delegators from './delegators'
import Validators from './validators'
import { StakeChart } from './stake-chart'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { useAnalytics } from 'reactfire';

export default function ({location}) {
  const analytics = useAnalytics();
  useEffect(() => {
    analytics.logEvent('stakes', { path_name: location.pathname });
  }, [location.pathname]);
  return <>
    <Row>
      <Col xl={{ span: 10, offset: 1 }}>
        <StakeChart />
      </Col>
    </Row>
    <Row>
      <Col md={6}>
        <Validators />
      </Col>
      <Col md={6}>
        <Delegators />
      </Col>
    </Row>
  </>
}
