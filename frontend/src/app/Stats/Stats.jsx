import React, { useState, useEffect } from "react";
import Worker from './worker'
import Validators from './validators'
import LiveNodeMaps from './LiveNodeMaps'
import axios from 'axios';
import { get, findIndex } from "lodash";
import useTimer from './../_helpers/useTimer'
import { useAnalytics } from 'reactfire';

export default function ({location}) {
  const [items, setItems] = useState()
  const [isLoading, setIsLoading] = useState(false)
  const analytics = useAnalytics();

  useEffect(() => {
    analytics.logEvent('stats', { path_name: location.pathname });
  }, [location.pathname]);

  useTimer(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();
    setIsLoading(true)

    axios.get(`${process.env.REACT_APP_API}/node_info`, {
      cancelToken: source.token,
      params: {
        limit: 1000,
        offset: 0,
      }
    })
      .then((res) => {
        const list = get(res, `data.data`, []);
        setIsLoading(false)
        const filteredItems = list.reduce((acc, item) => {
          const idx = findIndex(acc, (o) => o.ip === item.ip && o.hardware_info.host_name === item.hardware_info.host_name)
          if (idx === -1) {
            acc.push(item)
          } else {
            const isRegisterEqual = item.worker_registered == acc[idx].worker_registered
            const isReplace = !isRegisterEqual ? item.worker_registered : item.is_online && !acc[idx].is_online
            if (isReplace) {
              acc[idx] = item;
            }
          }
          return acc
        }, [])
        setItems(filteredItems);
      })
      .catch(() => setIsLoading(false))
    return () => source.cancel();
  }, 0, 5 * 1000);

  return <>
    <LiveNodeMaps nodeList={items} />
    {/* <Validators items={items ? items.filter((item) => item.node_type != 'worker').sort((a) => a.is_validator) : null} /> */}
    <Worker items={items} />
  </>
}
