import React, { useState, useEffect } from "react";
import axios from 'axios';
import { get } from "lodash";
import { formatDate } from "./../_helpers";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory, { PaginationProvider, PaginationListStandalone } from 'react-bootstrap-table2-paginator';
import { txnFormatter, validFormatter, blockFormatter, feeFormatter, txnTypeFormatter, cointFromMsgFormatter, txnsTypeFormatter } from './../_helpers/columnFormatter';
import ScaleLoader from "react-spinners/ScaleLoader";
import overlayFactory from 'react-bootstrap-table2-overlay';
import { useAnalytics } from 'reactfire';

export default function ({location}) {
  const limit = 15;

  const columns = [
    {
      dataField: "Hash",
      text: "Tx Hash",
      style: { maxWidth: '300px' },
      formatter: txnFormatter
    },
    {
      dataField: "type",
      text: "Type",
      style: { maxWidth: '200px', minWidth: '70px' },
      formatter: (_, row) => txnsTypeFormatter(row.Payload)
    },
    {
      dataField: "IsValid",
      text: "Result",
      formatter: validFormatter
    },
    {
      dataField: "Payload",
      text: "Amount",
      formatter: cointFromMsgFormatter,
    },
    {
      dataField: "Fee",
      style: { maxWidth: '80px' },
      classes: 'text-warning',
      text: "Fee",
      formatter: feeFormatter,
    },
    {
      dataField: "BlockHeight",
      text: "Block",
      formatter: blockFormatter,
      // classes: 'text-center',
    },
    {
      dataField: "BlockTime",
      classes: 'text-center text-info',
      style: { minWidth: '120px' },
      formatter: (cellContent) => formatDate(+cellContent * 1000),
      text: "Age",
    },
  ]

  const [items, setItems] = useState(null)
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(false);

  function getParams() {
    const offset = currentPage * limit - limit;
    const params = {
      limit,
      offset,
    }

    return params;
  }

  const analytics = useAnalytics();
  useEffect(() => {
    analytics.logEvent('txns', { path_name: location.pathname });
  }, [location.pathname]);

  useEffect(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    const params = getParams();
    setLoading(true)

    axios.get(`${process.env.REACT_APP_API}/transactions`, {
      cancelToken: source.token,
      params
    })
      .then((res) => {
        setLoading(false)
        const total = get(res, `data.total`, 0);
        const list = get(res, `data.data`, []);
        setTotal(total);
        setItems(list);
      })
      .catch(() => setLoading(false));

    return () => source.cancel();
  }, [currentPage]);

  return <>
    <div className="border">
      <Row>
        <Col xs={12}>
          <div className='bg-light p-3'>
            <span className='ico-transactions mr-2'></span>
            <span className='h5 text-uppercase'>Transactions</span>
          </div>
        </Col>
      </Row>
      {
        items ? <>
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
                    keyField="Hash"
                    wrapperClasses="table-responsive"
                    classes="table-vertical-center overflow-hidden"
                    data={items}
                    columns={columns}
                    overlay={overlayFactory({ spinner: true, styles: { overlay: (base) => ({ ...base, background: 'rgba(0,0,0,0.75)' }) } })}
                    onTableChange={(_, { page }) => { setCurrentPage(page) }}
                    {...paginationTableProps}
                  />
                  <div className='d-flex justify-content-center mt-3'>
                    <PaginationListStandalone
                      {...paginationProps}
                    />
                  </div>
                </div>
              </>
            }
          </PaginationProvider>
        </> : <div className='d-flex justify-content-center py-3 my-3'>
            <ScaleLoader
              width={3}
              height={27}
              color={"#fff"}
              loading={!items}
            />
          </div>
      }
    </div>
  </>
}
