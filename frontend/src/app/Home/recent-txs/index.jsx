import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import axios from 'axios';
import { get } from "lodash";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import { ReactComponent as IcExchange } from './../../../assets/svg/ic-exchange.svg';
import useTimer from './../../_helpers/useTimer'
import { txnFormatter, txnsTypeFormatter } from './../../_helpers/columnFormatter'
import ScaleLoader from "react-spinners/ScaleLoader";

export default function () {
  const limit = 5;
  const columns = [
    {
      dataField: "Payload",
      text: "Type",
      style: { maxWidth: '200px', minWidth: '70px' },
      formatter: txnsTypeFormatter
    },
    {
      dataField: "Hash",
      text: "Tx Hash",
      style: { maxWidth: '280px' },
      formatter: txnFormatter
    }
  ]

  const [items, setItems] = useState(null)

  function getParams() {
    const offset = 0;
    const params = {
      limit,
      offset,
    }

    return params;
  }

  useTimer(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    const params = getParams();

    axios.get(`${process.env.REACT_APP_API}/transactions`, {
      cancelToken: source.token,
      params
    })
      .then((res) => {
        const list = get(res, `data.data`, []);
        setItems(list);
      });
    return () => source.cancel();
  }, 0, 5 * 1000);

  return <>
    <div className="border my-4">
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
          <BootstrapTable
            striped
            bootstrap4
            remote
            keyField="Hash"
            wrapperClasses="table-responsive px-3 bg-secondary"
            classes="table-vertical-center overflow-hidden mb-0"
            data={items}
            columns={columns}
          />
          <div className='d-flex justify-content-center mt-2 mb-3 text-uppercase'>
            <Link to={`/transactions`}>View More</Link>
          </div>
        </> : <div className='d-flex justify-content-center py-5'>
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
