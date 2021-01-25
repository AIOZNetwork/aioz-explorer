import moment from "moment";
import BigNumber from 'bignumber.js'

export const atozToAioz = (atoz, fix = 2) => (new BigNumber(+atoz / 1e+18)).toFormat(fix)

export const lsCointToAmount = (ls) => {
  if (!ls || !ls.length) {
    return '0'
  }
  return ls.map((i) => {
    if (i.denom === 'atoz') {
      return atozToAioz(i.amount, null) + ' AIOZ'
    }
    return new Intl.NumberFormat().format(i.amount) + ` ${i.denom}`.toUpperCase()
  }).join(", ")
}

export const formatDate = (timestamp) => timestamp ? moment(new Date(timestamp)).fromNow() : '';

export const formatDateTs = (timestamp) => timestamp ? moment(new Date(timestamp)).format('YYYY-MM-DD hh:mm:ss A') : '';

