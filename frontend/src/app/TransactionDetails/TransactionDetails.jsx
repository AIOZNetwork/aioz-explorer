import React, { useState, useEffect } from 'react'
import { get } from 'lodash';
import axios from 'axios';
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import TransactionInfo from './transaction-info';
import ScaleLoader from "react-spinners/ScaleLoader";
import { formatDateTs, lsCointToAmount } from '../_helpers';
import { addressFormatter, blockFormatter, coinsFormatter, txnTypeFormatter, validFormatter } from '../_helpers/columnFormatter';
import { Link, } from "react-router-dom";
import BootstrapTable from 'react-bootstrap-table-next';
import SearchNotFound from './../SearchNotFound/SearchNotFound'
import { useAnalytics } from 'reactfire';

export default function (props) {
  const { match, history, location } = props;

  const hashId = get(match, 'params.hashId');
  const [transaction, setTransaction] = useState();
  const [msg, setMsg] = useState();
  const [isShowNotFound, setIsShowNotFound] = useState(false)
  const columns = [
    {
      dataField: "message_type",
      text: "Type",
      style: { maxWidth: '200px', minWidth: '110px' },
      formatter: (cell) => txnTypeFormatter(cell)
    },
    {
      dataField: "address_from[0]",
      style: { maxWidth: '240px' },
      text: "From Address",
      formatter: addressFormatter
    },
    {
      dataField: "address_to[0]",
      style: { maxWidth: '240px' },
      text: "To Address",
      formatter: addressFormatter
    },
    {
      dataField: "block_height",
      text: "Block",
      formatter: blockFormatter,
    },
    {
      dataField: "amount",
      text: "Amount",
      classes: 'text-success',
      formatter: (coin) => coinsFormatter(coin),
      // style: { maxWidth: '120px' },
    },
    {
      dataField: "is_valid",
      text: "Result",
      formatter: validFormatter
    },
    {
      dataField: "block_time",
      style: { minWidth: '120px' },
      classes: 'text-info',
      formatter: (cellContent) => formatDateTs(+cellContent * 1000) || '',
      text: "Age",
    },
  ]

  const analytics = useAnalytics();
  useEffect(() => {
    analytics.logEvent('txn-details', { path_name: location.pathname });
  }, [location.pathname]);

  useEffect(() => {
    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();
    setMsg()
    setTransaction()
    axios.get(`${process.env.REACT_APP_API}/transaction/${(hashId + '').toLowerCase()}`).then(res => {
      const data = get(res, 'data.data');
      const Fee = get(data, 'Fee', '{}');
      const fee = JSON.parse(Fee)
      const msg = get(data, 'Payload')
      setTransaction({
        type: msg,
        status: data.IsValid ? 'Success' : 'Fail',
        block: data.BlockHeight,
        time: formatDateTs(+data.BlockTime * 1000) || '',
        gas: `${Intl.NumberFormat().format(data.Gas)} / ${Intl.NumberFormat().format(fee.gas)}`,
        fee: lsCointToAmount(fee.amount)
      })

      setMsg(msg)
      
      if (isShowNotFound) {
        setIsShowNotFound(false)
      }
    }, ({ response }) => {
      if (response && response.status === 404) {
        setIsShowNotFound(true)
      }
    });

    return () => source.cancel();
  }, [hashId]);

  if (isShowNotFound) {
    return <SearchNotFound />
  }

  return (
    <>
      <div className="border">
        <Row>
          <Col xs={12}>
            <div className='bg-light p-3'>
              <span className='ico-transactions mr-2'></span>
              <span className='text-uppercase h5'>Transaction</span>
            </div>
          </Col>
        </Row>
        <TransactionInfo data={{ ...transaction, hash: hashId }} />
      </div>
      {
        msg ? msg.length ? <div className="border mt-4">
          <Row>
            <Col xs={12}>
              <div className='bg-light  p-3'>
                <span className='h5 text-uppercase'>{msg.length} messages in this transactions</span>
              </div>
            </Col>
          </Row>
          <div className='px-3 bg-secondary'>
            <BootstrapTable
              striped
              bootstrap4
              remote
              keyField='block_time'
              wrapperClasses="table-responsive"
              classes="table-vertical-center overflow-hidden"
              data={msg}
              columns={columns}
            />
          </div>
        </div> : null : <div className='d-flex justify-content-center py-5 my-5'>
            <ScaleLoader
              width={3}
              height={27}
              color={"#fff"}
              loading={!msg}
            />
          </div>
      }

    </>

  )
}