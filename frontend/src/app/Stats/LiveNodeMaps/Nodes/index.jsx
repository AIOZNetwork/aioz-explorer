import { useMap, useMapEvent, useMapEvents } from "react-leaflet"
import React, { useCallback, useEffect, useState } from "react"
import { LatLng } from "leaflet"
import { isEqualWith, omit } from "lodash"
import useTimer from './../../../_helpers/useTimer'
import axios from 'axios';
import { get } from "lodash";
import "./nodes.scss"
import Snap from 'snapsvg-cjs';
import {useMediaQuery, useMediaQueries} from '@react-hook/media-query'

export default function Nodes(props) {
  const map = useMap()
  const { svgMaps, bounds } = props
  const [nodes, setNodes] = useState(props.nodes)
  const {matches} = useMediaQueries({
    screen: 'screen',
    width: '(max-width: 991px)'
  })
  const isMobile = matches && matches.width;

  const addNode = (node) => {
    const { node_id } = node
    const { color, type, size, coord, halo, opacity } = getNodeOptions(node)

    const cirCenter = svgMaps.circle(0, 0, 1).attr({
      class: "node",
      fill: color,
      opacity
    })

    const haloBlur = svgMaps.filter(Snap.filter.blur(4))

    const cirHalo = svgMaps.circle(0, 0, 1).attr({
      class: "halo",
      fill: "transparent",
      strokeWidth: 2,
      stroke: color,
      filter: haloBlur,
      opacity
    })

    const svgNode = svgMaps.group(cirCenter, cirHalo).attr({
      id: node_id,
      class: type,
      transform: `translate(${coord.x}, ${coord.y})`
    })

    svgMaps.append(svgNode)
    cirHalo.animate({ r: size * 2 }, 600, () => cirCenter.animate({ r: size }, 600))
  }

  const online = (node) => {
    const { color, size, id } = getNodeOptions(node)
    const snapGroup = Snap.select(`[id="${id}"]`)
    const snapNode = snapGroup.selectAll(".node")
    const snapHalo = snapGroup.selectAll(".halo")

    snapNode.attr({ fill: color, stroke: color })
    snapNode.animate({ opacity: 1, r: size }, 600, () => {
      snapNode.animate({ r: size }, 400)
      snapHalo.animate({ opacity: 1 }, 400)
    })
  }

  const offline = (node) => {
    const { color, size, id } = getNodeOptions(node)
    const snapGroup = Snap.select(`[id="${id}"]`)
    const snapNode = snapGroup.select(".node")
    const snapHalo = snapGroup.select(".halo")

    snapNode.attr({ fill: color, stroke: color })
    snapNode.animate({ r: size * 2 }, 600, () => {
      snapNode.animate({ opacity: 0, r: size }, 400)
      snapHalo.animate({ opacity: 0 }, 400)
    })
  }

  const getNodeOptions = (node) => {
    const { node_type, is_validator, is_online, latitude, longitude, node_id } = node

    let type = ""
    let color = ""
    let size = isMobile ? 2 : 4;
    let opacity = 1
    let halo = false

    switch (node_type) {
      case "node":
        // if (is_validator) {
        //   type = "validator"
        //   color = "#fb6a07"
        //   halo = true
        // } else {
        //   type = "blockchain"
        //   color = "#10a0de"
        // }
        // break

      default:
        type = "worker"
        color = "#7BCC3A"
        halo = true
        size = isMobile ? 1 : 3
    }

    if (!is_online) {
      color = "#999"
      size = isMobile ? 1 : 2
      opacity = 0
    }
    const coord = map.latLngToLayerPoint(
      new LatLng(
        latitude,
        longitude
      ))

    return {
      id: node_id,
      color,
      type,
      coord,
      size,
      opacity,
      halo
    }
  }

  const freshNodes = (data) => {
    const newNodes = new Map()
    const dataNodes = data.map(node => omit(node, "last_update"))
    let isChanged = false

    dataNodes.forEach((newNode) => {
      const { node_id } = newNode
      const oldNode = nodes.get(node_id)

      if (oldNode) {
        if (!isEqualWith(oldNode, newNode)) {
          console.log("Change Node Status !!!")
          newNode.is_online ? online(newNode) : offline(newNode)
          isChanged = true
        }
      } else {
        addNode(newNode)
        isChanged = true
      }

      newNodes.set(node_id, newNode)
    })

    if (isChanged) {
      setNodes(newNodes)
    }
  }

  const updateCoordNodes = useCallback(() => {
    nodes.forEach((node) => {
      const { id, coord } = getNodeOptions(node);

      Snap.select(`[id="${id}"]`).attr({
        transform: `translate(${coord.x}, ${coord.y})`
      });

    });
  }, [nodes]);

  useMapEvent('resize', () => map.fitBounds(bounds));
  useMapEvent('zoomend', () => updateCoordNodes());
  useEffect(() => nodes.forEach((node) => addNode(node)), [])
  useTimer(() => {
    const { CancelToken } = axios;
    const source = CancelToken.source();

    axios.get(`${process.env.REACT_APP_API}/node_info`, {
      cancelToken: source.token,
      params: {
        limit: 1000,
        offset: 0,
      }
    })
      .then((res) => {
        const data = get(res, `data.data`, []);
        freshNodes(data)
      })
      .catch(() => { })
    return () => source.cancel();
  }, 0, 5 * 1000);

  return <></>
}
