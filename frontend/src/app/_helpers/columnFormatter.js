import { Link } from "react-router-dom";
import { atozToAioz, lsCointToAmount } from './index';
import { ReactComponent as Aioz } from './../../assets/svg/aioz-currency-medium.svg';
import { ReactComponent as Stake } from './../../assets/svg/stake-currency-medium.svg';
import { get } from 'lodash'

const denomEnum = Object.freeze({
  aioz: 'aioz',
  atoz: 'atoz',
  stake: 'stake'
})

const msgEnum = Object.freeze({
  send: 'send',
  delegate: 'delegate',
  multiSend: 'multisend', // currently not in-use
  createValidator: 'create_validator',
  undelegate: 'undelegate',
  beginRedelegate: 'begin-redelegate',
  unjail: 'unjail',
  withdrawDelegatorReward: 'withdraw_delegator_reward',
  withdrawValidatorCommission: 'withdraw_validator_commission',

})

export const sharesFormatter = (cellContent) => `${new Intl.NumberFormat().format(Math.round(+cellContent * 100) / 100)}`

export const addressFormatter = (cellContent) => <Link className='text-truncate d-block text-lowercase' to={`/address/${cellContent}`}>{cellContent}</Link>
export const blockFormatter = (cellContent) => <Link className='text-truncate d-block' to={`/blocks/${cellContent}`}>{cellContent}</Link>
export const txnFormatter = (cellContent) => <Link className='text-truncate d-block text-lowercase' to={`/transactions/${cellContent}`}>{cellContent}</Link>

export const validFormatter = (cellContent) => cellContent ? <span className="badge badge-success text-uppercase">Success</span> : <span className="badge badge-warning text-uppercase">Fail</span>


export const feeFormatter = (cellContent) => {
  const fee = JSON.parse(cellContent)
  return coinsFormatter(fee.amount)
}

const coinFormatter = ({ denom, amount, onlyAioz }) => {
  switch (denom) {
    case denomEnum.aioz:
      return <span className='d-flex align-items-center' key={denom}>
        {new Intl.NumberFormat().format(amount)}&nbsp;<Aioz width='15' />
      </span>
    case denomEnum.atoz:
      return <span className='d-flex align-items-center' key={denom}>
        {amount >= 1e+16
          ? <>{atozToAioz(amount)}&nbsp;<Aioz width='15' /></>
          : 0
        }
      </span>
    case denomEnum.stake:
      return <span className={onlyAioz ? 'd-none' : 'd-flex align-items-center'} key={denom}>
        {new Intl.NumberFormat().format(amount)}&nbsp;
        <Stake width='15' />
      </span>

    default:
      return
  }
}

export const coinsFormatter = (cellContent, onlyAioz) => {
  try {
    const json = typeof cellContent === 'string' ? JSON.parse(cellContent) : cellContent

    if (Array.isArray(json)) {
      if (json.length === 0) {
        return '0'

      }
      return json.map(coin => coinFormatter({ ...coin, onlyAioz }))
    }
    return coinFormatter(json)

  } catch (error) {
    return '0'
  }

}

export const cointFromMsgFormatter = (msgs, row) => {
  if (msgs.length > 1) {
    return <Link className='text-truncate d-block text-capitalize' to={`/transactions/${row.Hash}`}>More</Link>
  }
  const coin = get(msgs, '[0].amount')
  return coinsFormatter(coin)
}

export const txnTypeFormatter = (type, hasMore) => {
  switch (type) {
    case msgEnum.send:
      return <div style={{ color: "#D0E658" }} className='text-truncate'><span className='mr-1 ico-transfer'></span>SEND{hasMore ? ' +' : ''}</div>
    case msgEnum.delegate:
      return <div style={{ color: "#F2994A" }} className='text-truncate'><span className='mr-1 ico-deliver'></span>DELEGATE{hasMore ? ' +' : ''}</div>
    case msgEnum.multiSend:
      return <div style={{ color: "#F2C23E" }} className='text-truncate'><span className='mr-1 ico-transfer'></span>MULTISEND{hasMore ? ' +' : ''}</div>
    case msgEnum.createValidator:
      return <div style={{ color: "#B6DF72" }} className='text-truncate'><span className='mr-1 ico-plus' style={{ fontSize: '10px' }}></span>CREATE VALIDATOR{hasMore ? ' +' : ''}</div>
    case msgEnum.undelegate:
      return <div style={{ color: "#B95585" }} className='text-truncate'><span className='mr-1 ico-trash'></span>UNDELEGATE{hasMore ? ' +' : ''}</div>
    case msgEnum.beginRedelegate:
      return <div style={{ color: "#8A5B56" }} className='text-truncate'><span className='mr-1 ico-deliver'></span>BEGIN REDELEGATE{hasMore ? ' +' : ''}</div>
    case msgEnum.unjail:
      return <div style={{ color: "#CD4445" }} className='text-truncate'><span className='mr-1 ico-store'></span>UNJAIL{hasMore ? ' +' : ''}</div>
    case msgEnum.withdrawDelegatorReward:
      return <div style={{ color: "#D4717A" }} className='text-truncate'><span className='mr-1 ico-long-arrow-alt-up-solid'></span>WITHDRAW DELEGATOR REWARD{hasMore ? ' +' : ''}</div>
    case msgEnum.withdrawValidatorCommission:
      return <div style={{ color: "#958A9C" }} className='text-truncate'><span className='mr-1 ico-long-arrow-alt-up-solid'></span>WITHDRAW VALIDATOR COMMISSION{hasMore ? ' +' : ''}</div>
    default:
      return
  }
}

export const txnsTypeFormatter = (msgs) => {
  const type = get(msgs, '[0].message_type')
  return txnTypeFormatter(type, msgs && msgs.length > 1)
}

export const statusFormatter = (cellContent, row) => {
  if(row.is_online) {
    return <><span className='rounded-circle p-1 bg-success d-inline-block'></span>&nbsp;{cellContent}</>
  }
  return <><span className='rounded-circle p-1 bg-light d-inline-block'></span>&nbsp;{cellContent}</>
}

export const bytesToSize = (bytes) => {
  var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  if (bytes == 0) return '0 Byte';
  var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
  return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
}