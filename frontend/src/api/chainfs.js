import { fetchURL } from './utils'
import { getApiPath, doubleEncode } from '@/utils/url.js'
import { notify } from '@/notify'

export async function protectFile(source, path, hours) {
  try {
    const params = {
      path: doubleEncode(path),
      source: doubleEncode(source),
    }
    if (hours) {
      params.hours = String(hours)
    }
    const apiPath = getApiPath('api/chainfs/protect', params)
    const res = await fetchURL(apiPath, { method: 'POST' })
    return await res.json()
  } catch (err) {
    notify.showError(err.message || 'Failed to protect file')
    throw err
  }
}
