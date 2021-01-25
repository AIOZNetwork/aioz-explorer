import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom';
import axios from 'axios';
import { get } from 'lodash';
import { ReactComponent as IcBlocks } from './../../assets/svg/ic-blocks.svg';
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory, { PaginationProvider, PaginationListStandalone } from 'react-bootstrap-table2-paginator';
import BlockInfo from './block-info';
import { formatDateTs, lsCointToAmount } from "./../_helpers";
import { txnFormatter, validFormatter, feeFormatter, cointFromMsgFormatter, txnsTypeFormatter } from "./../_helpers/columnFormatter";
import ScaleLoader from "react-spinners/ScaleLoader";
import overlayFactory from 'react-bootstrap-table2-overlay';
import SearchNotFound from './../SearchNotFound/SearchNotFound'
import { useAnalytics } from 'reactfire';

export default function (props) {
  const analytics = useAnalytics();
  const limit = 10;
  const { match, history, location } = props;
  const height = get(match, 'params.height');

  const [detail, setDetail] = useState();
  const [txs, setTxs] = useState();
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [isShowNotFound, setIsShowNotFound] = useState(false)

  const columns = [
    {
      dataField: "type",
      text: "Type",
      style: { maxWidth: '200px', minWidth: '70px' },
      formatter: (_, row) => txnsTypeFormatter(row.Payload)
    },
    {
      dataField: "Hash",
      text: "Tx Hash",
      style: { maxWidth: '400px' },
      formatter: txnFormatter
    },
    {
      dataField: "IsValid",
      text: "Result",
      style: { minWidth: '50px' },
      formatter: validFormatter
    },
    {
      dataField: "Payload",
      text: "Amount",
      formatter: cointFromMsgFormatter,
    },
    {
      dataField: "Fee",
      style: { minWidth: '50px' },
      text: "Fee",
      classes: 'text-warning',
      formatter: feeFormatter,
    },
    {
      dataField: "BlockTime",
      classes: 'text-info',
      style: { minWidth: '120px' },
      formatter: (cellContent) => formatDateTs(+cellContent * 1000),
      text: "Time",
    },
  ]

  // get block detail
  useEffect(() => {
    const CancelToken = axios.CancelToken;
    const source = CancelToken.source();
    setTxs()
    setDetail()
    axios
      .get(`${process.env.REACT_APP_API}/block/height/${height}`, {
        cancelToken: source.token
      })
      .then((res) => {
        const data = get(res, 'data.data');
        // DataHash -> TXNS HASH if blocks has tx
        // Block hash -> HeaderHash
        // prevBlock hash -> LastBlockId
        setDetail({
          status: 'Finalized',
          timestamp: formatDateTs(+data.Time) || '',
          hash: data.HeaderHash || '',
          amount: data.NumTxs || 0,
          prevBlock: data.LastBlockId,
          proposer: data.ProposerAddress,
          txnsHash: data.DataHash,
          stateHash: data.AppHash
        });
        if (isShowNotFound) {
          setIsShowNotFound(false)
        }
      }, ({ response }) => {

        if (response && response.status === 404) {
          setIsShowNotFound(true)
        }
      });

    return () => source.cancel();

  }, [height]);

  // get transaction list
  useEffect(() => {
    const CancelToken = axios.CancelToken;
    const source = CancelToken.source();
    const offset = page * limit - limit;
    setLoading(true)
    axios
      .get(`${process.env.REACT_APP_API}/transactions/height/${height}`, {
        params: {
          offset,
          limit: limit,
        },
        cancelToken: source.token
      })
      .then(res => {
        setLoading(false)
        const total = get(res, `data.total`, 0);
        const list = get(res, `data.data`, []);
        setTotal(total);
        setTxs(list);
      })
      .catch(() => setLoading(false));

    return () => source.cancel();

  }, [page, height]);

  useEffect(() => {
    analytics.logEvent('block-details', { path_name: location.pathname });
  }, [location.pathname]);

  if (isShowNotFound) {
    return <SearchNotFound />
  }

  return (
    <>
      <div className="border">
        <Row>
          <Col xs={12}>
            <div className='bg-light  p-3'>
              <span className='ico-block mr-2'></span><span className='h5 text-uppercase'>Block&nbsp;{detail ?
                <span className="hash  font-weight-normal">#{height}</span> :
                null}</span>
            </div>
          </Col>
        </Row>
        <BlockInfo data={{ ...detail, height }} />
      </div>

      {txs ? txs.length ? <>
        <div className="border mt-4">
          <Row>
            <Col xs={12}>
              <div className='bg-light  p-3'>
                <span className='h5 text-uppercase'>{total} transactions in this block</span>
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
            }) => <div className='px-3 bg-secondary'>
                <BootstrapTable
                  loading={loading}
                  striped
                  bootstrap4
                  remote
                  keyField="Hash"
                  wrapperClasses="table-responsive"
                  classes="table-vertical-center overflow-hidden mb-0"
                  data={txs}
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
            }
          </PaginationProvider>
        </div>
      </> : null : <div className='d-flex justify-content-center py-5 my-5'>
          <ScaleLoader
            ScaleLoader
            width={3}
            height={27}
            color={"#fff"}
            loading={!txs}
          />
        </div>
      }
    </>
  )
}
