import React from "react";
import { Link } from "react-router-dom";
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import ScaleLoader from "react-spinners/ScaleLoader";
import { sharesFormatter, statusFormatter, bytesToSize } from './../../_helpers/columnFormatter'
import { formatDate } from "../../_helpers";


export default function ({
  items,
}) {
  const columns = [
    {
      dataField: "node_id",
      text: "",
      classes: 'text-center',
      style: { width: '30px' },
      formatter: (field,row,rowIndex) => rowIndex+1
    },
    {
      dataField: "ip",
      text: "Ip",
      style: { minWidth: '160px', width: '160px' },
      formatter: statusFormatter
    },
    {
      dataField: "is_validator",
      text: "Type",
      style: { width: '100px' },
      formatter: (isValidator, row) => isValidator ? 'Validator' : 'Node'
    },
    {
      dataField: "location",
      classes: 'text-truncate',
      style: { width: '160px' },
      text: "Location",
    },
    {
      dataField: "hardware_info.host_name",
      text: "Host name",
      classes: 'text-truncate',
      style: { width: '150px', maxWidth: '150px' },
    },

    {
      dataField: "hardware_info",
      text: "Hardware info",
      style: { width: '150px', maxWidth: '400px' },
      formatter: (hardware, row) => {
        const gpu = hardware.gpus ? (', ' + hardware.gpus.map((gpu) => gpu.model)) : ''
        const cpu = hardware.cpus ? hardware.cpus.map((cpu) => cpu.model) : ''
        const sizes = hardware.hard_drives ? hardware.hard_drives.reduce((acc, drive) => acc + drive.size, 0) : ''
        const storage = bytesToSize(sizes * 1000000);
        const str = cpu + gpu + ', ' + storage + ', ' + hardware.os
        return <div title={str} className='text-truncate'>{str}</div>
      },
    },
    {
      dataField: "last_update",
      headerClasses: 'text-right',
      classes: 'text-truncate text-success text-right',
      text: "Last Online",
      formatter: (unix) => formatDate(unix*1000),
    },
    
  ]
  return <>
    <div className="border mt-5">
      <Row>
        <Col xs={12}>
          <div className='bg-light  p-3'>
            <span className='h5 text-uppercase'>{items ? items.length : ''} BLOCKCHAIN NODES</span>
          </div>
        </Col>
      </Row>
      {
        items ? <div className='px-3 bg-secondary'>
          <BootstrapTable
            striped
            bootstrap4
            remote
            keyField="node_id"
            wrapperClasses="table-responsive"
            classes="table-vertical-center overflow-hidden mb-0"
            data={items}
            columns={columns}
          />
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
