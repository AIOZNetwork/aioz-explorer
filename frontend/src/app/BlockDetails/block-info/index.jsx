import React from 'react'
import { Link } from 'react-router-dom';
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Skeleton, { SkeletonTheme } from "react-loading-skeleton";
import { addressFormatter, validFormatter } from './../../_helpers/columnFormatter'
export default function BlockInfo({ data }) {
  const {
    height,
    status,
    timestamp,
    hash,
    amount,
    prevBlock,
    proposer,
    stateHash,
    txnsHash
  } = data;

  return (
    <SkeletonTheme color="#141414" highlightColor="#222">
      <div className="px-3 py-md-2 bg-dark">
        <Row className='py-md-2'>
          <Col md={6}>
            <div className='py-1 text-white-50'>Status</div><div className='py-1 border-top text-white'>{status || <Skeleton width={60} />}</div>
          </Col>
          <Col md={6}>
            <div className='py-1 text-white-50'>Time</div><div className='py-1 border-top text-white'>{timestamp || <Skeleton width={130} />}</div>
          </Col>
        </Row>
        <Row className='py-md-2'>
          <Col md={6}>
            <div className='py-1 text-white-50'>Hash</div><div className='py-1 text-lowercase border-top text-white text-truncate'>{hash || <Skeleton width={550} />}</div>
          </Col>
          <Col md={6}>
            <div className='py-1 text-white-50'>Transactions</div><div className='py-1 border-top text-white'>{amount === undefined ? <Skeleton width={30} /> : amount}</div>
          </Col>
        </Row>
        <Row className='py-md-2'>
          <Col md={6}>
            <div className='py-1 text-white-50'>Proposer</div><div className='py-1 text-lowercase border-top text-white text-truncate'>{proposer || <Skeleton width={350} />}</div>
          </Col>
          <Col md={6}>
            <div className='py-1 text-white-50'>State Hash</div><div className='py-1 text-lowercase border-top text-white text-truncate'>{stateHash === undefined ? <Skeleton width={550} /> : stateHash}</div>
          </Col>
        </Row>
        <Row className='py-md-2'>
          <Col md={6}>
            <div className='py-1 text-white-50'>Previous Block</div><div className='py-1 text-lowercase border-top text-white text-truncate'>{height > 1 ? <Link to={`/blocks/${height - 1}`}>{prevBlock === undefined ? <Skeleton width={540} /> : prevBlock}</Link> : 'None'}</div>
          </Col>
          <Col md={6}>
            <div className='py-1 text-white-50'>Txns Hash</div><div className='py-1 text-lowercase border-top text-white text-truncate'>{txnsHash === undefined ? <Skeleton width={540} /> : txnsHash}</div>
          </Col>
        </Row>
      </div>
    </SkeletonTheme>
  )
}
