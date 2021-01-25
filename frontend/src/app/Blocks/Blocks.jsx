import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import axios from 'axios';
import { get } from "lodash";
import { formatDate } from "./../_helpers";
import { coinsFormatter } from './../_helpers/columnFormatter'
import { ReactComponent as IcBlocks } from './../../assets/svg/ic-blocks.svg';
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory, { PaginationProvider, PaginationListStandalone } from 'react-bootstrap-table2-paginator';
import ScaleLoader from "react-spinners/ScaleLoader";
import overlayFactory from 'react-bootstrap-table2-overlay';
import { useAnalytics } from 'reactfire';

export default function ({ location }) {
  const analytics = useAnalytics();
  const limit = 15;
  const columns = [
    {
      dataField: "Height",
      text: "Height",
    },
    {
      dataField: "AppHash",
      text: "Block Hash",
      style: { maxWidth: '500px' },
      formatter: (cellContent, row) => <Link className='text-truncate text-lowercase d-block' to={`/blocks/${row.Height}`}>{cellContent}</Link>
    },
    {
      dataField: "NumTxs",
      classes: 'text-center',
      text: "TXNS",
    },
    {
      dataField: "AvgFee",
      text: "Avg Fee",
      classes: 'text-truncate text-center text-warning',
      formatter: (coin) => coinsFormatter(coin),
    },
    {
      dataField: "Time",
      text: "Age",
      classes: 'text-truncate text-info',
      formatter: formatDate,
    },
  ]

  const [items, setItems] = useState(null)
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  function getParams() {
    const offset = currentPage * limit - limit;
    const params = {
      limit,
      offset,
    }

    return params;
  }

  useEffect(() => {
    analytics.logEvent('blocks', { path_name: location.pathname });
  }, [location.pathname]);

  useEffect(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    const params = getParams();
    setLoading(true)

    axios.get(`${process.env.REACT_APP_API}/blocks`, {
      cancelToken: source.token,
      params
    })
      .then((res) => {
        const total = get(res, `data.total`, 0);
        const list = get(res, `data.data`, []);
        setLoading(false)
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
          <div className='bg-light  p-3'>
            <span className='ico-block mr-2'></span><span className='h5 text-uppercase'>Blocks</span>
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
                    striped
                    bootstrap4
                    loading={loading}
                    remote
                    keyField="Height"
                    wrapperClasses="table-responsive"
                    classes="table-vertical-center overflow-hidden mb-0"
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
