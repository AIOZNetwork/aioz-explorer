import React, { useState, useEffect } from "react";
import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory, { PaginationProvider, PaginationListStandalone } from 'react-bootstrap-table2-paginator';
import axios from 'axios';
import { get } from "lodash";
import AddressDetail from "./address-detail";
import { formatDateTs } from "../_helpers";
import { addressFormatter, blockFormatter, txnFormatter, txnTypeFormatter, coinsFormatter, validFormatter } from "../_helpers/columnFormatter";
import ScaleLoader from "react-spinners/ScaleLoader";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import overlayFactory from 'react-bootstrap-table2-overlay';
import SearchNotFound from './../SearchNotFound/SearchNotFound'
import { useAnalytics } from 'reactfire';

const limit = 10;

export default function ({ match, history, location }) {
  const analytics = useAnalytics();
  const { address } = match.params;

  const [addressDetail, setAddressDetail] = useState();
  const [transactions, setTransactions] = useState();
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [isShowNotFound, setIsShowNotFound] = useState(false)
  const isValoper = String(address).startsWith('aiozvaloper-')
  // ]
  const columns = [
    {
      dataField: "message_type",
      text: "Type",
      formatter: (cell) => txnTypeFormatter(cell),
      style: { maxWidth: '200px', minWidth: '70px' },
    },
    {
      dataField: "address_from[0]",
      style: { maxWidth: '150px' },
      text: "From Address",
      formatter: addressFormatter
    },
    {
      dataField: "address_to[0]",
      style: { maxWidth: '150px' },
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
    },
    {
      dataField: "is_valid",
      text: "Result",
      formatter: validFormatter
    },
    {
      dataField: "transaction_hash",
      text: "Tx Hash",
      style: { maxWidth: '160px' },
      formatter: txnFormatter,
    },
    {
      dataField: "block_time",
      style: { minWidth: '120px' },
      classes: 'text-info',
      formatter: (cellContent) => formatDateTs(+cellContent * 1000) || '',
      text: "Age",
    },
  ]

  useEffect(() => {
    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();
    setAddressDetail()
    setTransactions()
    axios
      .get(`${process.env.REACT_APP_API}/wallet/${(address + '').toLowerCase()}`, {
        cancelToken: source.token
      })
      .then(res => {
        const { staked_info, StakedInfo, ...rest } = get(res, 'data.data');
        if (isValoper) {
          rest.stakeInfo = staked_info || [];
        } else {
          rest.stakeInfo = StakedInfo || [];
        }
        setAddressDetail(rest)
        if (isShowNotFound) {
          setIsShowNotFound(false)
        }
      })
      .catch((err) => {
        const errStatus = get(err, 'response.status');
        if (errStatus === 404 || errStatus === 400) {
          setIsShowNotFound(true)
        }
      });

    return () => source.cancel();
  }, [address]);

  useEffect(() => {
    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();
    const offset = (page - 1) * limit
    setLoading(true)

    axios
      .get(`${process.env.REACT_APP_API}/msgs/${address}`, {
        params: { limit, offset },
        cancelToken: source.token
      })
      .then(res => {
        setLoading(false)
        const total = get(res, `data.total`, 0);
        const list = get(res, `data.data`, []);
        setTotal(total);
        setTransactions(list);
      })
      .catch(() => setLoading(false));

    return () => source.cancel();
  }, [address, page]);

  useEffect(() => {
    analytics.logEvent('address-details', { path_name: location.pathname });
  }, [location.pathname]);

  if(isShowNotFound) {
    return <SearchNotFound />
  }

  return <>
    <AddressDetail {...addressDetail} address={address} isValoper={isValoper} />

    {transactions ? transactions.length ? <div className="border mt-4">
      <Row>
        <Col xs={12}>
          <div className='info-details__heading bg-light  p-3'>
            <span className='h5 text-uppercase'>{total} messages from this address</span>
          </div>
        </Col>
      </Row>
      <PaginationProvider pagination={paginationFactory({
        custom: true,
        totalSize: total,
        sizePerPage: limit,
        withFirstAndLast : false
      })}>
        {({
          paginationProps,
          paginationTableProps
        }) => <>
            <div className='px-3 bg-secondary'>
              <BootstrapTable
                loading={loading}
                striped
                bootstrap4
                remote
                keyField="transaction_hash"
                wrapperClasses="table-responsive"
                classes="table-vertical-center overflow-hidden mb-0"
                data={transactions}
                columns={columns}
                overlay={overlayFactory({ spinner: true, styles: { overlay: (base) => ({ ...base, background: 'rgba(0,0,0,0.75)' }) } })}
                onTableChange={(_, { page }) => { setPage(page) }}
                {...paginationTableProps}
              />
              {
                total > limit ? <div className='d-flex justify-content-center mt-3'>
                  <PaginationListStandalone
                    {...paginationProps}
                  />
                </div> : null
              }
            </div>
          </>
        }
      </PaginationProvider>
    </div> : null : <div className='d-flex justify-content-center py-5 my-5'>
        <ScaleLoader
          ScaleLoader
          width={3}
          height={27}
          color={"#fff"}
          loading={!transactions}
        />
      </div>}
  </>
}
