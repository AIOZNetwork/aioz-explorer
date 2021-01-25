import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import axios from 'axios';
import { get } from "lodash";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import ScaleLoader from "react-spinners/ScaleLoader";
import Button from 'react-bootstrap/Button'
import HashLoader from "react-spinners/HashLoader";
import { sharesFormatter } from './../../_helpers/columnFormatter'

export default function () {
  const limit = 10;

  const columns = [
    {
      dataField: "address",
      text: "Address",
      style: { maxWidth: '280px' },
      formatter: (cellContent, row) => cellContent === 'TOTAL' ? <div className='text-center'>{cellContent}</div> : <Link className='text-truncate d-block' to={`/address/${cellContent}`}>{cellContent}</Link>
    },
    {
      dataField: "tokens",
      text: "Token Staked",
      headerClasses: 'text-right',
      classes: 'text-truncate text-success text-right',
      style: { maxWidth: '150px' },
      formatter: sharesFormatter,
    },
    {
      dataField: "%Staked",
      headerClasses: 'text-right',
      classes: 'text-right',
      text: "%Staked",
      formatter: (cellContent, row) => `${Math.round((+row.tokens || 0) / (totalStake || 1) * 10000) / 100} %`
    },
  ]

  const [totalStake, setTotalStake] = useState(0)
  const [items, setItems] = useState(null)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [isLoading, setIsLoading] = useState(false)

  function getParams() {
    const offset = page * limit - limit;
    const params = {
      limit,
      offset,
    }

    return params;
  }

  useEffect(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    const params = getParams();
    setIsLoading(true)

    axios.get(`${process.env.REACT_APP_API}/staking/validators`, {
      cancelToken: source.token,
      params
    })
      .then((res) => {
        const total = get(res, 'data.total', 0)
        const list = get(res, `data.data.Validators`, []);
        const totalTokens = get(res, `data.data.TotalTokens`, 0);
        setIsLoading(false)
        setTotal(total)
        setTotalStake(totalTokens)
        setItems(items ? [...items, ...list] : list);
      })
      .catch(() => setIsLoading(false))
    return () => source.cancel();
  }, [page]);

  return <>
    <div className="border mt-5">
      <Row>
        <Col xs={12}>
          <div className='bg-light  p-3'>
            <span className='h5 text-uppercase'>TOP VALIDATOR / GUARDIAN NODES</span>
          </div>
        </Col>
      </Row>
      {
        items ? <div className='px-3 bg-secondary'>
          <BootstrapTable
            striped
            bootstrap4
            remote
            keyField="address"
            wrapperClasses="table-responsive"
            classes="table-vertical-center overflow-hidden mb-0"
            data={[...items, {
              address: 'TOTAL',
              tokens: totalStake,
            }]}
            columns={columns}
          />
          {
            total > (page) * limit ? <div className='d-flex justify-content-center mt-2 mb-3'>
              <div className='my-1 mr-2'><HashLoader
                size={15}
                color={"#fff"}
                loading={isLoading}
              /></div>
              <Button disabled={isLoading} size="sm" variant="link" onClick={() => setPage(page + 1)} className='p-0 text-uppercase'>View More</Button>
            </div> : null
          }

        </div> : <div className='d-flex justify-content-center py-5'>
            <ScaleLoader
              width={3}
              height={27}
              color={"#fff"}
              loading={!items}
            />
          </div>
      }

    </div>
  </>;
}
