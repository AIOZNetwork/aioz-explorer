import React from "react";
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import BootstrapTable from 'react-bootstrap-table-next';
import ScaleLoader from "react-spinners/ScaleLoader";
import { sharesFormatter, statusFormatter, bytesToSize } from './../../_helpers/columnFormatter'
import { formatDate } from "../../_helpers";
import { MdCancel, MdCheckCircle } from "react-icons/md";
import { useState } from "react";

export default function ({
  items,
}) {
  const [filter, setFilter] = useState(1);
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
      dataField: "worker_registered",
      text: "Joined",
      classes: 'text-center',
      style: { width: '80px' },
      formatter: (isRegistered, row) => isRegistered ? <MdCheckCircle color='#35D687' /> : <MdCancel color='#D4717A' />
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
      style: { width: '200px', maxWidth: '200px' },
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
      formatter: (unix) => formatDate(unix * 1000),
    },
  ]

  const filterredItems = items ? items.filter((i) => !filter || (filter ==1 && i.is_online) || (filter==2 && !i.is_online)) : []

  return <>
    <div className="border mt-5">
      <Row>
        <Col xs={12}>
          <div className='bg-light p-3 d-flex flex-row justify-content-between align-items-center'>
            <span className='h5 text-uppercase mb-0'>{items ? filterredItems.length : ''} {filter == 1 ? 'Online ' : filter == 2 ? 'Offline ' : ''}AIOZ NODES</span>
            <div className=''>
              <Button variant="link" size='sm' className={`text-decoration-none ${!filter ? 'text-info' : ''}`} onClick={()=> setFilter(0)}>All</Button>|
              <Button variant="link" size='sm' className={`text-decoration-none ${filter == 1 ? 'text-info' : ''}`} onClick={()=> setFilter(1)}>Online</Button>|
              <Button variant="link" size='sm' className={`text-decoration-none ${filter == 2 ? 'text-info' : ''}`} onClick={()=> setFilter(2)}>Offline</Button>
            </div>
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
            data={filterredItems}
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
