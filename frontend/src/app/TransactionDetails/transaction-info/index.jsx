import React from 'react';
import { Link } from 'react-router-dom';
import Skeleton, { SkeletonTheme } from "react-loading-skeleton";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { txnsTypeFormatter } from '../../_helpers/columnFormatter';

export default function TransactionInfo({ data }) {
  const {
    hash,
    type,
    status,
    block,
    time,
    gas,
    fee
  } = data;

  return <SkeletonTheme color="#141414" highlightColor="#222">
    <div className="px-3 py-2 bg-dark">
      <Row className='py-md-2'>
        <Col md={8}>
          <div className='py-1 text-white-50'>Hash</div><div className='py-1 border-top text-white text-truncate text-truncate'>{hash || <Skeleton width={170} />}</div>
        </Col>
        <Col md={4}>
          <div className='py-1 text-white-50'>Type</div><div className='py-1 border-top text-white text-truncate'>{txnsTypeFormatter(type) || <Skeleton width={60} />}</div>
        </Col>
      </Row>
      <Row className='py-md-2'>
        <Col md={8}>
          <div className='py-1 text-white-50'>Time</div><div className='py-1 border-top text-white text-truncate'>{time || <Skeleton width={170} />}</div>
        </Col>
        <Col md={4}>
          <div className='py-1 text-white-50'>Gas (used / wanted)</div><div className='py-1 border-top text-white text-truncate'>{gas || <Skeleton width={120} />}</div>
        </Col>
      </Row>
      <Row className='py-md-2'>
        <Col md={8}>
          <div className='py-1 text-white-50'>Block</div><div className='py-1 border-top text-white text-truncate'>{block === undefined ? <Skeleton width={50} /> : <Link to={`/blocks/${block}`}>{block}</Link>}</div>
        </Col>
        <Col md={4}>
          <div className='py-1 text-white-50'>Fee</div><div className='py-1 border-top text-white text-truncate'>{fee || <Skeleton width={60} />}</div>
        </Col>
      </Row>
      <Row className='py-md-2'>
        <Col md={8}>
          <div className='py-1 text-white-50'>Status</div><div className='py-1 border-top text-white text-truncate'>{status || <Skeleton width={60} />}</div>
        </Col>
      </Row>
    </div>
  </SkeletonTheme>
}
