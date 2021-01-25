import React, { useState, useEffect } from "react";
import useTimer from './../../_helpers/useTimer'
import axios from 'axios';
import { get } from "lodash";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import ScaleLoader from "react-spinners/ScaleLoader";
import { addressFormatter, coinsFormatter, sharesFormatter } from '../../_helpers/columnFormatter';
import { ReactComponent as GoldCup } from './../../../assets/svg/gold-cup.svg';
import Skeleton, { SkeletonTheme } from "react-loading-skeleton";

export default function () {
  const [stats, setStats] = useState()

  useTimer(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();

    axios.get(`${process.env.REACT_APP_API}/statistic`, {
      cancelToken: source.token,
    })
      .then((res) => {
        const data = get(res, `data.data`, {});
        setStats(data);
      });
    return () => source.cancel();
  }, 0, 5 * 1000);

  return <>
    <SkeletonTheme color="#141414" highlightColor="#222">
      <Row>
        <Col md={7} xl={8}>
          <div className="border my-4">
            <Row>
              <Col xs={12}>
                <div className='bg-light p-3'>
                  <span className='ico-globe mr-2'></span>
                  <span className='h5 text-uppercase'>Network</span>
                </div>
              </Col>
            </Row>
            <div className='px-3 py-4 bg-secondary'>
              <Row>
                <Col xs={6}>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Inflation</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.Inflation ? stats.Inflation : <Skeleton width={100} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Total Wallet</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.TotalWallets ? stats.TotalWallets : <Skeleton width={100} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Total Block</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.TotalBlock !== null ? sharesFormatter(stats.TotalBlock) : <Skeleton width={100} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Circulating Supply</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.CirculatingSupply ? coinsFormatter(stats.CirculatingSupply, true) : <Skeleton width={180} />}</div>
                  </div>

                  <div className="py-2">
                    <div className='text-white-50 py-1'>Market Cap</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.MarketCap ? stats.MarketCap : <Skeleton width={100} />}</div>
                  </div>
                </Col>
                <Col xs={6}>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Block time</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.BlockTime ? stats.BlockTime : <Skeleton width={100} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Avg Fee</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.AvgFee !== null ? coinsFormatter(stats.AvgFee) : <Skeleton width={50} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Total Transaction</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.TotalTransaction !== null ? sharesFormatter(stats.TotalTransaction) : <Skeleton width={100} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Total Supply</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.TotalSupply ? coinsFormatter(stats.TotalSupply, true) : <Skeleton width={180} />}</div>
                  </div>
                  <div className="py-2">
                    <div className='text-white-50 py-1'>Volume 24h</div><div className='border-top text-white py-1 text-truncate '>{stats && stats.Volume24h !== null ? coinsFormatter(stats.Volume24h, true) : <Skeleton width={200} />}</div>
                  </div>
                </Col>
              </Row>
            </div>
          </div>

        </Col>
        <Col md={5} xl={4}>
          <div className="border my-4">
            <Row>
              <Col xs={12}>
                <div className='bg-light p-3'>
                  <span className='ico-currency mr-2'></span>
                  <span className='h5 text-uppercase'>Top Account</span>
                </div>
              </Col>
            </Row>
            <div className='p-3 bg-secondary'>
              {
                stats
                  ? stats.TopAccount && stats.TopAccount.length
                    ? <div>
                      {
                        stats.TopAccount.map(({ address, coins, sequence, account_number }, index) => <div className={`${index ? 'border-top' : ''}`} key={index}>
                          <Row>
                            <Col xs={1} className='d-flex align-items-center pr-0'><div className='text-center w-100' style={{ fontSize: '16px' }}>{index ? index + 1 : <GoldCup />}</div></Col>
                            <Col xs={11} className='py-1'>
                              <div>{addressFormatter(address)}</div>
                              <div className='text-success'>{coinsFormatter(JSON.stringify(coins))}</div>
                            </Col>
                          </Row>

                        </div>)
                      }
                    </div>
                    : null
                  : <div className='d-flex justify-content-center py-5'>
                    <ScaleLoader
                      width={3}
                      height={27}
                      color={"#fff"}
                      loading={!stats}
                    />
                  </div>
              }
            </div>

          </div>
        </Col>
      </Row>
    </SkeletonTheme>

  </>;
}
