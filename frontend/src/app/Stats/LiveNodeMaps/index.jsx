import React, { useEffect, useState } from "react"
import "./LiveNodeMaps.scss"
import { omit } from "lodash";
import { GeoJSON, MapContainer } from "react-leaflet";
import Nodes from "./Nodes";
import globeGeoJson from './globe.json';
import Snap from 'snapsvg-cjs';

export default function LiveNodeMaps({ nodeList }) {
  const [svgMaps, setSvgMaps] = useState();
  const [nodes, setNodes] = useState([])
  const bounds = [[66, -175], [16, 172]]

  const onCreated = (map) => {

    setSvgMaps(Snap('.worker-nodes-maps svg'));
  }
  useEffect(() => {
    const nodes = new Map();
    (nodeList || []).map((node) => {
      nodes.set(node.node_id, omit(node, 'last_update'));
    })

    setNodes(nodes)
  }, [nodeList])

  return (<section className="live-nodes mt-5" id={"live-nodes"}>
    <div className="live-nodes-heading">
      <h2 className="title-1 mb-4">
        Live Node Maps
                </h2>

      {/* <div className="node-pannel">
        <span className="node node--worker">Node</span>
        <span className="node node--validator">Validator</span>
        <span className="node node--offline">Offline</span>
      </div> */}
    </div>

    {Boolean(globeGeoJson) &&
      <div className="worker-nodes-maps">
        <MapContainer
          bounds={bounds}
          doubleClickZoom={false}
          scrollWheelZoom={false}
          whenCreated={onCreated}
          zoomControl={false}
          touchZoom={false}
          dragging={false}>
          <GeoJSON attribution="Globe Geo" data={globeGeoJson} />
          {svgMaps && <Nodes nodes={nodes} svgMaps={svgMaps} bounds={bounds} />}
        </MapContainer>
      </div>
    }
  </section>
  )
};