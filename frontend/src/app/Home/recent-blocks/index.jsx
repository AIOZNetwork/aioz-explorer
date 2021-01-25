import React, { useState, useEffect } from "react";
import useTimer from './../../_helpers/useTimer'
import { blockFormatter } from './../../_helpers/columnFormatter'
import { Link } from "react-router-dom";
import axios from 'axios';
import { get } from "lodash";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { ReactComponent as IcBlocks } from './../../../assets/svg/ic-blocks.svg';
import BootstrapTable from 'react-bootstrap-table-next';
import ScaleLoader from "react-spinners/ScaleLoader";

export default function () {
  const limit = 5;

  const columns = [
    {
      dataField: "Height",
      text: "Height",
    },
    {
      dataField: "AppHash",
      text: "Block Hash",
      classes: 'text-truncate',
      style: { maxWidth: '300px' },
      formatter: (cellContent, row) => <Link className='text-truncate text-lowercase d-block' to={`/blocks/${row.Height}`}>{cellContent}</Link>
    },
    {
      dataField: "NumTxs",
      classes: 'text-center',
      text: "TXNS",
    },
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

    axios.get(`${process.env.REACT_APP_API}/blocks`, {
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
            <span className='ico-block mr-2'></span>
            <span className='h5 text-uppercase'>Blocks</span>
          </div>
        </Col>
      </Row>
      {
        items ? <>
          <BootstrapTable
            striped
            bootstrap4
            remote
            keyField="Height"
            wrapperClasses="table-responsive px-3 bg-secondary"
            classes="table-vertical-center overflow-hidden mb-0"
            data={items}
            columns={columns}
          />
          <div className='d-flex justify-content-center mt-2 mb-3 text-uppercase'>
            <Link to={`/blocks`}>View More</Link>
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
