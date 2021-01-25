import React, {useEffect} from "react";
import RecentTxs from './recent-txs'
import RecentBlocks from './recent-blocks'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Stats from './stats'
import { useAnalytics } from 'reactfire';

export default function ({location}) {
  const analytics = useAnalytics();
  useEffect(() => {
    analytics.logEvent('home', { path_name: location.pathname });
  }, [location.pathname]);
  return <>
    <Stats />
    <Row>
      <Col md={6}>
        <RecentBlocks />
      </Col>
      <Col md={6}>
        <RecentTxs />
      </Col>
    </Row>
  </>
}
