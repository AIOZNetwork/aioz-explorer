import React from "react";
import BootstrapTable from 'react-bootstrap-table-next';
import { lsCointToAmount } from "../../_helpers";
import { sharesFormatter, addressFormatter, coinsFormatter } from "../../_helpers/columnFormatter";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import Skeleton, { SkeletonTheme } from "react-loading-skeleton";
import ScaleLoader from "react-spinners/ScaleLoader";

export default ({
  address,
  coins,
  jailed,
  sequence,
  power,
  isValoper,
  status,
  tokens,
  stakeInfo
}) => {
  const sharesColumn = {
    dataField: "shares",
    headerClasses: 'text-right',
    classes: 'text-truncate text-success text-right',
    style: { width: '80px' },
    formatter: sharesFormatter,
    text: "Shares",
  }

  const columns = isValoper ?
    [
      {
        dataField: "delegator_address",
        text: "Delegator Address",
        formatter: addressFormatter
      },
      sharesColumn,
    ] : [
      {
        dataField: "validator_address",
        text: "Validator Address",
        style: { maxWidth: '240px' },
        formatter: addressFormatter
      },
      sharesColumn,
    ]

  return <>
    <div className="border">
      <Row>
        <Col xs={12}>
          <div className='bg-light  p-3'>
            <span className='h5 text-uppercase'>Address</span>
          </div>
        </Col>
      </Row>
      <SkeletonTheme color="#141414" highlightColor="#222">
        <div className="px-3 py-2 bg-dark">
          <Row>
            <Col md={6}>
              <div className="py-md-2">
                <div className='text-white-50 py-1'>Address</div><div className='border-top text-white py-1 text-truncate '>{address}</div>
              </div>
              {
                isValoper ? <>
                  <div className="py-md-2">
                    <div className='text-white-50 py-1'>Tokens</div>
                    <div className='border-top text-white py-1 text-truncate '>{tokens !== undefined ? new Intl.NumberFormat().format(+tokens) : <Skeleton width={100} />}</div>
                  </div>
                </> : <div className="py-md-2">
                    <div className='text-white-50 py-1'>Sequence</div>
                    <div className='border-top text-white py-1 text-truncate'>{sequence === undefined ? <Skeleton width={130} /> : sequence}</div>
                  </div>
              }

            </Col>
            <Col md={6}>
              {isValoper ? <>
                <div className="py-md-2">
                  <div className='text-white-50 py-1'>Status</div>
                  <div className='border-top text-white py-1 text-truncate '>{status || <Skeleton width={60} />}</div>
                </div>
                <div className="py-md-2">
                  <div className='text-white-50 py-1'>Power</div>
                  <div className='border-top text-white py-1 text-truncate '>{power || <Skeleton width={50} />}</div>
                </div>
                <div className="py-md-2">
                  <div className='text-white-50 py-1'>Jailed</div>
                  <div className='border-top text-white py-1 text-truncate '>{jailed !== undefined ? (jailed + '') : <Skeleton width={50} />}</div>
                </div>
              </> : <>
                  <div className="py-md-2">
                    <div className='text-white-50 py-1'>Balance</div>
                    <div className='border-top text-white py-1 text-truncate '>{coinsFormatter(coins) || <Skeleton width={60} />}</div>
                  </div>
                </>}
            </Col>
          </Row>
        </div>
      </SkeletonTheme>
    </div>

    {
      stakeInfo ? stakeInfo.length ? <div className="border mt-4">
        <Row>
          <Col xs={12}>
            <div className='bg-light  p-3'>
              <span className='h5 text-uppercase'>{isValoper ? <>Tokens staked to this node by {stakeInfo.length} delegators</> : <>Tokens staked by this node to {stakeInfo.length} validator</>}
              </span>
            </div>
          </Col>
        </Row>
        <div className='px-3 bg-secondary'>
          <BootstrapTable
            striped
            bootstrap4
            remote
            keyField={isValoper ? "delegator_address" : 'validator_address'}
            wrapperClasses="table-responsive"
            classes="table-vertical-center overflow-hidden"
            data={stakeInfo}
            columns={columns}
          />
        </div>
      </div> : null : <div className='d-flex justify-content-center py-5 my-5'>
          <ScaleLoader
            ScaleLoader
            width={3}
            height={27}
            color={"#fff"}
            loading={!stakeInfo}
          />
        </div>
    }
  </>

};
