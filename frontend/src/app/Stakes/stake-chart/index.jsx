import { ResponsivePieCanvas } from '@nivo/pie'
import React, { useState, useEffect } from "react";
import axios from 'axios';
import { get } from "lodash";
import './index.scss'
import ScaleLoader from "react-spinners/ScaleLoader";
import { useMediaQuery } from 'react-responsive'

export const StakeChart = () => {
  const [data, setData] = useState(null)
  const [totalToken, setTotalToken] = useState(0)
  const chartItemNumber = 10
  const color = ['#D4717A', '#D0E658', '#F2994A', '#F2C23E', '#B6DF72', '#CD4445', '#8A5B56', '#B95585', '#6B6E83', '#907653', '#8652C9']

  useEffect(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    const params = {
      limit: 10,
      offset: 0,
    };

    axios.get(`${process.env.REACT_APP_API}/staking/wallets`, {
      cancelToken: source.token,
      params
    })
      .then((res) => {
        const items = get(res, 'data.data.Delegators', [])
        const total = +get(res, 'data.data.TotalTokens', 0)
        const result = items.reduce((acc, item, index) => {
          if (index < chartItemNumber) {
            acc.push({
              "id": item.delegator_address,
              "label": item.delegator_address,
              "value": Math.round(+item.shares * 100) / 100,
              "color": color[index % color.length]
            })
          } else if (index === chartItemNumber) {
            const restNode = {
              "id": 'rest-node',
              "label": 'Rest Node',
              "value": acc.reduce((acc, data) => acc - data.value, total),
              "color": '#A5ACB9',
            }
            acc.push(restNode)
          }
          return acc
        }, [])
        setData(result)
        setTotalToken(Math.round(total * 100) / 100)
      });
    return () => source.cancel();
  }, []);

  const isTabletOrMobile = useMediaQuery({ query: '(max-width: 1224px)' })
  console.log(isTabletOrMobile)

  if (!data) {
    return <div className='d-flex justify-content-center py-5'>
      <ScaleLoader
        width={3}
        height={27}
        color={"#fff"}
        loading={!data}
      />
    </div>
  }

  if (!data.length) {
    return null
  }

  return <div style={{ height: 300, position: 'relative' }}>
    <ResponsivePieCanvas
      data={data}
      margin={{ top: 40, right: isTabletOrMobile ? 0 : 200, bottom: isTabletOrMobile ? 120 : 40, left: isTabletOrMobile ? 15 : 80 }}
      innerRadius={0.7}
      colors={{ scheme: 'paired' }}
      borderColor={{ from: 'color', modifiers: [['darker', 0.6]] }}
      radialLabelsSkipAngle={10}
      radialLabelsTextColor="#fff"
      enableRadialLabels={false}
      enableSliceLabels={false}
      radialLabelsLinkColor={{ from: 'color' }}
      sliceLabelsSkipAngle={10}
      isInteractive={false}
      colors={d => {
        return d.data.color
      }}
      defs={[]}
      fill={[]}
      legends={[
        {
          anchor: isTabletOrMobile ? 'bottom' : 'right',
          direction: 'column',
          justify: false,
          translateX: isTabletOrMobile ? -120 : 0,
          translateY: isTabletOrMobile ? 120 : 0,
          itemWidth: 100,
          itemHeight: 20,
          itemsSpacing: 2,
          symbolSize: 20,
          itemDirection: 'left-to-right',
          itemTextColor: '#E6E6E6',
        }
      ]}
    />
    <div className='chart-layer d-xl-flex flex-column text-center font-weight-light d-none'>
      <span className='h5 font-weight-light'>{new Intl.NumberFormat().format(+totalToken)}</span>
      <span>tokens</span>
    </div>
  </div>
}